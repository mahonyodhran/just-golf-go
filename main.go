package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer db.Close()
	InitializeApp()
	log.Println("App starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
