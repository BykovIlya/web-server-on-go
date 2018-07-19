package main

import (
	"fmt"
	"github.com/BykovIlya/web-server-on-go/models"
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
)

var posts map[string]*models.Post
var counter int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "write", nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
		return
	}

	t.ExecuteTemplate(w, "write", post)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	http.Redirect(w, r, "/", 302)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	delete(posts, id)

	http.Redirect(w, r, "/", 302)
}

func main() {
	fmt.Println("Listening on port : 3000")

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write" {
			counter++
		}
	})
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assest", staticOptions))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit", editHandler)
	m.Get("/delete", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Get("/test", func() string {
		return "test"
	})

	m.Run()
}
