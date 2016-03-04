package main

import (
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var port int
var cwd string
var templates string
var data string

type Posts struct {
	Posts []Post
}

type Post struct {
	Filename string
	Content  string
}

func init() {
	flag.IntVar(&port, "port", 8000, "listening port")
	flag.StringVar(&templates, "templates", ".", "templates directory")
	flag.StringVar(&data, "data", "data", "data directory, where .md files are")
	flag.Parse()

	cwd, _ = os.Getwd()

	log.Println(fmt.Sprintf("cwd: %s", cwd))
	log.Println(fmt.Sprintf("templates directory: %s", templates))
	log.Println(fmt.Sprintf("data directory: %s", data))
}

func before(r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getFilename(urlPath string) string {
	return path.Base(urlPath)
}

func formatPostTitle(title string) string {
	a := strings.Replace(title, "-", " ", -1)
	b := strings.Title(a)

	return b
}

func readPost(file string) (Post, error) {
	md, err := ioutil.ReadFile(filepath.Join(data, "/"+file))
	if err != nil {
		return Post{Filename: "", Content: ""}, err
	}
	return Post{Filename: file, Content: string(md)}, nil
}

func Markdown(input string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(input))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(html)
}

var funcMap = template.FuncMap{
	"Markdown": Markdown,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	before(r)
	indexTemplate := filepath.Join(templates, "/index.html")
	t, _ := template.New("index.html").ParseFiles(indexTemplate)

	files, _ := filepath.Glob(filepath.Join(data, "/*.md"))

	posts := []Post{}
	for _, file := range files {
		post, err := readPost(getFilename(file))
		if err == nil {
			posts = append(posts, post)
		}
	}

	t.Execute(w, Posts{Posts: posts})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	before(r)
	postTemplate := filepath.Join(templates, "/post.html")
	t, _ := template.New("post.html").Funcs(funcMap).ParseFiles(postTemplate)

	post, err := readPost(getFilename(r.URL.Path))
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	} else {
		t.Execute(w, post)
	}
}

func noneHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	log.Println(fmt.Sprintf("listening on %d", port))

	http.HandleFunc("/favicon.ico", noneHandler)
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
