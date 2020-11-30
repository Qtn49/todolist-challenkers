package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var db, _ = gorm.Open("mysql", "root:@/todolist_challenkers?charset=utf8&parseTime=True&loc=Local")

/**

Structure définissant le modèle d'une tâche

*/
type TodoItem struct {
	Id          int `gorm:"primary_key"`
	Titre       string
	Nom         string
	Description string
	Etat        string `gorm:"default:'A Faire'"`
	DateRendu   time.Time
}

/**

Fonction permettant de créer une tâche accessible par la méthode POST

*/
func CreateItem(w http.ResponseWriter, r *http.Request) {

	titre, nom, description := r.FormValue("titre"), r.FormValue("nom"),
		r.FormValue("description")
	dateString := r.FormValue("date_rendu")
	dates := strings.Split(dateString, "-")
	year, _ := strconv.Atoi(dates[0])
	month, _ := strconv.Atoi(dates[1])
	day, _ := strconv.Atoi(dates[2])
	dateRendu := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	todo := &TodoItem{Titre: titre, Nom: nom, Description: description, Etat: "A Faire",
		DateRendu: dateRendu}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)

}

/**

Fonction permettant de vérifier l'existence d'une tâche à partir d'un identifiant

*/
func GetItemId(id int) bool {

	todo := &TodoItem{}
	result := db.First(&todo, id)

	if result.Error != nil {

		log.Warn("Item not Found")
		return false

	}

	return true

}

/**

Fonction permettant de gérer la suppression d'une tâche dans la base de donnée, appelée par la méthode DELETE

*/
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

/**

Mise à jour de l'État d'une tâche à partir de son identifiant récupéré par la méthode POST

*/
func UpdateItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if !GetItemId(id) {

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false}, "error": "Item not found"`)

	} else {

		etat := r.FormValue("etat")
		if ValideState(etat) {
			todo := &TodoItem{}
			db.First(&todo, id)
			todo.Etat = etat
			db.Save(&todo)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"updated": true}`)
		}

	}

}

/**

Renvoie au format json les tâches ratés correspondant au critere envoyé

*/
func SearchMissed(w http.ResponseWriter, r *http.Request) {

	var todos []TodoItem
	critere := r.FormValue("critere")
	TodoItemsMissed := db.Where("nom LIKE ? AND date_rendu < ? AND Etat != 'Fait'", "%"+critere+"%", time.Now()).Find(&todos).Value
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TodoItemsMissed)

}

/**

Renvoie les tâches non manquées correspondant au critere demandé

*/
func SearchNotMissed(w http.ResponseWriter, r *http.Request) {

	var todos []TodoItem
	critere := r.FormValue("critere")
	TodoItemsNotMissed := db.Where("nom LIKE ? AND date_rendu > ? OR nom LIKE ? AND Etat = 'Fait'", "%"+critere+"%", time.Now(), "%"+critere+"%").Find(&todos).Value
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TodoItemsNotMissed)

}

/**

appelée par la méthode GET, méthode qui renvoie les États possible d'une tâche

*/
func GetValuesState(w http.ResponseWriter, r *http.Request) {

	values := []string{"A Faire", "En Cours", "Fait"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)

}

/**

Vérifie que l'état renvoyé correspond à ceux disponibles

*/
func ValideState(etat string) bool {

	values := []string{"A Faire", "En Cours", "Fait"}

	for _, v := range values {
		if etat == v {

			return true

		}
	}

	return false

}

/**

fonction d'initialisation permettant d'initialiser le log et son format

*/
func init() {

	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)

}

/**

Fonction main permettant de définir les différentes routes et leur fonction correspondantes

*/
func main() {

	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItem{})
	db.Debug().AutoMigrate(&TodoItem{})

	todo := &TodoItem{
		Titre:       "Test",
		Nom:         "Effectuer un test",
		Description: "Je remplis les cases pour le test",
		DateRendu:   time.Date(2020, 12, 23, 22, 10, 0, 0, time.Local),
	}
	db.Create(&todo)
	todo2 := TodoItem{
		Titre:       "Bruh",
		Nom:         "Sheesh",
		Description: "bibi",
		Etat:        "En Cours",
		DateRendu:   time.Date(2020, 2, 21, 22, 10, 0, 0, time.Local),
	}
	db.Create(&todo2)
	todo3 := TodoItem{
		Titre:       "Bruh",
		Nom:         "Sheesh",
		Description: "bibi",
		Etat:        "Fait",
		DateRendu:   time.Date(2020, 2, 21, 22, 10, 0, 0, time.Local),
	}
	db.Create(&todo3)

	router := mux.NewRouter()

	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
	router.HandleFunc("/searchMissed", SearchMissed).Methods("POST")
	router.HandleFunc("/searchNotMissed", SearchNotMissed).Methods("POST")
	router.HandleFunc("/states", GetValuesState).Methods("GET")

	log.Info("Starting TodoList")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)

}
