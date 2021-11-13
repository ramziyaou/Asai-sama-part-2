// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os/signal"
// 	"rest-go-demo/controllers"
// 	"rest-go-demo/database"
// 	"rest-go-demo/entity"
// 	"rest-go-demo/middleware"
// 	"syscall"

// 	"github.com/gorilla/mux"
// 	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
// )

// // README.md хранит в себе все инструкции по эксплуатации программы
// // Используется БД под названием "test_two", которая создавалась по команде create database test_two;
// // Создаются три таблицы: по одной для юзеров, кошельков и для статуса майнинга (последнее во избежание дата рейса). К сожалению, связать таблицы не удалось,
// // в связи с проблемами с ассоциационными функциями в gorm. Если успею, постараюсь исправить этот недочет
// // GetWallet возвращает обновленные данные о балансе кошелька только по завершении майнинга, во время майнинга данные не обновляются во избежание дата рейса.
// // Также постмайнинговый баланс можно увидеть с небольшой задержкой во времени, так как майнинг - процесс непростой :)
// // ID 0 запрещен в пределах данного сервиса

// func main() {
// 	initDB()
// 	user := entity.User{
// 		ID:       1,
// 		Username: "Ramziya",
// 		Password: "1234",
// 	}
// 	if err := database.Connector.Where("username = ?", user.Username).First(&user).Error; err != nil {
// 		database.Connector.Create(&user)
// 	}
// 	// Set stop and start status of mining to false in case of incorrect exit from program earlier (e.g. program exited without stopping mining previously)
// 	// var ss entity.StartStopCheck
// 	// if err := database.Connector.Model(&ss).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	log.Println("Starting the HTTP server on port 8090")

// 	router := mux.NewRouter().StrictSlash(true)
// 	initaliseHandlers(router)
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGINT)
// 	go func() {
// 		<-c
// 		close(c)
// 		var ss entity.StartStopCheck
// 		if err := database.Connector.Model(&ss).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Println("Oops! I received Ctrl+C signal!")
// 		os.Exit(1)
// 	}()
// 	if err := http.ListenAndServe(":8080", router); err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for {
// 		// infinity loop
// 	}
// }

// func initaliseHandlers(router *mux.Router) {
// 	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.GetUser).Methods("GET")
// 	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.SaveUser).Methods("POST")
// 	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.GetWallet).Methods("GET")
// 	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.SaveWallet).Methods("POST")
// 	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/start", controllers.StartMining).Methods("OPTIONS")
// 	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/stop", controllers.StopMining).Methods("OPTIONS")
// 	router.Use(middleware.TimerMiddleware, middleware.HTTPMethodsCheckMiddleware, middleware.AuthMiddleware)
// }

// func initDB() {
// 	config :=
// 		database.Config{
// 			ServerName: "localhost:3306", // Change to your root localhost
// 			User:       "root",
// 			Password:   "admin", // Change to your root password
// 			DB:         "test_two",
// 		}

// 	connectionString := database.GetConnectionString(config)
// 	err := database.Connect(connectionString)
// 	if err != nil {
// 		log.Println("initDB:", err)
// 		return
// 	}
// 	database.Migrate(&entity.User{}, &entity.CryptoWallet{}, &entity.StartStopCheck{})
// }

package main

import (
	"fmt"

	"os"
	"log"
	"net/http"
	"rest-go-demo/controllers"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"rest-go-demo/middleware"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
)

// README.md хранит в себе все инструкции по эксплуатации программы
// Используется БД под названием "test_two", которая создавалась по команде create database test_two;
// Создаются три таблицы: по одной для юзеров, кошельков и для статуса майнинга (последнее во избежание дата рейса). К сожалению, связать таблицы не удалось,
// в связи с проблемами с ассоциационными функциями в gorm. Если успею, постараюсь исправить этот недочет
// GetWallet возвращает обновленные данные о балансе кошелька только по завершении майнинга, во время майнинга данные не обновляются во избежание дата рейса.
// Также постмайнинговый баланс можно увидеть с небольшой задержкой во времени, так как майнинг - процесс непростой :)
// ID 0 запрещен в пределах данного сервиса


func main() {
	initDB()
	user := entity.User{
		ID:       1,
		Username: "Ramziya",
		Password: "1234",
	}
	if err := database.Connector.Where("username = ?", user.Username).First(&user).Error; err != nil {
		database.Connector.Create(&user)
	}
	// Set stop and start status of mining to false in case of incorrect exit from program earlier (e.g. program exited without stopping mining previously)
	var ss entity.StartStopCheck
	if err := database.Connector.Model(&ss).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Starting the HTTP server on port 8090")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlers(router)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println(err)
		return
	}
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.SaveUser).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.GetWallet).Methods("GET")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.SaveWallet).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/start", controllers.StartMining).Methods("OPTIONS")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/stop", controllers.StopMining).Methods("OPTIONS")
	router.Use(middleware.TimerMiddleware, middleware.HTTPMethodsCheckMiddleware, middleware.AuthMiddleware)
}

func initDB() {
	user := os.Getenv("MYSQL_USER")
    pass := os.Getenv("MYSQL_PASSWORD")
    host := os.Getenv("MYSQL_HOST") 
    dbname := os.Getenv("MYSQL_DATABASE")
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA user", user, "pass", pass, "host", host, "dbname", dbname)
	config :=
		database.Config{
			ServerName: host, 
			User:       user,  // "root"
			Password:   pass,  // "admin"
			DB:         dbname,
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		log.Println("initDB:", err)
		return
	}
	database.Migrate(&entity.User{}, &entity.CryptoWallet{}, &entity.StartStopCheck{})
}
