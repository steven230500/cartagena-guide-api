// Datos transcritos de lib/data/{commerces,events,parishes,plans}.ts en cartagena-guide,
// ya verificados con investigación real el 2026-08-27 (ver ese repo para las fuentes).
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository/postgres"
	"github.com/steven230500/cartagena-api/internal/service"
)

func strp(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está seteada")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a Postgres: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	businessSvc := service.NewBusinessService(postgres.NewBusinessRepository(queries))
	eventSvc := service.NewEventService(postgres.NewEventRepository(queries))
	parishSvc := service.NewParishService(postgres.NewParishRepository(queries))
	planSvc := service.NewPlanService(postgres.NewPlanRepository(queries))
	routeSvc := service.NewRouteService(postgres.NewRouteRepository(queries))
	achievementSvc := service.NewAchievementService(
		postgres.NewAchievementRepository(queries),
		postgres.NewFavoriteRepository(queries),
		postgres.NewRouteProgressRepository(queries),
	)

	seedBusinesses(ctx, businessSvc)
	businessBySlug := indexBusinesses(ctx, businessSvc)
	seedEvents(ctx, eventSvc, businessBySlug)
	seedParishes(ctx, parishSvc)
	seedPlans(ctx, planSvc)
	seedRoutes(ctx, routeSvc)
	seedAchievements(ctx, achievementSvc)

	log.Println("seed completo")
}

