package main

import (
	"fmt"
	"hello/countries"
	"net/http"

	set "github.com/deckarep/golang-set"
)

func helloWorld(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello World \n")
}

func main() {
	fmt.Println("Hello World!!!")
	fmt.Println("Hello", countries.GetCountry("FR"))

	yellowFruit := set.NewSet()
	yellowFruit.Add("banana")
	yellowFruit.Add("lemon")
	yellowFruit.Add("pineapple")
	fmt.Println(yellowFruit)
	redFruit := set.NewSetFromSlice(
		[]interface{}{"apple", "cherry", "strawberry"})
	fmt.Println(redFruit)
	fruit := yellowFruit.Union(redFruit)
	fmt.Println(fruit)

	// create route
	http.HandleFunc("/hello", helloWorld)
	// run server
	http.ListenAndServe(":8282", nil)
}
