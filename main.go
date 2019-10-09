package main

import (
	"log"
	"os"
	"time"

	nr "github.com/newrelic/go-agent"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	license := os.Getenv("NEWRELIC_LICENSE_KEY")
	log.Println(license)
	config := nr.NewConfig("my-golang-app", license)
	app, err := nr.NewApplication(config)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := app.RecordCustomEvent("my_type", map[string]interface{}{"hello": "world"}); err != nil {
		log.Fatalf("%+v", err)
	}
	log.Println("ok")
	time.Sleep(time.Hour)
}