func seedBusinesses(ctx context.Context, svc *service.BusinessService) {
	businesses := []domain.Business{
		{
			Name: "La Cevichería", Slug: "la-cevicheria-cartagena", Type: "food", Subtype: "seafood",
			Barrio: "San Diego", Lat: 10.4239, Lng: -75.5504,
			Image:       "/traditional-ceviche-restaurant-cartagena.jpg",
			Tags:        []string{"mariscos", "ceviche", "caribe", "fresco"},
			Description: "Famoso restaurante de mariscos y ceviches frescos del Caribe.",
			Hours:       strp("12:00 - 23:00"), PriceHint: strp("$$"),
			Phone: strp("+57 314 894 9877"), Web: strp("https://lacevicheriacartagena.com"),
			Email: strp("lacevicheriacartagena@gmail.com"), Instagram: strp("@lacevicheria"),
		},
		{
			Name: "Restaurante 1621", Slug: "restaurante-1621", Type: "food", Subtype: "fine-dining",
			Barrio: "San Diego", Lat: 10.4261, Lng: -75.5499,
			Image:       "/caribbean-food-festival-cartagena.jpg",
			Tags:        []string{"alta-cocina", "caribeña", "hotel", "sofisticado"},
			Description: "Alta cocina caribeña en el interior del Hotel Sofitel Legend Santa Clara.",
			Hours:       strp("18:30 - 23:30 (miércoles a domingo)"), PriceHint: strp("$$$"),
			Phone: strp("+57 311 491 2308"),
			Web:   strp("https://www.sofitellegendsantaclara.com/restaurants-bars/1621-restaurant/"),
			Email: strp("reservations.santaclara@sofitel.com"),
		},
		{
			// No es un solo negocio: portal colonial con decenas de puestos independientes.
			// Sin teléfono/web/instagram oficial único — no inventar uno.
			Name: "Portal de los Dulces", Slug: "portal-de-los-dulces", Type: "artisan", Subtype: "sweets",
			Barrio: "Centro Histórico", Lat: 10.4226, Lng: -75.5492,
			Image: "/cartagena-colonial-route-historic-buildings.jpg",
			Tags:  []string{"dulces", "tradicional", "artesanal", "típico"},
			Description: "Portal colonial frente a la Puerta del Reloj, en Plaza de los Coches. " +
				"Decenas de puestos independientes de dulces típicos: cocadas, alegrías, caballitos.",
			Hours: strp("08:00 - 18:00"), PriceHint: strp("$"),
		},
		{
			// Igual que Portal de los Dulces: complejo de ~20-25 tiendas independientes.
			Name: "Las Bóvedas", Slug: "las-bovedas-cartagena", Type: "artisan", Subtype: "crafts",
			Barrio: "San Diego", Lat: 10.429, Lng: -75.545,
			Image: "/cartagena-city-walls-fortifications.jpg",
			Tags:  []string{"artesanías", "colonial", "histórico", "galerías"},
			Description: "Antiguos calabozos militares (1798) entre los baluartes de Santa Catalina y " +
				"Santa Clara, hoy ~20-25 tiendas independientes de artesanías y souvenirs.",
			Hours: strp("09:00 - 19:00"), PriceHint: strp("$$"),
		},
		{
			Name: "Caribe Plaza", Slug: "caribe-plaza", Type: "mall", Subtype: "shopping",
			Barrio: "Manga", Lat: 10.4152, Lng: -75.5358,
			Image:       "/modern-shopping-mall-cartagena.jpg",
			Tags:        []string{"centro-comercial", "tiendas", "cine", "restaurantes"},
			Description: "Centro comercial moderno con tiendas de moda, cine y restaurantes.",
			Hours:       strp("Dom-Jue 10:00-20:00, Vie-Sáb 10:00-21:00"), PriceHint: strp("$$"),
			Phone: strp("+57 605 669 2332"), Web: strp("https://www.cccaribeplaza.com"),
		},
		{
			Name: "Mall Plaza El Castillo", Slug: "mall-plaza-el-castillo", Type: "mall", Subtype: "shopping",
			Barrio: "Castillogrande", Lat: 10.4143, Lng: -75.5329,
			Image:       "/bocagrande-modern-buildings-beach.jpg",
			Tags:        []string{"centro-comercial", "mar", "exclusivo", "comida"},
			Description: "Centro comercial frente al mar con tiendas exclusivas y zona de comida.",
			Hours:       strp("Lun-Sáb 10:00-21:00, Dom 11:00-20:00"), PriceHint: strp("$$"),
			Phone: strp("+57 605 660 2020"), Web: strp("https://www.mallplaza.com/co/cartagena"),
		},
		{
			Name: "Café del Mar", Slug: "cafe-del-mar-cartagena", Type: "shop", Subtype: "cocktail",
			Barrio: "Centro Histórico", Lat: 10.4233, Lng: -75.554,
			Image:       "/cartagena-de-indias-colonial-walls-sunset-caribbea.jpg",
			Tags:        []string{"bar", "murallas", "atardecer", "vistas", "icónico"},
			Description: "Bar icónico sobre el Baluarte de Santo Domingo, en las murallas, con vistas al atardecer.",
			Hours:       strp("16:30 - 23:00"), PriceHint: strp("$$"),
			Phone: strp("+57 605 664 2945"), Web: strp("https://www.cafedelmarcartagena.com.co"),
			Email: strp("reservas@cafedelmarcartagena.com.co"),
		},
		{
			// OJO: fuentes recientes dicen que reabrió en 2024 reubicado en Plaza Santo Domingo
			// (Casa Cruxada), no en Getsemaní. Verificar antes de publicar cuál dirección vale hoy.
			Name: "Bazurto Social Club", Slug: "bazurto-social-club", Type: "shop", Subtype: "club",
			Barrio: "Getsemaní", Lat: 10.4252, Lng: -75.5478,
			Image:       "/champeta-concert-cartagena-walls-sunset.jpg",
			Tags:        []string{"champeta", "salsa", "caribeño", "emblemático", "música"},
			Description: "Club emblemático que mezcla champeta, salsa y sonidos caribeños.",
			Hours:       strp("20:00 - 04:00"), PriceHint: strp("$"),
			Phone: strp("+57 317 648 1183"), Web: strp("https://www.bazurtosocialclub.com"),
		},
		{
			// No se encontró teléfono ni web propia confirmados — vacío en vez de inventar.
			Name: "Casa Abba", Slug: "casa-abba-gallery", Type: "culture", Subtype: "contemporary",
			Barrio: "San Diego", Lat: 10.4281, Lng: -75.547,
			Image: "/art-fair-getsemani-cartagena-street-art.jpg",
			Tags:  []string{"galería", "arte", "contemporáneo", "local", "internacional"},
			Description: "Galería-boutique con arte, zapatos y bolsos pintados a mano, joyería y " +
				"pintura de artistas y artesanos locales.",
			Hours: strp("10:00 - 18:00"), PriceHint: strp("$$"),
		},
		{
			Name: "St. Dom Boutique", Slug: "st-dom-boutique", Type: "shop", Subtype: "luxury",
			Barrio: "Centro Histórico", Lat: 10.4224, Lng: -75.5507,
			Image: "/puerta-del-reloj-cartagena-clock-tower.jpg",
			Tags:  []string{"boutique", "lujo", "diseño", "latinoamericano", "sostenible"},
			Description: "Concept store de 700m² en casa colonial de 300 años: moda, accesorios y " +
				"decoración de diseñadores colombianos.",
			Hours: strp("Lun-Sáb 10:00-20:00, Dom 12:00-20:00"), PriceHint: strp("$$$"),
			Phone: strp("+57 605 664 0197"), Instagram: strp("@stdomcartagena"),
		},
	}

	for _, b := range businesses {
		if _, err := svc.Create(ctx, b); err != nil {
			log.Fatalf("seed business %s: %v", b.Name, err)
		}
	}
	log.Printf("sembrados %d comercios", len(businesses))
}

