package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/john-k-ge/homunculus"
	"github.com/john-k-ge/homunculus/cf"
)

const (
	dbError = "db-error"
	timeout = "connection-timeout"
)

var (
	homie      *h.Homunculus
	conditions = map[string]int64{
		dbError: 5,
		timeout: 3,
	}
	cfConf = cf.CfEnv{
		Inx:     0,
		Guid:    "60b31aa7-6ee8-42a7-97b6-bfb514b42f04",
		CFUser:  os.Getenv("uid"),
		CFPass:  os.Getenv("pass"),
		APIHost: "api.system.aws-usw02-pr.ice.predix.io",
		UAAHost: "uaa.system.aws-usw02-pr.ice.predix.io",
	}
)

func pullCFEnv() *cf.CfEnv {

	cfUser := os.Getenv("CFU")
	cfPass := os.Getenv("CFP")
	api := os.Getenv("API")
	uaa := os.Getenv("UAA")

	cfEnvVals := cf.GetCFEnvVals()

	cfEnvVals.APIHost = api
	cfEnvVals.UAAHost = uaa
	cfEnvVals.CFUser = cfUser
	cfEnvVals.CFPass = cfPass

	return cfEnvVals
}

func generateDbError(w http.ResponseWriter, req *http.Request) {
	log.Print("Triggering a DB error")
	current, err := homie.Increment(dbError)
	if err != nil {
		log.Printf("Failed to increment %v: %v", dbError, err)
		fmt.Fprintf(w, "Failed to increment `%v`: %v", dbError, err)
		return
	}
	fmt.Fprintf(w, "`%v` is now %v out of %v", dbError, current, conditions[dbError])
}

func generateTimeout(w http.ResponseWriter, req *http.Request) {
	log.Print("Triggering a timeout")
	current, err := homie.Increment(timeout)
	if err != nil {
		log.Printf("Failed to increment %v: %v", timeout, err)
		fmt.Fprintf(w, "Failed to increment `%v`: %v", timeout, err)
		return
	}
	fmt.Fprintf(w, "`%v` is now %v out of %v", timeout, conditions[timeout], current)
}

func main() {
	var err error
	if len(os.Getenv("LOCAL")) != 0 {
		log.Print("Running on a local machine")
		homie, err = h.NewHomunculus(&cfConf)
	} else {
		homie, err = h.NewHomunculus(pullCFEnv())
	}
	if err != nil {
		log.Printf("Failed to initialize Homunculus: %v", err)
		os.Exit(1)
	}

	err = homie.AddBulkConditions(conditions)
	if err != nil {
		log.Printf("Failed to add bulk conditions: %v", err)
		os.Exit(2)
	}

	log.Printf("String represention: %v", homie)

	fmt.Println("Starting...")
	port := os.Getenv("PORT")
	log.Printf("Listening on port %v", port)
	if len(port) == 0 {
		port = "9000"
	}

	http.HandleFunc("/db", generateDbError)
	http.HandleFunc("/timeout", generateTimeout)

	log.Printf("Running on port: %v", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("Listen and serve err: %v", err)
	}
}
