package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/webhelp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "modernc.org/sqlite"

	"domi.ninja/example-project/internal/db_generated"
	"github.com/patrickmn/go-cache"
)

type App struct {
	cfg *webhelp.AppConfig

	db *db_generated.Queries

	version  string
	memcache *cache.Cache
}

func Run() {
	// pretty verbose logger but it helps to find where panics happend as well as logging all request
	webhelp.UseLogger()

	// toml config file for set settings, url etc.
	cfg := webhelp.MustLoadConfig("./app.toml")

	// db connection
	dbConn, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal("db error ", err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// give our http route handlers the stuff they need
	app := &App{
		db:  db_generated.New(dbConn),
		cfg: cfg,

		version:  webhelp.BuildRandomNumber,
		memcache: cache.New(5*time.Minute, 10*time.Minute),
	}

	router := chi.NewRouter()

	// define golang app logger
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(webhelp.LoggerMiddleware)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// maybe some jwt auth stuff on the router?
	// router.Use(jwtauth.Verifier(webhelp.tokenAuth))

	// Serve static files directory
	fileServerStatic := http.FileServer(http.Dir("./frontend/assets"))
	fileServerUploads := http.FileServer(http.Dir("./uploads"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fileServerStatic))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServerUploads))

	// Authentication routes (public)
	router.Get("/login", app.HandleLogin_GET)
	router.Get("/register", app.HandleRegister_GET)
	router.Post("/auth/register", app.HandleRegister_POST)
	router.Post("/auth/login", app.HandleLogin_POST)
	router.Post("/auth/logout", app.HandleLogout_POST)
	router.Get("/auth/me", app.HandleMe_GET)

	// parameterised routes need to be registered before the root route
	router.Get("/post/{postId}", app.HandlePost_PostId_GET)
	router.Delete("/post/{postId}", app.HandlePost_PostId_DELETE)

	// crud routes
	router.Post("/posts", app.HandlePosts_POST)

	// Protected routes (require authentication)
	router.Group(func(r chi.Router) {
		r.Use(app.RequireAuth)
		// Add protected routes here, for example:
		// r.Get("/dashboard", app.HandleDashboard_GET)
		// r.Post("/posts", app.HandlePosts_POST) // if you want posts to be protected
	})

	if webhelp.DevMode() {
		router.Get("/reload", app.HandleReload_WS)
	}

	// root routes
	router.Get("/", app.HandleIndex)

	// health check route
	router.Get("/health", app.HandleHealth)

	// maybe some admin routes on /admin?
	// Bundle your routes in a separate file and chi route bundle thing
	// router.Route("/admin", adminRoutes)

	// Listen
	bindAddr := app.cfg.Server.BindAddress + ":" + fmt.Sprint(app.cfg.Server.Port)
	server := &http.Server{
		Handler: router,
		Addr:    bindAddr,
	}
	log.Print("Webserver starting")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("http.Server failed to start: ", err)
	}
}