// indexBusinesses arma slug -> id para poder linkear eventos con su comercio relacionado.
func indexBusinesses(ctx context.Context, svc *service.BusinessService) map[string]domain.Business {
	all, err := svc.List(ctx, domain.BusinessFilter{})
	if err != nil {
		log.Fatalf("listar comercios para indexar: %v", err)
	}
	out := make(map[string]domain.Business, len(all))
	for _, b := range all {
		out[b.Slug] = b
	}
	return out
}

func seedEvents(ctx context.Context, svc *service.EventService, businessBySlug map[string]domain.Business) {
	type seedEvent struct {
		title, slug, eventType, venue, image, description, content string
		start                                                      string
		end                                                        string
		relatedSlug                                                string
		lat, lng                                                   float64
		tags                                                       []string
	}

	events := []seedEvent{
		{
			title: "Feria de Arte en Getsemaní", slug: "feria-arte-getsemani",
			start: "2025-07-12", end: "2025-07-14",
			eventType: "feria", venue: "Plaza de la Trinidad",
			// "taller-mojojoy" en el dato original no corresponde a ninguno de los 10
			// comercios verificados hoy — se deja sin relacionar en vez de inventar el link.
			relatedSlug: "",
			lat:         10.4212, lng: -75.5483,
			image:       "/art-fair-getsemani-cartagena-street-art.jpg",
			tags:        []string{"arte", "calle", "local"},
			description: "Feria de arte local con artistas de Getsemaní y talleres interactivos.",
			content: "Una celebración del arte local donde artistas del barrio Getsemaní exponen sus obras. " +
				"Incluye talleres de pintura, música en vivo y gastronomía típica.",
		},
		{
			title: "Festival Gastronómico del Caribe", slug: "festival-gastronomico-caribe",
			start: "2025-08-15", end: "2025-08-18",
			eventType: "gastronomia", venue: "Plaza Santo Domingo",
			lat: 10.4242, lng: -75.5489,
			image:       "/caribbean-food-festival-cartagena.jpg",
			tags:        []string{"gastronomía", "caribe", "cultura"},
			description: "Festival que celebra la rica gastronomía caribeña con chefs locales e internacionales.",
			content: "Cuatro días de celebración culinaria con los mejores chefs del Caribe. " +
				"Degustaciones, talleres de cocina y presentaciones culturales.",
		},
		{
			title: "Concierto de Champeta en las Murallas", slug: "concierto-champeta-murallas",
			start:     "2025-09-21",
			eventType: "concierto", venue: "Murallas de Cartagena",
			lat: 10.4234, lng: -75.5512,
			image:       "/champeta-concert-cartagena-walls-sunset.jpg",
			tags:        []string{"música", "champeta", "murallas"},
			description: "Concierto de champeta al atardecer con vista al mar Caribe.",
			content: "Una noche mágica de champeta, el ritmo más auténtico de Cartagena, con las murallas " +
				"históricas como escenario y el atardecer caribeño de fondo.",
		},
	}

	for _, e := range events {
		event := domain.Event{
			Title: e.title, Slug: e.slug, StartDate: e.start,
			Type: e.eventType, Venue: e.venue, Lat: e.lat, Lng: e.lng, Image: e.image,
			Tags: e.tags, Description: e.description, Content: strp(e.content),
		}
		if e.end != "" {
			event.EndDate = strp(e.end)
		}
		if e.relatedSlug != "" {
			if b, ok := businessBySlug[e.relatedSlug]; ok {
				event.RelatedBusinessID = &b.ID
			}
		}
		if _, err := svc.Create(ctx, event); err != nil {
			log.Fatalf("seed event %s: %v", e.title, err)
		}
	}
	log.Printf("sembrados %d eventos", len(events))
}

