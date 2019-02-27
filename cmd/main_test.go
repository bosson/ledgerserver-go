package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"encoding/base64"
	"encoding/json"

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

	m.Run()
}

func readBody(r *bytes.Buffer) (body []byte, err error) {
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

func getState(t *testing.T, handler http.Handler, session string) string {

	resp, err := http.Get("/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var bodyMap map[string]string
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		t.Fatal(err)
	}

	locationValues, err := url.Parse(bodyMap["redirect"])
	if err != nil {
		t.Fatal(err)
	}

	code := locationValues.Query().Get("code")
	t.Logf("code: %s", code)

	return code
}

func TestRegistration(t *testing.T) {
	handler := loadHandler()
	headers, _ := getAuthHeaders(handler)
	time.Sleep(time.Second)

	reg := registerClient(t, handler, headers)

	t.Logf("Client: %#v", reg)
}
