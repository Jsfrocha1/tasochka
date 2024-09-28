package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

var once sync.Once
var instance *App

func GetDBInstance() *App {
	once.Do(func() {
		dsn := "host=localhost user=postgres password=root dbname=task_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to the database")
		}
		instance = &App{DB: db}
	})
	return instance
}

type PersonReq struct {
	Name    string `json:"name"`
	Age     uint   `json:"age"`
	NewAge  uint   `json:"newage"`
	NewName string `json:"newname"`
}

type Person struct {
	gorm.Model
	Id   uint
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func handlReq(app *App) {
	http.HandleFunc("/main", app.mainHandler)
	http.ListenAndServe("127.0.0.1:3333", nil)
}

func (app *App) mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.postPerson(w, r)
	case http.MethodGet:
		app.getPerson(w, r)
	case http.MethodPatch:
		app.patchPerson(w, r)
	case http.MethodDelete:
		app.deletePerson(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) deletePerson(w http.ResponseWriter, r *http.Request) {
	var personreq PersonReq

	err := json.NewDecoder(r.Body).Decode(&personreq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := app.DB.Where("name = ?", personreq.Name).Delete(&Person{}).Error; err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
}

func (app *App) patchPerson(w http.ResponseWriter, r *http.Request) {
	var personreq PersonReq
	var person Person

	err := json.NewDecoder(r.Body).Decode(&personreq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := app.DB.Model(&person).Where("name = ?", personreq.Name).Updates(Person{Name: personreq.NewName, Age: personreq.NewAge}).Error; err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *App) postPerson(w http.ResponseWriter, r *http.Request) {
	var person Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if person.Name == "" || person.Age == 0 {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	if res := app.DB.Where("name= ?", person.Name).Where("age = ?", person.Age).Find(&person).RowsAffected; res != 0 {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	} else {

		result := app.DB.Create(&person)
		if result.Error != nil {
			http.Error(w, "Failed to create person", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func (app *App) getPerson(w http.ResponseWriter, r *http.Request) {
	var persons []Person
	var person Person
	var id uint

	s := r.RequestURI

	if s == "/main" {
		if err := app.DB.Find(&persons).Error; err != nil {
			http.Error(w, "Failed to retrieve records", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(persons)
		return
	}

	_, err := fmt.Sscanf(s, "/main?id=%d", &id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := app.DB.First(&person, id).Error; err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(person)
}

func main() {
	app := GetDBInstance()

	// Автоматическая миграция
	app.DB.AutoMigrate(&Person{})

	handlReq(app)
}
