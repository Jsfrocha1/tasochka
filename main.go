package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tasochka/utils"

	"gorm.io/gorm"
)

var AppUrl = "127.0.0.1:3333"

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

func handlReq() {
	http.HandleFunc("/main", mainHandler)
	http.ListenAndServe(AppUrl, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postPerson(w, r)
	case http.MethodGet:
		getPerson(w, r)
	case http.MethodPatch:
		patchPerson(w, r)
	case http.MethodDelete:
		deletePerson(w, r)
	default:
		http.Error(
			w, 
			utils.ErrMessages.MethodNotAllowed, 
			http.StatusMethodNotAllowed,
		)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	var personreq PersonReq

	err := json.NewDecoder(r.Body).Decode(&personreq)
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidRequest, http.StatusBadRequest)
		return
	}

	err = utils.DB.Where("name = ?", personreq.Name).Delete(&Person{}).Error
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidData, http.StatusBadRequest)
		return
	}
}

func patchPerson(w http.ResponseWriter, r *http.Request) {
	var personreq PersonReq
	var person Person

	err := json.NewDecoder(r.Body).Decode(&personreq)
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidRequest, http.StatusBadRequest)
		return
	}
	err = utils.DB.Model(&person).Where("name = ?", personreq.Name).
		Updates(Person{Name: personreq.NewName, Age: personreq.NewAge}).Error
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidData, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func postPerson(w http.ResponseWriter, r *http.Request) {
	var person Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidRequest, http.StatusBadRequest)
		return
	}
	if person.Name == "" || person.Age == 0 {
		http.Error(w, utils.ErrMessages.MissingFields, http.StatusBadRequest)
		return
	}

	res := utils.DB.Where("name= ?", person.Name).
		Where("age = ?", person.Age).Find(&person).RowsAffected

	if res != 0 {
		http.Error(
			w,
			utils.ErrMessages.FailedToCreatePerson,
			http.StatusInternalServerError,
		)
		return
	}
	result := utils.DB.Create(&person)
	if result.Error != nil {
		http.Error(
			w,
			utils.ErrMessages.FailedToCreatePerson,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(person)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	var persons []Person
	var person Person
	var id uint

	url := r.RequestURI

	if url == "/main" {
		err := utils.DB.Find(&persons).Error
		if err != nil {
			http.Error(
				w,
				utils.ErrMessages.FailedToRetrieveRecords,
				http.StatusInternalServerError,
			)
			return
		}

		json.NewEncoder(w).Encode(persons)

		return
	}

	_, err := fmt.Sscanf(url, "/main?id=%d", &id)
	if err != nil {
		http.Error(w, utils.ErrMessages.InvalidIDFormat, http.StatusBadRequest)
		return
	}
	err = utils.DB.First(&person, id).Error
	if err != nil {
		http.Error(w, utils.ErrMessages.RecordNotFound, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(person)
}

func main() {
	utils.GetDBInstance()

	utils.DB.AutoMigrate(&Person{})

	handlReq()
}
