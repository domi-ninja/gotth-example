package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "modernc.org/sqlite"

	"domi.ninja/example-project/internal/db_generated"
	"domi.ninja/example-project/webhelp"
	"github.com/patrickmn/go-cache"
)

type App struct {
	webhelp.Wapp

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
		db: db_generated.New(dbConn),

		Wapp: webhelp.Wapp{
			Cfg: cfg,
		},

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
		AllowedOrigins:   []string{"https://*", "ht  tp://*"},
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
	router.Get("/login", app.HandleLogin_VIEW)
	router.Post("/login", app.HandleLogin_POST)

	router.Get("/register", app.HandleRegister_VIEW)
	router.Post("/register", app.HandleRegister_POST)

	router.Get("/logout", app.HandleLogout_GET)

	// parameterised routes need to be registered before the root route
	router.Get("/post/{postId}", app.HandlePost_PostId_VIEW)

	// Protected routes (require authentication)
	router.Group(func(authR chi.Router) {
		authR.Use(app.RequireAuth)

		authR.Get("/me", app.HandleMe_GET)

		// crud routes
		router.Post("/posts", app.HandlePosts_POST)

		// manage
		// TODO more stuff here ... authR.Delete("/post/{postId}", app.HandlePost_PostId_DELETE)
		authR.Delete("/post/{postId}", app.HandlePost_PostId_DELETE)

	})

	if webhelp.DevMode() {
		router.Get("/reload", app.HandleReload_WS)
	}

	// root routes
	router.Get("/", app.HandleIndex_VIEW)

	// health check route
	router.Get("/health", app.HandleHealth)

	// Listen
	bindAddr := app.Cfg.Server.BindAddress + ":" + fmt.Sprint(app.Cfg.Server.Port)
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
