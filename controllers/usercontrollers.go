//https://gorm.io/docs/transactions.html
//TRANSACTION!!!!!!!!!

package controllers

import (
	// "encoding/json"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	InvalidInput      error = errors.New("Invalid username or password")
	DuplicateID       error = errors.New("ID exists")
	DuplicateUsername error = errors.New("Username exists")
	Commit            error = errors.New("Transaction's been committed")
)

func NewUser(strId, username, password string) (entity.User, error) {
	var user entity.User

	if username == "" || password == "" {
		return user, InvalidInput
	}

	// Convert ID into int
	id, err := strconv.Atoi(strId)
	if err != nil {
		return user, err
	}

	// Check if ID exists
	if err := database.Connector.Where("id = ?", id).First(&user).Error; err == nil {
		return user, DuplicateID
	}

	// Check if username exists
	if err := database.Connector.Where("username = ?", username).First(&user).Error; err == nil {
		return user, DuplicateUsername
	}

	// If all checks above passed
	return entity.User{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}

//GetAllPerson get all person data
// func GetAllPerson(w http.ResponseWriter, r *http.Request) {
// 	var users []entity.Person
// 	database.Connector.Find(&users)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(users)
// }

//GetPersonByID returns person with specific ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from path
	vars := mux.Vars(r)
	strId := vars["id"]

	// Incorrect ID format
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is not integer")
		return
	}

	// Check if user exists
	var user entity.User
	if err := database.Connector.Where("id = ?", id).First(&user).Error; err != nil {
		// if err := database.Connector.First(&user, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with ID %d not found  %s\n", id, err)
		GetAllPerson(w, r)
		return
	}

	// Display user record if found
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// TODO: how to print wallet names only? maybe encode into json only certain fields and then encode wallet names only? or display via fprintf
	// or maybe create a separate func for printing data for user, includig wallet names
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// fmt.Fprintf(w, "ID: %d, username: %s\nWallets:\n", u.ID, u.Username)
	// if len(u.Wallets) == 0 {
	// 	fmt.Fprintln(w, "empty")
	// }
	// for _, v := range u.Wallets {
	// 	fmt.Fprintln(w, v.Name)
	// }
}

func GetAllPerson(w http.ResponseWriter, r *http.Request) {
	var persons []entity.User
	database.Connector.Find(&persons)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(persons)
}

func GetAllWallet(w http.ResponseWriter, r *http.Request) {
	var wallets []entity.CryptoWallet
	database.Connector.Find(&wallets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallets)
	fmt.Println(wallets)
}

//CreatePerson creates person
func SaveUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from path
	vars := mux.Vars(r)
	strId := vars["id"]

	// Get username and password from params
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := NewUser(strId, username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	// If no error
	database.Connector.Create(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	//   Or

	//   if result := db.Where("name = ?", "jinzhu").First(&user); result.Error != nil {
	// 	// error handling...	// Store user to middleware map to check upon authentication
	// s.Amw.Populate(username, password)

	//   GORM returns ErrRecordNotFound when failed to find data with First, Last, Take, if there are several errors happened, you can check the ErrRecordNotFound error with errors.Is, for example:

	//   // Check if returns RecordNotFound error
	//   err := db.First(&user, 100).Error
	//   errors.Is(err, gorm.ErrRecordNotFound)

	fmt.Fprintln(w, "Success")
}

// Get wallet amount
func GetWallet(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	// Get username from value passed in request context upon authentication
	username := r.Context().Value("username").(string)
	database.Connector.Where("username = ?", username).First(&user)

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]
	// Check if user has any wallets
	if user.Wallets == "" {
		// if user.Wallets == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "You have no cryptowallets")
		return
	}

	// Check if user has given wallet
	for _, name := range strings.Split(user.Wallets, " ") {
		if name == wallet {
			// Get wallet info if everything above ok
			var v entity.CryptoWallet
			if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
				fmt.Println("wallet error!", err)
			}
			// fmt.Println(&v)
			v.RLock()
			defer v.RUnlock()
			GetAllWallet(w, r)
			fmt.Fprintf(w, "Wallet %s, amount: %d", wallet, v.Amount)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&v)
			fmt.Println(v.Amount)
			return
		}
	}
	// if not ok
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Couldn't find wallet %s under your username", wallet)
	return

	// To save wallet:
	// user := entity.User{Username: *username}
	// Get wallet name from path
	// vars := mux.Vars(r)
	// wallet := vars["name"]
	// user.Wallets = append(user.Wallets, entity.NewWallet(wallet))
	// database.Connector.Save(&user)

	// db.Model(&user).Where("foobar = ?", "abc").Related(&orders)
	// https://blog.davidvassallo.me/2019/04/08/lessons-learned-golang-gorm-filtering-associations/
	// database.Connector.Model(&user).Where("name = ?", wallet).Related(&Wallets) // ili wallets i sozdat' do etogo wallets variable????
}

