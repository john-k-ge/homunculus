package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/john-k-ge/homunculus"
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
)

func pullCFEnv() h.HConfig {
	cfUser := os.Getenv("CFU")
	cfPass := os.Getenv("CFP")
	api := os.Getenv("API")
	uaa := os.Getenv("UAA")

	return h.HConfig{
		CFUser:  cfUser,
		CFPass:  cfPass,
		APIHost: api,
		UAAHost: uaa,
	}
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
	homie, err := h.NewHomunculus(pullCFEnv())
	if err != nil {
		log.Printf("Failed to initialize Homunculus: %v", err)
		os.Exit(1)
	}

	err = homie.AddBulkConditions(conditions)
	if err != nil {
		log.Printf("Failed to add bulk conditions: %v", err)
		os.Exit(2)
	}

	fmt.Println("Starting...")
	port := os.Getenv("PORT")
	log.Printf("Listening on port %v", port)
	if len(port) == 0 {
		port = "9000"
	}

	http.HandleFunc("/db", generateDbError)
	http.HandleFunc("timeout", generateTimeout)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("Listen and serve err: %v", err)
	}
}
