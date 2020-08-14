package main

import (
  "./controllers"
  "github.com/gin-gonic/gin"
  "github.com/itsjamie/gin-cors"
  "time"
  "os"
)

func main() {
  GOLANG_PORT := os.Getenv("GOLANG_PORT")
  r := gin.Default()

  r.Use(cors.Middleware(cors.Config{
    Origins:        "*",
    Methods:        "GET, PUT, POST, DELETE",
    RequestHeaders: "Origin, Authorization, Content-Type",
    ExposedHeaders: "",
    MaxAge: 50 * time.Second,
    Credentials: true,
    ValidateHeaders: false,
  }))

	r.GET("/", controllers.GetIndex)
	r.Run(":"+GOLANG_PORT)
}