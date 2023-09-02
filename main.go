package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/database"
	models2 "rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/controllers"
	"rateMyRentalBackend/http/middlewares"
	"rateMyRentalBackend/http/response"
	"rateMyRentalBackend/http/routes/misc"
	routesV1 "rateMyRentalBackend/http/routes/v1"
	"rateMyRentalBackend/http/ws"
)

func main() {
	environment := flag.String("e", "development", "")

	//flag.Parse()
	app := config.App(*environment)
	env := app.Env
	db := database.ConnectDB(env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)
	err := db.AutoMigrate(
		&models2.User{},
		&models2.Property{},
		&models2.PropertyImage{},
		&models2.PropertyType{},
		&models2.PropertyDetachedType{},
		&models2.Otp{},
		&models2.Rating{},
		&models2.App{},
		&models2.FavoriteProperty{},
		&models2.ChatMessage{},
		&models2.Channel{},
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go ws.HubInstance.Run()

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	//router.ForwardedByClientIP = true
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	router.Static("/static", "./public")
	router.LoadHTMLGlob("templates/*.html")

	health := new(controllers.HealthController)

	router.GET("/", health.Status)

	router.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "chat.html", gin.H{})
	})

	router.NoRoute(func(c *gin.Context) {
		response.ErrorResponse(http.StatusNotFound, "404 resource not found", c)
		return
	})

	router.NoMethod(func(c *gin.Context) {
		response.ErrorResponse(http.StatusMethodNotAllowed, "405 method not allowed", c)
		return
	})

	v1 := router.Group("v1")

	routesV1.Setup(env, db, v1)

	globalRoute := router.Group("")

	ws.WebsocketRouter(env, db, globalRoute)

	misc.MiscellaneousRouter(env, db, globalRoute)

	log.Printf("The api is running in %v environment\n", *environment)

	router.Run(":" + env.ServerPort)

}
