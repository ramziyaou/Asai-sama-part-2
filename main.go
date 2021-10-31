package main

import (
	"log"
	"net/http"
	"rest-go-demo/controllers"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"rest-go-demo/middleware"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
)

func main() {
	initDB()
	user := entity.User{
		ID:       0,
		Username: "Ramziya",
		Password: "1234",
	}
	if err := database.Connector.Where("username = ?", user.Username).First(&user).Error; err != nil {
		database.Connector.Create(&user)
	}

	log.Println("Starting the HTTP server on port 8090")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlers(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.SaveUser).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.GetWallet).Methods("GET")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.SaveWallet).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/start", controllers.StartMining).Methods("OPTIONS")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/stop", controllers.StopMining).Methods("OPTIONS")

	router.HandleFunc("/delete/{id}", controllers.DeleteU).Methods("DELETE")
	router.Use(middleware.TimerMiddleware, middleware.HTTPMethodsCheckMiddleware, middleware.AuthenticationMiddlewareAuthMiddleware)
}

func initDB() {
	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "root",
			Password:   "admin",
			DB:         "test_two",
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}
	database.Migrate(&entity.User{}, &entity.CryptoWallet{}, &entity.StartStopCheck{})
	// database.Migrate(&entity.CryptoWallet{})
}
