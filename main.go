package main

import (
	"log"
	"os"
	"net/http"
	"html/template"
) 

func GeneralMiddleware(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	rf := func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}

	return rf
}

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/post.html")

	if err != nil {
		log.Fatal("Failure!!", err)
	}

	files, err := os.ReadDir("templates/posts")

	if err != nil {
		log.Fatal("Failure!!!", err)
	}

	html := make([]template.HTML, 0)

	for _, dirEntry := range files {
		fileBytes, err := os.ReadFile("templates/posts/" + dirEntry.Name())

		if err != nil {
			log.Fatal("Failure!!!", err)
		}

		html = append(html, template.HTML(string(fileBytes)))
	}

	t.Execute(w, html)
}

func main() {
	fs_css := http.FileServer(http.Dir("./assets/css"))

	mux := http.NewServeMux()

	mux.Handle("/css/", http.StripPrefix("/css/", fs_css))

	mux.HandleFunc("/", GeneralMiddleware(Home))

	err := http.ListenAndServe(":8080", mux)
	
	if err != nil {
		log.Fatalf("Unable to start web server %v\n", err)
	}
}