func seedParishes(ctx context.Context, svc *service.ParishService) {
	type seedParish struct {
		name, address, neighborhood, phone string
		schedules                          []domain.Schedule
	}

	parishes := []seedParish{
		{
			name: "Catedral de Santa Catalina de Alejandría", address: "Plaza de la Proclamación, Centro Histórico",
			neighborhood: "Centro Histórico", phone: "+57 5 664 3299",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"7:00 AM", "12:00 PM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"7:00 AM", "12:00 PM", "6:00 PM", "7:00 PM"}},
				{Day: "Domingos", Times: []string{"7:00 AM", "9:00 AM", "11:00 AM", "12:00 PM", "6:00 PM", "7:00 PM"}},
			},
		},
		{
			name: "Iglesia de San Pedro Claver", address: "Calle de la Factoria, Centro Histórico",
			neighborhood: "Centro Histórico", phone: "+57 5 664 4991",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"7:00 AM", "12:00 PM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Domingos", Times: []string{"8:00 AM", "10:00 AM", "12:00 PM", "6:00 PM"}},
			},
		},
		{
			name: "Iglesia de Santo Domingo", address: "Plaza de Santo Domingo, Centro Histórico",
			neighborhood: "Centro Histórico", phone: "+57 5 664 3965",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Domingos", Times: []string{"8:00 AM", "10:00 AM", "12:00 PM", "6:00 PM"}},
			},
		},
		{
			name: "Iglesia de la Trinidad", address: "Barrio Getsemaní",
			neighborhood: "Getsemaní", phone: "+57 5 664 2715",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"6:30 AM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"6:30 AM", "6:00 PM"}},
				{Day: "Domingos", Times: []string{"7:00 AM", "9:00 AM", "11:00 AM", "6:00 PM"}},
			},
		},
		{
			name: "Iglesia de San Francisco", address: "Plaza de San Francisco, Centro Histórico",
			neighborhood: "Centro Histórico", phone: "",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"7:00 AM", "12:00 PM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Domingos", Times: []string{"8:00 AM", "10:00 AM", "12:00 PM", "6:00 PM"}},
			},
		},
		{
			name: "Santuario de San Roque", address: "Barrio San Diego",
			neighborhood: "San Diego", phone: "+57 5 664 3643",
			schedules: []domain.Schedule{
				{Day: "Lunes a Viernes", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Sábados", Times: []string{"7:00 AM", "6:00 PM"}},
				{Day: "Domingos", Times: []string{"8:00 AM", "10:00 AM", "12:00 PM", "6:00 PM"}},
			},
		},
	}

	for _, p := range parishes {
		parish := domain.Parish{
			Name: p.name, Address: p.address, Neighborhood: p.neighborhood,
			Phone: strp(p.phone), Schedules: p.schedules,
		}
		if _, err := svc.Create(ctx, parish); err != nil {
			log.Fatalf("seed parish %s: %v", p.name, err)
		}
	}
	log.Printf("sembradas %d parroquias", len(parishes))
}

