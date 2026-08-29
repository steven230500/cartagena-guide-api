package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/handler"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/repository/postgres"
	"github.com/steven230500/cartagena-api/internal/service"
	"github.com/steven230500/cartagena-api/migrations"
)

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está seteada")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET no está seteada")
	}

	if err := migrations.Run(databaseURL); err != nil {
		log.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a Postgres: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	// repository (postgres) -> service -> handler, por recurso.
	businessRepo := postgres.NewBusinessRepository(queries)
	businessHandler := handler.NewBusinessHandler(service.NewBusinessService(businessRepo))
	eventHandler := handler.NewEventHandler(service.NewEventService(postgres.NewEventRepository(queries)))
	parishHandler := handler.NewParishHandler(service.NewParishService(postgres.NewParishRepository(queries)))
	planHandler := handler.NewPlanHandler(service.NewPlanService(postgres.NewPlanRepository(queries)))
	routeHandler := handler.NewRouteHandler(service.NewRouteService(postgres.NewRouteRepository(queries)))

	userRepo := postgres.NewUserRepository(queries)
	authSvc := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authSvc)
	favoriteRepo := postgres.NewFavoriteRepository(queries)
	favoriteHandler := handler.NewFavoriteHandler(service.NewFavoriteService(favoriteRepo))
	routeProgressRepo := postgres.NewRouteProgressRepository(queries)
	routeProgressHandler := handler.NewRouteProgressHandler(service.NewRouteProgressService(routeProgressRepo))
	achievementHandler := handler.NewAchievementHandler(
		service.NewAchievementService(postgres.NewAchievementRepository(queries), favoriteRepo, routeProgressRepo),
	)
	businessClaimHandler := handler.NewBusinessClaimHandler(
		service.NewBusinessClaimService(postgres.NewBusinessClaimRepository(queries), businessRepo, userRepo),
	)

	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("no se pudo configurar trusted proxies: %v", err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.GET("/businesses", businessHandler.List)
		api.GET("/businesses/:slug", businessHandler.GetBySlug)
		api.GET("/events", eventHandler.List)
		api.GET("/events/:slug", eventHandler.GetBySlug)
		api.GET("/parishes", parishHandler.List)
		api.GET("/plans", planHandler.List)
		api.GET("/routes", routeHandler.List)
		api.GET("/routes/:slug", routeHandler.GetBySlug)
		api.GET("/achievements", achievementHandler.List)

		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		authed := api.Group("")
		authed.Use(middleware.RequireUser(authSvc))
		{
			authed.GET("/me", authHandler.Me)
			authed.GET("/me/favorites", favoriteHandler.List)
			authed.POST("/me/favorites/:business_id", favoriteHandler.Add)
			authed.DELETE("/me/favorites/:business_id", favoriteHandler.Remove)
			authed.GET("/me/route-progress/:route_id", routeProgressHandler.Get)
			authed.PUT("/me/route-progress/:route_id", routeProgressHandler.Upsert)
			authed.GET("/me/achievements", achievementHandler.Progress)
			authed.GET("/me/businesses", businessHandler.ListMine)
			authed.PUT("/me/businesses/:id", businessHandler.UpdateMine)
			authed.GET("/me/business-claims", businessClaimHandler.ListMine)
			authed.POST("/me/business-claims", businessClaimHandler.Create)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAdminKey())
		{
			admin.POST("/businesses", businessHandler.Create)
			admin.PUT("/businesses/:id", businessHandler.Update)
			admin.DELETE("/businesses/:id", businessHandler.Delete)

			admin.POST("/events", eventHandler.Create)
			admin.PUT("/events/:id", eventHandler.Update)
			admin.DELETE("/events/:id", eventHandler.Delete)

			admin.POST("/parishes", parishHandler.Create)
			admin.PUT("/parishes/:id", parishHandler.Update)
			admin.DELETE("/parishes/:id", parishHandler.Delete)

			admin.POST("/plans", planHandler.Create)
			admin.PUT("/plans/:id", planHandler.Update)
			admin.DELETE("/plans/:id", planHandler.Delete)

			admin.POST("/routes", routeHandler.Create)
			admin.PUT("/routes/:id", routeHandler.Update)
			admin.DELETE("/routes/:id", routeHandler.Delete)

			admin.POST("/achievements", achievementHandler.Create)
			admin.PUT("/achievements/:id", achievementHandler.Update)
			admin.DELETE("/achievements/:id", achievementHandler.Delete)

			admin.GET("/business-claims", businessClaimHandler.ListPending)
			admin.POST("/business-claims/:id/approve", businessClaimHandler.Approve)
			admin.POST("/business-claims/:id/reject", businessClaimHandler.Reject)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("escuchando en :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
