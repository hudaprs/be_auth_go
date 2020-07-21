package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"be_auth_go/app/middlewares"
	"be_auth_go/app/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (app *App) Routes() {
	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router.Use(middlewares.SetContentTypeJSON)

	app.Router.HandleFunc("/api/users", app.CreateUser).Methods("POST")
	app.Router.HandleFunc("/api/auth", app.Login).Methods("POST")
}

func (app *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	app.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Printf("\n Cannot connect to database %s", DbName)
		log.Fatal("This is the error:", err)
	}

	fmt.Println("Database Connected...")
	app.DB.Debug().AutoMigrate(&models.User{})
	app.Routes()
}

func (app *App) RunServer() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("There's something wrong loading .env")
		return
	}

	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	fmt.Println("Server started at port 8000")
	log.Fatal(http.ListenAndServe(":8000", app.Router))
}
