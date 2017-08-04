package main

import (
	"net/http"
	"html/template"
	"os"
	"bufio"
	"log"

	"github.com/gorilla/mux"
	"fmt"
	"strconv"
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

// Write the given string array to the specified file, overwriting any lines
// that were previously in the file.
func writeFile( lines []*Todo ) {
	file, err := os.Create( fileName ) //example of multiple results from a function where one is the error code
	if err != nil {
		panic( "could not open todo file" )
	}
	defer file.Close() //will call file's close function at the end of writeFile

	w := bufio.NewWriter( file )
	defer w.Flush() //interesting, two deferred funcs, one needs to be called first....

	for _, each := range lines { //ignore the first param with "_"
		fmt.Fprint( w, each.Line + "\n" )
	}
}

func deleteFromFile(i int, t []*Todo) {
	i = i - 1
	if i < len(t) {
		t = append(t[:i], t[i+1:]...) // https://github.com/golang/go/wiki/SliceTricks ??
	}

	writeFile(t)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, _ := strconv.Atoi(vars["id"])
	c := ReadFile()
	deleteFromFile(i, c)
	ViewHandler(w, r)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("handler called")
	value := r.FormValue("entry")
	current := ReadFile()
	current = append(current, &Todo{Line: value})
	writeFile(current)
	ViewHandler(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", AddHandler).Methods("POST")
	r.HandleFunc("/{id}", DeleteHandler)
	r.HandleFunc("/", ViewHandler)
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/assets/").Handler(s)
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
