package main

import (
	"net/http"
)

func main() {

	port := "5000"
	m := MakeHandler()
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}
