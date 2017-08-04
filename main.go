package main

import (
	"net/http"
	"html/template"
	"os"
	"bufio"
	"log"

	"github.com/gorilla/mux"
	"fmt"
)

type Todo struct {
	Line 	string
	Index 	int
}

const fileName = "./todo.dat"
func ReadFile() []*Todo {
	file, err := os.Open(fileName)
	if err != nil {
		return []*Todo{}
	}
	defer file.Close()

	s := bufio.NewScanner(bufio.NewReader(file))

	var tmp []*Todo
	c := 1
	for s.Scan() {
		tmp = append(tmp, &Todo{Line: s.Text(), Index: c})
		c = c + 1
	}

	return tmp

}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")

	td := ReadFile()
	t.Execute(w, td)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("entry")
	fmt.Printf("%s", value)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todo", ViewHandler)
	r.HandleFunc("/todo", AddHandler).Methods("POST")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
