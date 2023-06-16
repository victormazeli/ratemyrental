package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/database"
	"rateMyRentalBackend/http/controllers"
	"rateMyRentalBackend/http/middlewares"
	"rateMyRentalBackend/http/response"
	"rateMyRentalBackend/http/routes/misc"
	routesV1 "rateMyRentalBackend/http/routes/v1"
	"rateMyRentalBackend/http/ws"
	"rateMyRentalBackend/models"
)

func main() {
	environment := flag.String("e", "development", "")

	//flag.Parse()
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
		&models.ChatMessage{},
		&models.Channel{},
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
		response.ErrorResponse(http.StatusNotFound, "404 page not found", c)
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
