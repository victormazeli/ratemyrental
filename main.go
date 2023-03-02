package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rateMyRentalBackend/common"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/controllers"
	"rateMyRentalBackend/database"
	"rateMyRentalBackend/middlewares"
	"rateMyRentalBackend/models"
	routesV1 "rateMyRentalBackend/routes/v1"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
	}
	flag.Parse()
	app := config.App(*environment)
	env := app.Env
	db := database.ConnectDB(env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)
	err := db.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.PropertyImage{},
		&models.PropertyType{},
		&models.PropertyDetachedType{},
		&models.Otp{},
		&models.Rating{},
		&models.App{},
		&models.FavoriteProperty{},
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.ForwardedByClientIP = true
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	router.Static("/static", "./public")

	health := new(controllers.HealthController)

	router.GET("/", health.Status)

	router.NoRoute(func(c *gin.Context) {
		common.ErrorResponse(http.StatusNotFound, "404 page not found", c)
		return
	})

	router.NoMethod(func(c *gin.Context) {
		common.ErrorResponse(http.StatusMethodNotAllowed, "405 method not allowed", c)
		return
	})

	v1 := router.Group("v1")

	routesV1.Setup(env, db, v1)

	log.Printf("The App is running in %v environment\n", *environment)

	router.Run(":" + env.ServerPort)

}
