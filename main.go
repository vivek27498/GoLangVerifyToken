package main

import (
	"fmt"
	"net/http"
	"verifyToken/Controllers"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	verifyHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subpath := vars["subpath"]
		fmt.Println("subpath", subpath)
		Controllers.SetVerifyRes(w, r) //function name should be capital
	}

	router.HandleFunc("/verifyToken/{subpath:.*}", verifyHandler).Methods("POST")

	fmt.Println("running on 9090")
	http.ListenAndServe("127.0.0.1:9090", router)

	select {}

}
