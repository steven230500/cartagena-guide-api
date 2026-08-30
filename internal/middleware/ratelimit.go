package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// clientIP saca la IP real detrás de Caddy leyendo X-Forwarded-For directo —
// no depende de gin.Context.ClientIP()/TrustedProxies, porque la API solo es
// alcanzable dentro de la red docker interna (nunca directo desde internet),
// así que el "proxy" siempre es Caddy y no hace falta la lista de confianza
// que gin usaría para decidir si mirar el header o no.
func clientIP(c *gin.Context) string {
	if fwd := c.GetHeader("X-Forwarded-For"); fwd != "" {
		if ip := strings.TrimSpace(strings.Split(fwd, ",")[0]); ip != "" {
			return ip
		}
	}
	return c.ClientIP()
}

type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	burst    int
}

func newIPLimiter(r rate.Limit, burst int) *ipLimiter {
	l := &ipLimiter{limiters: make(map[string]*rate.Limiter), r: r, burst: burst}
	go l.cleanupLoop()
	return l
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(l.r, l.burst)
		l.limiters[ip] = lim
	}
	return lim
}

// cleanupLoop evita que el mapa crezca sin límite en un proceso de larga
// duración — una IP que agotó su burst hace rato ya está "llena" de nuevo,
// así que es seguro tirarla y recrearla desde cero la próxima vez que pegue.
func (l *ipLimiter) cleanupLoop() {
	for range time.Tick(10 * time.Minute) {
		l.mu.Lock()
		for ip, lim := range l.limiters {
			if lim.Tokens() >= float64(l.burst) {
				delete(l.limiters, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RequireRateLimit limita a `perMinute` requests por minuto por IP, con un
// burst inicial del mismo tamaño — pensado para endpoints de auth públicos
// (login/register) donde no hace falta memoria entre reinicios del proceso.
func RequireRateLimit(perMinute int) gin.HandlerFunc {
	limiter := newIPLimiter(rate.Every(time.Minute/time.Duration(perMinute)), perMinute)

	return func(c *gin.Context) {
		if !limiter.get(clientIP(c)).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "demasiados intentos, probá de nuevo en un rato"})
			c.Abort()
			return
		}
		c.Next()
	}
}
