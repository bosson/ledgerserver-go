package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bosson/ledgerserver/pkg/api"
	"github.com/goadapp/goad/version"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

func main() {

	_ = godotenv.Load(".env")
	viper.AutomaticEnv()
	viper.SetDefault("log_formatter", "text")
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("prod", false)
	viper.SetDefault("addr", ":9000")
	viper.SetDefault("base", "http://localhost:6502")
	viper.SetDefault("cryptoprovider", "file")
	viper.SetDefault("gorm_dialect", "sqlite3")
	viper.SetDefault("gorm_options", "file:./ledger.db?cache=shared")
	viper.SetDefault("stats", "none")

	// viper.SetDefault("logstash", "")
	// if viper.GetString("logstash") != "" {
	// 	conn, err := gas.Dial("tcp", viper.GetString("logstash"))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "app"}))
	// }

	log.SetOutput(os.Stdout)

	switch viper.GetString("log_formatter") {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
		break
	case "text":
		log.SetFormatter(
			&log.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			})
	}

	lvl, err := log.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Fatal("unknodn log level: " + viper.GetString("log_level"))
	}

	log.SetLevel(lvl)

	log.WithFields(log.Fields{
		"service": "ledgerserver",
		"version": version.Version,
	})

	// stats.FromEnv("ledgerserver")
	// w := log.Writer()
	// defer w.Close()

	addr := viper.GetString("addr")

	server := &graceful.Server{
		Timeout: time.Duration(15) * time.Second,
		Server: &http.Server{
			Addr:        addr,
			Handler:     loadHandler(),
			ReadTimeout: time.Duration(10) * time.Second,
			// ErrorLog:    log.Logger,
		},
	}

	log.Info("server starting at " + addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

var w *io.PipeWriter
var db *gorm.DB

func loadHandler() http.Handler {

	// prod := viper.GetBool("prod")
	r := httprouter.New()

	var err error
	if db == nil {
		db, err = gorm.Open(viper.GetString("gorm_dialect"), viper.GetString("gorm_options"))
		if err != nil {
			log.Fatal(err)
		}
		if viper.GetBool("gorm_debug") {
			db.LogMode(true)
		}
		// db.SetLogger(log.New(w, "database", 0))
		db.DB().SetMaxIdleConns(1)
		db.DB().SetMaxOpenConns(5)
	}

	r.GET("/", api.Version)
	r.GET("/health", api.Version)
	r.GET("/web/*filepath", api.Static)

	ledger := api.NewLedgerPoster(viper.GetString("base"))
	r.POST("/service/ledger/:book", ledger.LedgerPost)

	return r
}