func seedPlans(ctx context.Context, svc *service.PlanService) {
	plans := []domain.Plan{
		{
			Title: "Concierto de Champeta en las Murallas",
			Description: "Disfruta de música champeta en vivo con vista al mar Caribe. " +
				"Un plan perfecto para el fin de semana.",
			Type: "cultural", Price: "free", Date: "Sábados", Time: "6:00 PM - 9:00 PM",
			Location: "Murallas de Cartagena", Neighborhood: "Centro Histórico",
		},
		{
			Title: "Tour Gastronómico por Getsemaní",
			Description: "Recorre los mejores restaurantes y puestos de comida callejera del barrio " +
				"más bohemio de Cartagena.",
			Type: "gastronomic", Price: "paid", PriceAmount: strp("$80.000 COP"),
			Date: "Domingos", Time: "10:00 AM - 2:00 PM",
			Location: "Getsemaní", Neighborhood: "Getsemaní",
		},
		{
			Title:       "Yoga al Amanecer en Bocagrande",
			Description: "Sesión de yoga frente al mar para comenzar el día con energía positiva.",
			Type:        "outdoor", Price: "free", Date: "Sábados y Domingos", Time: "6:30 AM - 7:30 AM",
			Location: "Playa de Bocagrande", Neighborhood: "Bocagrande",
		},
		{
			Title:       "Noche de Salsa en Café Havana",
			Description: "Baila salsa y disfruta de cócteles en uno de los bares más icónicos de Cartagena.",
			Type:        "nightlife", Price: "paid", PriceAmount: strp("$50.000 COP (cover)"),
			Date: "Viernes y Sábados", Time: "9:00 PM - 2:00 AM",
			Location: "Café Havana", Neighborhood: "Centro Histórico",
		},
		{
			Title:       "Mercado de Artesanías en Las Bóvedas",
			Description: "Explora artesanías locales, joyería y souvenirs únicos en las históricas bóvedas.",
			Type:        "shopping", Price: "free", Date: "Todos los días", Time: "9:00 AM - 6:00 PM",
			Location: "Las Bóvedas", Neighborhood: "Centro Histórico",
		},
		{
			Title:       "Cine al Aire Libre en Parque Centenario",
			Description: "Proyección de películas colombianas bajo las estrellas. Trae tu manta y disfruta.",
			Type:        "cultural", Price: "free", Date: "Primer viernes del mes", Time: "7:00 PM - 10:00 PM",
			Location: "Parque Centenario", Neighborhood: "Centro",
		},
	}

	for _, p := range plans {
		if _, err := svc.Create(ctx, p); err != nil {
			log.Fatalf("seed plan %s: %v", p.Title, err)
		}
	}
	log.Printf("sembrados %d planes", len(plans))
}

func latlng(lat, lng float64) (*float64, *float64) {
	return &lat, &lng
}