// Create new wallet
func SaveWallet(w http.ResponseWriter, r *http.Request) {
	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]
	var user entity.User
	// Get username from value passed in request context upon authentication
	username := r.Context().Value("username").(string)
	if err := database.Connector.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println("ERRRRRRRRRRRRRRRRR:", err)
		return
	}
	// database.Connector.Where("username = ?", username).First(&user)
	// Check if wallet already exists
	for _, name := range strings.Split(user.Wallets, " ") {
		if name == wallet {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Wallet under name %s already exists", wallet)
			return
		}
	}

	// // Check if user has any wallets
	// if user.Wallets == nil {
	// 	user.Wallets = []*entity.CryptoWallet{}
	// }

	// Create wallet under given name if everything above ok
	// user.Wallets = append(user.Wallets, *entity.NewWallet(wallet))
	// user.Wallets = append(user.Wallets, wallet)
	user.Wallets += wallet + " "
	// for _, v := range user.Wallets {
	// 	fmt.Println(v)
	// }
	// Save updated changes
	if err := database.Connector.Save(&user).Error; err != nil {
		fmt.Println("ERR@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@2: ", err)
	}
	v := entity.NewWallet(wallet)
	v.Username = username
	fmt.Println("amount when creating:", v.Amount)
	if err := database.Connector.Create(&v).Error; err != nil {
		fmt.Println("ERRииииииииииии3: ", err)
	}
	ss := entity.NewStartStop(username, wallet)
	if err := database.Connector.Create(&ss).Error; err != nil {
		fmt.Println("ERRииииииииииии3: ", err)
	}
	GetAllWallet(w, r)
	GetAllPerson(w, r)
	fmt.Fprintln(w, "Success")
}

// Start cryptomining
func StartMining(w http.ResponseWriter, r *http.Request) {
	// Get username from request context
	if err := dbTransaction(database.Connector, w, r); err != nil && err != Commit {
		fmt.Println(err)
	} else if err == Commit || err == nil {
		fmt.Println("Transaction completed")
	}
}

func dbTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Get username from request context
	username := r.Context().Value("username")

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]

	// Get user from db

	// ok = true
	var v entity.CryptoWallet
	if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		// tx.Error = fmt.Errorf("Couldn't find wallet %s under your username", wallet)
		return err
	}
	// var result []bool
	// database.Connector.Select([]string{"notstarted", "stop"}).Where("username = ?", username).Where("name = ?", wallet).Find(&crypto_wallets)
	// database.Connector.Table("crypto_wallet").Preload("Orders").Find(&APIUser{}) preloads associations with given conditions
	// if err := database.Connector.Table("crypto_wallets").Select([]string{"stop"}).Where("username = ?", username).Where("name = ?", wallet).Find(&result).Error; err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	// db.Select([]string{"name", "age"}).Find(&users)
	// SELECT name, age FROM users;
	// SELECT name, age FROM users;
	// database.Connector.Raw("SELECT amount FROM crypto_wallets WHERE name = ? and username = ? ", wallet, username).Scan(&result)
	// the above code returns default values! bummer :)
	// fmt.Println(result)
	// if result[0] == false {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "Mining already started\n")
	// 	// tx.Error = fmt.Errorf("Mining already started\n")
	// 	return nil // ???
	// }
	// fmt.Println(v)
	// db.Where("name = ?", "jinzhu").Where("age = ?", 18).First(&user)

	// if err := database.Connector.Select("username", "name").Where(" = ?", 18}).Find(&User{})
	// if err := database.Connector.Where("username = ?", username).First(&v).Error; err != nil {
	// 	tx.Rollback()
	// }

	// Start mining
	// v.RLock()
	// defer v.RUnlock()
	// if !v.Notstarted {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	tx.Error = fmt.Errorf("Mining already started\n")
	// 	tx.Rollback()
	// }
	var ss entity.StartStopCheck
	database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss) // check for error!
	if ss.Start {
		w.WriteHeader(http.StatusBadRequest)
		// tx.Error = fmt.Errorf("Mining already started\n")
		return fmt.Errorf("Mining not stopped yet\n") //?????
	}
	// v.Lock()
	// defer v.Unlock()
	// v.Notstarted = false
	// fmt.Println(v.Notstarted, v.Amount)
	// db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"notstarted": false}).Error; err != nil {
	// if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"notstarted": false}).Error; err != nil {
	ss.Start = true
	if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"start": true}).Error; err != nil {
		// if err := database.Connector.Save(&ssp).Error; err != nil {
		fmt.Println("wallet save error", err) // nil
		return err
	}
	fmt.Println("starting mining info", &v)
	// v.Start = true
	// database.Connector.Save(&v)
	// fmt.Println(v.Name, v.Start)
	v.RLock()
	defer v.RUnlock()
	fmt.Fprintf(w, "Starting mining, current amount: %d\n", v.Amount)
	// log.Printf("Starting mining, current amount: %d\n", v.Amount)

	// v.RUnlock()

	// Wait for stop instructions, otherwise keep mining
	go func(v *entity.CryptoWallet) {
		wg := &sync.WaitGroup{}
		// wg2 := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg.Add(1)
		go func(ctx context.Context, v *entity.CryptoWallet, wg *sync.WaitGroup) {
			for {
				// wg2.Add(1)
				v.Mine()
				// wg2.Done()
				// db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
				// if err := database.Connector.Model(&v).Where("name = ?", wallet).Updates(map[string]interface{}{"notstarted": false}).Error; err != nil {
				// if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"amount": v.Amount}).Error; err != nil {
				// 	fmt.Println("wallet amount update error", err) // nil
				// 	return
				// }
				// fmt.Println(v)
				select {
				case <-ctx.Done():
					if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"amount": v.Amount}).Error; err != nil {
						fmt.Println("wallet amount update error", err) // nil
						return
					}
					wg.Done()
					fmt.Println("wgDone")
					return
				default:
				}
			}
		}(ctx, v, wg)
		fmt.Println("1", ss.Stop)
		for !ss.Stop {
			time.Sleep(time.Second * 10)
			// wg2.Wait()
			// if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
			if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss).Error; err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Println("NOW: Error")
				return
			}
			fmt.Println("2", ss.Stop)

			// fmt.Println(v.Stop)
			// select {
			// case _, ok := <-v.Stop:
			// 	if !ok {
			// 		tx.Error = fmt.Errorf("mining failed")
			// 		tx.Rollback()
			// 		return
			// 	}
			// 	cancel()
			// 	v.Lock()         ////
			// 	defer v.Unlock() ////
			// 	v.Notstarted = true
			// 	if err := database.Connector.Model(&v).Updates(map[string]interface{}{"amount": v.Amount}).Error; err != nil {
			// 		fmt.Println("wallet amount update error", err) // nil
			// 		return
			// 	}
			// 	return
			// }
		}
		cancel()
		fmt.Println("3", ss.Stop)
		fmt.Println("canceling contexts")
		wg.Wait()
		// if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"notstarted": true, "stop": false}).Error; err != nil {
		// ss.Start, ss.Stop = false, false
		if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
			// if err := database.Connector.Save(&ssp).Error; err != nil {
			fmt.Println("wallet status update error", err) // nil
			return
		}

		// time.Sleep(time.Second)
		return
	}(&v)
	if tx.Error != Commit {
		// tx.Error = Commit
		tx.Commit().Error = Commit
		return Commit
	} else {
		tx.Rollback()
	}
	return nil
}

// // Start cryptomining
// func StartMining(w http.ResponseWriter, r *http.Request) {
// 	// Get username from request context
// 	username := r.Context().Value("username")

// 	// Get wallet name from path
// 	vars := mux.Vars(r)
// 	wallet := vars["name"]

