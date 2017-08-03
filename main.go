package main

import (
	"net/http"
	"html/template"
	"os"
	"bufio"
)

type Todo struct {
	Line string
}

const fileName = "/tmp/todo.dat"
func ReadFile() []*Todo {
	file, err := os.Open(fileName)
	if err != nil {
		return []*Todo{}
	}
	defer file.Close()

	s := bufio.NewScanner(bufio.NewReader(file))

	var tmp []*Todo
	for s.Scan() {
		tmp = append(tmp, &Todo{Line: s.Text()})
	}

	return tmp

}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")

	td := ReadFile()
	t.Execute(w, td)
}

func main() {
	http.HandleFunc("/view", ViewHandler)
	http.ListenAndServe(":8080", nil)

}
