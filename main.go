package main

import (
	"log"
	"os"
	"net/http"
	"html/template"
	"strings"
) 

type BlogPost struct {
	Title string
	PostText template.HTML
}

func GeneralMiddleware(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	rf := func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}

	return rf
}

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/posts_page.html", "templates/post_card.html")

	if err != nil {
		log.Fatal("Failure!!", err)
	}

	files, err := os.ReadDir("templates/posts")

	if err != nil {
		log.Fatal("Failure!!!", err)
	}

	posts := make([]BlogPost, 0)

	for _, dirEntry := range files {
		post := BlogPost{}

		fileBytes, err := os.ReadFile("templates/posts/" + dirEntry.Name())

		if err != nil {
			log.Fatal("Failure!!!", err)
		}

		post.Title = strings.Split(dirEntry.Name(), ".")[0]
		post.PostText = template.HTML(string(fileBytes))

		posts = append(posts, post)
	}

	t.Execute(w, posts)
}

func Post(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/post_page.html", "templates/post.html")

	if err != nil {
		log.Fatal("Failed to parse templates")
	}

	name := r.URL.Query().Get("name")

	if len(name) < 1 {
		log.Fatal("Failed to get name.")
	}

	fileBytes, err := os.ReadFile("templates/posts/" + name + ".html")

	if err != nil {
		log.Fatal("Failed to read post")
	}

	post := BlogPost{}

	post.Title = name
	post.PostText = template.HTML(string(fileBytes))

	t.Execute(w, post)
}

func main() {
	fs_css := http.FileServer(http.Dir("./assets/css"))

	mux := http.NewServeMux()

	mux.Handle("/css/", http.StripPrefix("/css/", fs_css))

	mux.HandleFunc("/", GeneralMiddleware(Home))
	mux.HandleFunc("GET /post", GeneralMiddleware(Post))

	err := http.ListenAndServe(":8080", mux)
	
	if err != nil {
		log.Fatalf("Unable to start web server %v\n", err)
	}
}
