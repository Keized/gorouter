package main

import (
	"./gorouter"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type User struct {
	id			int
	Username 	string
	Email 		string
}

func SearchUser(users []User, id int) (u User, err error) {
	for _, u := range users {
		if id == u.id {
			return u, nil
		}
	}

	return u, errors.New("User not found")
}

var users = []User{{1, "john", "Doe"}, {2, "jane", "doe"}}

func getUsers(w http.ResponseWriter, req *http.Request, params gorouter.Params) {
	userJson, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func getUser(w http.ResponseWriter, req *http.Request, params gorouter.Params) {
	id, _ := strconv.Atoi(params.Get("id"))
	user, err := SearchUser(users, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "User not found"}`))
	} else {
		userJson, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	}
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	r := gorouter.Router{}
	r.GET("/users", getUsers)
	r.GET("/user/:id", getUser)
	err := http.ListenAndServe(":8001", &r)
	if err != nil {
		panic(err)
	}
}