package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tylerb/graceful"

	"encoding/base64"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {

	_ = godotenv.Load(".env.test")
	viper.AutomaticEnv()

	viper.SetDefault("gorm_debug", true)

	runtime.GOMAXPROCS(1) // this is required when using an inmem sqlite db

	var err error
	db, err = gorm.Open("sqlite3", ".test.db")
	defer os.Remove(".test.db")
	if err != nil {
		panic(err)
	}

	viper.SetDefault("addr", ":9900")
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

	log.Printf("server starting at %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestPing(t *testing.T) {

	resp, err := http.Get("/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// var bodyMap map[string]string
	// err = json.Unmarshal(body, &bodyMap)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	t.Logf("Client: %#v", string(body))

}

func readBody(r io.Reader) (body []byte, err error) {
	body, err = ioutil.ReadAll(io.LimitReader(r, 1048576))
	if err != nil {
		return
	}

	return
}

func basicAuth(username, password string) map[string]string {

	auth := username + ":" + password
	str := base64.StdEncoding.EncodeToString([]byte(auth))

	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + str

	return headers
}
