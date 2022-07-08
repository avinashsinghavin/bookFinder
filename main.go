package main

import (
	"RestAPI/cors"
	"RestAPI/model"
	"RestAPI/order"
	"RestAPI/providers"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
)

var AllowedOrigins = []string{"*", "localhost:30001"}

func main() {
	router := mux.NewRouter()
	http.Handle("/", router)

	//------- DB connection -------------------
	dbStructure := providers.DbConfig{}
	dbConfig := dbStructure.ConfigureDbConnection()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName, dbConfig.DbTimeOut)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Tidak Konek DB Errornya : %s", err)
	}
	er := db.AutoMigrate(&model.Book{})
	if er != nil {
		log.Fatalf("Error while creating Table: %s", er)
	}
	handler := order.AggregationOrderHandler{
		Router: router,
		DB:     db,
	}
	handler.Init()
	if err := http.ListenAndServe(":8000", cors.HandleCores(AllowedOrigins)(router)); err != nil {
		log.Fatal(err)
	}
}
