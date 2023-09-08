package main

import (
	"log"

	"github.com/MisterKirill/clover/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main () {
	err := godotenv.Load()
	if err != nil {
	  	log.Fatal("Failed to  load .env file")
	}

	r := gin.Default()
	r.MaxMultipartMemory = 5 << 20

	r.Use(cors.Default())

	r.POST("/avatars", routes.CreateAvatar)
	r.GET("/avatars/:avatarID", routes.GetAvatar)

	log.Println("Starting webserver")

	err = r.Run()
	if err != nil {
		log.Fatal("Failed to start webserver: ", err)
	}
}
