package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
)

var db, _ = gorm.Open("mysql", "root:@/todolist_challenkers?charset=utf8&parseTime=True&loc=Local")

type TodoItem struct {
	Id          int `gorm:"primary_key"`
	Nom         string
	Description string
	Etat        string
	DateRendu   time.Time
}

func CreateItem(w http.ResponseWriter, r *http.Request) {

	nom, description := r.FormValue("nom"), r.FormValue("description")
	dateRendu, _ := time.Parse(time.ANSIC, r.FormValue("date_rendu"))
	todo := &TodoItem{Nom: nom, Description: description, Etat: "À Faire", DateRendu: dateRendu}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)

}

func GetItemId(id int) bool {

	todo := &TodoItem{}
	result := db.First(&todo, id)

	if result.Error != nil {

		log.Warn("Item not Found")
		return false

	}

	return true

}

func GetItemFromId(id int) interface{} {

	if !GetItemId(id) {

		log.Warn("TodoItem not found in Database")
		return nil

	}

	todo := &TodoItem{}
	result := db.Where("Id = ?", id).Find(&todo)

	return result.Value

}

func DeleteItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if !GetItemId(id) {

		w.Header().Set("Content-type", "application/json")
		io.WriteString(w, `{"deleted": false}, "error": "Item not found"`)

	} else {

		todo := &TodoItem{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)

	}

}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if !GetItemId(id) {

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false}, "error": "Item not found"`)

	} else {

		etat := r.FormValue("etat")
		if ValideState(etat) {

		}

	}

}

func GetValuesState(w http.ResponseWriter, r *http.Request) {

	values := []string{"À Faire", "En Cours", "Fait"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)

}

func ValideState(etat string) bool {

	values := []string{"À Faire", "En Cours", "Fait"}

	for _, v := range values {
		if etat == v {

			return true

		}
	}

	return false

}

func init() {

	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)

}

func main() {

	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItem{})
	db.Debug().AutoMigrate(&TodoItem{})

	router := mux.NewRouter()

	router.HandleFunc("/todo", CreateItem).Methods("POST")

	log.Info("Starting TodoList")

	http.ListenAndServe(":8000", router)

}
