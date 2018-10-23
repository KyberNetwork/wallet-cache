package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{user_address}/log", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userAddr := vars["user_address"]

		fmt.Println(userAddr)
	})

	r.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("log")
	})

	http.ListenAndServe(":8080", r)
}