// 	// Get user from db
// 	var user entity.User
// 	database.Connector.Where("username = ?", username).First(&user)
// 	ok := false
// 	for _, name := range strings.Split(user.Wallets, " ") {
// 		if name == wallet {
// 			ok = true
// 		}
// 	}
// 	// Check if wallet exists
// 	for _, name := range strings.Split(user.Wallets, " ") {
// 		if name == wallet {
// 			// Get wallet from DB
// 			var v entity.CryptoWallet
// 			if err := database.Connector.Where("username = ?", username).First(&v).Error; err != nil {
// 				fmt.Println("Wallet error!", err)
// 			}
// 			fmt.Println(v.Name, v.Start)
// 			// Start mining
// 			v.RLock()
// 			defer v.RUnlock()
// 			if v.Start {
// 				w.WriteHeader(http.StatusBadRequest)
// 				fmt.Fprintf(w, "Mining already started\n")
// 				return
// 			}
// 			// v.Lock()
// 			// defer v.Unlock()
// 			v.Start = true
// 			database.Connector.Save(&user)
// 			fmt.Println(v.Name, v.Start)
// 			v.RLock()
// 			defer v.RUnlock()
// 			fmt.Fprintf(w, "Starting mining, current amount: %d\n", v.Amount)
// 			// log.Printf("Starting mining, current amount: %d\n", v.Amount)

// 			// v.RUnlock()

// 			// Wait for stop instructions, otherwise keep mining
// 			go func(v *entity.CryptoWallet) {
// 				ctx, cancel := context.WithCancel(context.Background())
// 				defer cancel()
// 				go func(ctx context.Context, v *entity.CryptoWallet) {
// 					for {
// 						v.Mine()
// 						database.Connector.Save(&user)
// 					}
// 				}(ctx, v)
// 				for {
// 					select {
// 					case _, ok := <-v.Stop:
// 						if !ok {
// 							fmt.Fprintf(w, "mining failed")
// 							database.Connector.Save(&user)
// 							return
// 						}
// 						cancel()
// 						v.Lock()         ////
// 						defer v.Unlock() ////
// 						v.Start = false
// 						database.Connector.Save(&user)
// 						return
// 					}
// 				}
// 			}(&v)
// 		}
// 	}

// 	// If failed to find wallet
// 	if !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		// WHY BEING PRINTED????
// 		fmt.Fprintf(w, "Couldn't find wallet %s under your username", wallet)
// 		return
// 	}
// }

func StopMining(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(time.Second * 5)
	// Get username from request context
	username := r.Context().Value("username")

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]

	// var v entity.CryptoWallet

	// if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
	// 	// If failed to find wallet
	// 	w.WriteHeader(http.StatusNotFound)
	// 	fmt.Fprintf(w, "Couldn't find wallet %s under your username", wallet)
	// 	return
	// }

	// Check if mining's started
	// v.RLock()
	// defer v.RUnlock()

	var ss entity.StartStopCheck
	database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss) // check for error!
	if !ss.Start {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Can't stop mining that hasn't started\n")
		return
	}
	// fmt.Println("all info on v I wanna stop:", &v)
	// if v.Notstarted {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "Can't stop mining that hasn't started\n")
	// 	return
	// }

	// Send signal to StartMining method to stop mining
	log.Println("closing wallet", wallet)
	// v.Stop <- struct{}{}
	// v.Stop = false
	// v.RLock()
	// defer v.RUnlock()

	// ss.Stop = true

	if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"stop": true}).Error; err != nil {
		// if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"stop": true}).Error; err != nil {
		fmt.Println("wallet amount update error", err) // nil
		return
	}
	time.Sleep(time.Second * 10)
	fmt.Fprintf(w, "Mining stopped for wallet %s\n", wallet)
}

// //
// // // //DeletPersonByID delete's person with specific ID
// // // func DeletPersonByID(w http.ResponseWriter, r *http.Request) {
// // // 	vars := mux.Vars(r)
// // // 	key := vars["id"]

// // // 	var person entity.Person
// // // 	id, _ := strconv.ParseInt(key, 10, 64)
// // // 	database.Connector.Where("id = ?", id).Delete(&person)
// // // 	w.WriteHeader(http.StatusNoContent)
// // // }
// //

func DeleteU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var person entity.User
	id, _ := strconv.ParseInt(key, 10, 64)
	database.Connector.Where("id = ?", id).Delete(&person)
	w.WriteHeader(http.StatusNoContent)
}
