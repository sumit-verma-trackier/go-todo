package main

import (
	"fmt"
	"net/http"

	"github.com/sumit-verma-trackier/todo/router"
)

func main() {

	router := router.Router()

	fmt.Println("Server started at 4000")

	http.ListenAndServe(":4000", router)

}
