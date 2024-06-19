package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"commune/config"

	"github.com/go-chi/chi"
	"github.com/robfig/cron/v3"

	"github.com/rs/zerolog"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

type App struct {
	Config   *config.Config
	Router   *chi.Mux
	HTTP     *http.Server
	MatrixDB *MatrixDB
	Cron     *cron.Cron
	Cache    *Cache
	Log      *zerolog.Logger
}

func (c *App) Activate() {
	c.Log.Info().Msg("Started Commune server")

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		if err := c.HTTP.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
			log.Printf("Shutdown by user")
		}
		close(idleConnsClosed)
	}()

	if err := c.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

type StartRequest struct {
	Config    string
	MakeViews bool
}

var CONFIG_FILE string
var MATRIX_CONFIG map[string]interface{}
var PRODUCTION_MODE bool
var AssetFiles map[string]string

func Start(s *StartRequest) {

	CONFIG_FILE = s.Config

	conf, err := config.Read(s.Config)
	if err != nil {
		panic(err)
	}

	PRODUCTION_MODE = conf.Mode == "production"

	mdb, err := NewMatrixDB()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	cron := cron.New()

	cache, err := NewCache(conf)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		//WriteTimeout: 60 * time.Second,
		IdleTimeout: 120 * time.Second,
		Addr:        fmt.Sprintf(`:%d`, conf.App.Port),
		Handler:     router,
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	c := &App{
		MatrixDB: mdb,
		Config:   conf,
		HTTP:     server,
		Router:   router,
		Cron:     cron,
		Cache:    cache,
		Log:      &logger,
	}

	c.Routes()

	// c.Build()

	// go c.Cron.AddFunc("*/15 * * * *", c.RefreshCache)
	// go c.Cron.Start()

	c.Activate()
}
