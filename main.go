package main

import (
	"log"
	"net/http"

	"github.com/MisterKirill/clover/routes"
	"github.com/joho/godotenv"
)

func main () {
	err := godotenv.Load()
	if err != nil {
	  	log.Fatal("Failed to  load .env file")
	}

	http.HandleFunc("/avatars", routes.CreateAvatar)

	log.Println("Starting webserver")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start webserver: ", err)
	}
}
