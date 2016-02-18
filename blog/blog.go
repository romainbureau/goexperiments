package main

import (
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var port int
var cwd string
var templates string

type Post struct {
	Title   string
	Content string
}

func init() {
	flag.IntVar(&port, "port", 8000, "listening port")
	flag.StringVar(&templates, "templates", ".", "templates directory")
	flag.Parse()

	cwd, _ = os.Getwd()

	log.Println(fmt.Sprintf("cwd: %s", cwd))
	log.Println(fmt.Sprintf("templates directory: %s", templates))
}

func before(r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	before(r)
	postTemplate := filepath.Join(templates, "./index.html")

	t, _ := template.ParseFiles(postTemplate)

	post := Post{Title: "my title", Content: string([]byte("my post content"))}
	t.Execute(w, post)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	before(r)
	funcMap := template.FuncMap{
		"Markdown": Markdown,
	}

	postTemplate := filepath.Join(templates, "./post.html")

	t, _ := template.New("post.html").Funcs(funcMap).ParseFiles(postTemplate)

	post := Post{Title: "my title", Content: "# this is my content\n - test\n - test"}
	t.Execute(w, post)
}

func Markdown(input string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(input))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(html)
}

func main() {
	log.Println(fmt.Sprintf("listening on %d", port))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post/", postHandler)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