// seedRoutes transcribe las 6 rutas de components/routes/routes-grid.tsx.
// rating/reviews NO se migran — eran números inventados. Solo "Centro Histórico
// Esencial" trae sus 3 pasos reales (de app/routes/[id]/page.tsx); las otras 5
// quedan con su resumen y sin pasos, para completar después desde el admin.
func seedRoutes(ctx context.Context, svc *service.RouteService) {
	lat1, lng1 := latlng(10.4236, -75.5518)
	lat2, lng2 := latlng(10.4242, -75.5516)
	lat3, lng3 := latlng(10.4248, -75.5512)

	routes := []domain.Route{
		{
			Slug: "centro-historico-esencial", Title: "Centro Histórico Esencial",
			Description: "Recorre los sitios más emblemáticos de la ciudad amurallada",
			Duration:    "2-3 horas", Distance: "3.2 km", Difficulty: "Fácil", Category: "Historia",
			Image:      "/cartagena-colonial-route-historic-buildings.jpg",
			Highlights: []string{"Plaza de Armas", "Catedral", "Murallas", "Balcones Coloniales"},
			AudioGuide: true, Offline: true, Price: "Gratis",
			Steps: []domain.RouteStep{
				{
					Title: "Plaza de Armas", Description: "Punto de partida en el corazón de la ciudad colonial",
					AudioURL: strp("/audio/plaza-armas.mp3"), Duration: strp("8 min"),
					Lat: lat1, Lng: lng1, Image: strp("/placeholder.svg?key=plazaarmas"),
				},
				{
					Title: "Catedral de Cartagena", Description: "Majestuosa catedral construida en el siglo XVI",
					AudioURL: strp("/audio/catedral.mp3"), Duration: strp("12 min"),
					Lat: lat2, Lng: lng2, Image: strp("/placeholder.svg?key=catedral"),
				},
				{
					Title: "Murallas Coloniales", Description: "Camina por las históricas fortificaciones",
					AudioURL: strp("/audio/murallas.mp3"), Duration: strp("15 min"),
					Lat: lat3, Lng: lng3, Image: strp("/placeholder.svg?key=murallas"),
				},
			},
		},
		{
			Slug: "sabores-del-caribe", Title: "Sabores del Caribe",
			Description: "Degusta la auténtica gastronomía cartagenera",
			Duration:    "3-4 horas", Distance: "2.8 km", Difficulty: "Fácil", Category: "Gastronomía",
			Image:      "/traditional-ceviche-restaurant-cartagena.jpg",
			Highlights: []string{"Mercado Bazurto", "Portal de los Dulces", "Cevichería", "Café del Mar"},
			AudioGuide: true, Offline: true, Price: "$25 USD",
		},
		{
			Slug: "palenque-y-libertad", Title: "Palenque y Libertad",
			Description: "Conoce la herencia afrodescendiente y la lucha por la libertad",
			Duration:    "2.5 horas", Distance: "2.1 km", Difficulty: "Fácil", Category: "Cultura",
			Image: "/palenquera-woman-traditional-dress-cartagena.jpg",
			Highlights: []string{
				"Monumento a la India Catalina", "Barrio San Diego", "Casa de la Cultura", "Plaza de la Trinidad",
			},
			AudioGuide: true, Offline: true, Price: "Gratis",
		},
		{
			Slug: "murallas-al-atardecer", Title: "Murallas al Atardecer",
			Description: "Camina por las fortificaciones mientras el sol se oculta en el Caribe",
			Duration:    "1.5 horas", Distance: "4.5 km", Difficulty: "Moderado", Category: "Paisaje",
			Image: "/cartagena-de-indias-colonial-walls-sunset-caribbea.jpg",
			Highlights: []string{
				"Baluarte de San Francisco", "Café del Mar", "Baluarte de Santa Catalina", "Plaza de los Coches",
			},
			AudioGuide: true, Offline: true, Price: "Gratis",
		},
		{
			Slug: "getsemani-bohemio", Title: "Getsemaní Bohemio",
			Description: "Explora el barrio más vibrante y artístico de Cartagena",
			Duration:    "2 horas", Distance: "1.8 km", Difficulty: "Fácil", Category: "Arte",
			Image:      "/getsemani-street-art-colorful-buildings.jpg",
			Highlights: []string{"Plaza de la Trinidad", "Calle del Arsenal", "Murales Urbanos", "Bares Locales"},
			AudioGuide: true, Offline: true, Price: "Gratis",
		},
		{
			Slug: "castillo-san-felipe", Title: "Castillo San Felipe",
			Description: "Descubre la fortaleza más imponente de América",
			Duration:    "1.5 horas", Distance: "1.2 km", Difficulty: "Moderado", Category: "Fortaleza",
			Image:      "/castillo-san-felipe-fortress-cartagena.jpg",
			Highlights: []string{"Túneles Subterráneos", "Cañones Históricos", "Vista Panorámica", "Museo Militar"},
			AudioGuide: true, Offline: true, Price: "$8 USD",
		},
	}

	for _, r := range routes {
		if _, err := svc.Create(ctx, r); err != nil {
			log.Fatalf("seed route %s: %v", r.Title, err)
		}
	}
	log.Printf("sembradas %d rutas", len(routes))
}

// seedAchievements: insignias que se desbloquean solas por acciones reales
// (rutas completadas, favoritos) — nada de ubicaciones/pistas inventadas.
func seedAchievements(ctx context.Context, svc *service.AchievementService) {
	achievements := []domain.Achievement{
		{
			Code: "primera-ruta", Title: "Primera ruta completada",
			Description: "Completaste tu primera ruta guiada por Cartagena.",
			Icon:        "route", CriteriaType: "routes_completed", Threshold: 1,
		},
		{
			Code: "explorador-rutas", Title: "Explorador de rutas",
			Description: "Completaste 3 rutas guiadas distintas.",
			Icon:        "compass", CriteriaType: "routes_completed", Threshold: 3,
		},
		{
			Code: "coleccionista", Title: "Coleccionista",
			Description: "Guardaste 5 comercios como favoritos.",
			Icon:        "heart", CriteriaType: "favorites_count", Threshold: 5,
		},
		{
			Code: "fan-cartagena", Title: "Fan de Cartagena",
			Description: "Guardaste 15 comercios como favoritos.",
			Icon:        "star", CriteriaType: "favorites_count", Threshold: 15,
		},
	}

	for _, a := range achievements {
		if _, err := svc.Create(ctx, a); err != nil {
			log.Fatalf("seed achievement %s: %v", a.Title, err)
		}
	}
	log.Printf("sembrados %d logros", len(achievements))
}
