package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/newrelic/go-agent"
)

func main() {
	cfg := newrelic.NewConfig("my-golang-app", mustGetEnv("NEWRELIC_LICENSE_KEY"))
	cfg.Logger = newrelic.NewDebugLogger(os.Stdout)
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for the application to connect.
	if err := app.WaitForConnection(5 * time.Second); nil != err {
		panic(err)
	}

	// Do the tasks at hand.  Perhaps record them using transactions and/or
	// custom events.
	tasks := []string{"white", "black", "red", "blue", "green", "yellow"}
	for _, task := range tasks {
		txn := app.StartTransaction("task", nil, nil)
		if err := txn.NoticeError(fmt.Errorf("fake error")); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Millisecond)
		if err := txn.End(); err != nil {
			panic(err)
		}
		if err := app.RecordCustomEvent("task", map[string]interface{}{
			"color": task,
		}); err != nil {
			panic(err)
		}
		for i := 0; i <= 1000; i++ {
			if err := app.RecordCustomMetric("my_metric", float64(rand.Intn(1000))); err != nil {
				log.Printf("%+v", err)
			}
		}
	}
	log.Println("ok")
	// Shut down the application to flush data to New Relic.
	app.Shutdown(10 * time.Second)
	log.Println("good")
}

func mustGetEnv(key string) string {
	if val := os.Getenv(key); "" != val {
		return val
	}
	panic(fmt.Sprintf("environment variable %s unset", key))
}
