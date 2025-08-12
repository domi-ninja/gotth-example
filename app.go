package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/webapp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "modernc.org/sqlite"

	"domi.ninja/example-project/internal/db_generated"
	"github.com/patrickmn/go-cache"
)

type WebApp struct {
	db                *db_generated.Queries
	memcache          *cache.Cache
	siteCfg           *webapp.AppConfig
	buildRandomNumber string
	appPath           string
}

func Run() {
	webapp.UseLogger()

	// toml config file for set settings, url etc.
	cfg := webapp.MustLoadConfig("./app.toml")

	// db connection
	dbConn, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal("db error ", err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// give our http route handlers the stuff they need
	app := &WebApp{
		db:                db_generated.New(dbConn),
		siteCfg:           cfg,
		buildRandomNumber: webapp.BuildRandomNumber,
	}

	router := chi.NewRouter()

	// define golang app logger
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(webapp.LoggerMiddleware)
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
	// router.Use(jwtauth.Verifier(webapp.tokenAuth))

	// Serve static files directory
	fileServerStatic := http.FileServer(http.Dir("./frontend/assets"))
	fileServerUploads := http.FileServer(http.Dir("./uploads"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fileServerStatic))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServerUploads))

	// parameterised routes need to be registered before the root route
	router.Get("/post/{postId}", app.HandlePost_PostId_GET)
	router.Delete("/post/{postId}", app.HandlePost_DELETE)

	// root routes
	router.Get("/", app.HandleIndex)

	//
	router.Post("/posts", app.HandlePosts_POST)

	// maybe some admin routes on /admin?
	// router.Route("/admin", webapp.adminRoutes)

	// Listen
	bindAddr := app.siteCfg.Server.BindAddress + ":" + fmt.Sprint(app.siteCfg.Server.Port)
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
