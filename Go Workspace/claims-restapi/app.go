package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func AllClaimsEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func FindClaimEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func CreateClaimEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var claim Claim
	if err := json.NewDecoder(r.Body).Decode(&claim); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	claim.ID = bson.NewObjectId()
	if err := dao.Insert(claim); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, claim)
}

func UpdateClaimEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func DeleteClaimEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/claims", AllClaimsEndPoint).Methods("GET")
	r.HandleFunc("/claims", CreateClaimEndPoint).Methods("POST")
	r.HandleFunc("/claims", UpdateClaimEndPoint).Methods("PUT")
	r.HandleFunc("/claims", DeleteClaimEndPoint).Methods("DELETE")
	r.HandleFunc("/claims/{id}", FindClaimEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
