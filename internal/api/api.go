package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"cballProject/path/postgress"

	"github.com/gorilla/mux"
)

type Api struct {
	db *postgress.Storage
	r  *mux.Router
}

func New(db *postgress.Storage) *Api {
	api := Api{db: db, r: mux.NewRouter()}
	api.endpoints()
	return &api
}

func (api *Api) Router() *mux.Router {
	return api.r
}

func (api *Api) endpoints() {
	// получить все данные пользователя по id
	api.r.HandleFunc("/users/{id}", api.GetUserHandler).Methods(http.MethodGet, http.MethodOptions)
	//добавить пользователя
	api.r.HandleFunc("/adduser/", api.GetUserHandler).Methods(http.MethodGet, http.MethodOptions)

}

func (api *Api) GetNewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	s := mux.Vars(r)["id"]
	userId, _ := strconv.Atoi(s)

	user, err := api.db.GetUser(userId)
	if err != nil {
		log.Printf("API get articles error - %v", err)
	}
	json.NewEncoder(w).Encode(user)
}

func (api *Api) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Acess-Control-Allow-Methods", http.MethodPost)
	w.Header().Set("Acces-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body before adding new user to DB", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}
	err = api.db.AddUser(user)
	if err != nil {
		log.Printf("Can't add user to database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("the user has been added successfully!"))
}
