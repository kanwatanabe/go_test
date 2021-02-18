package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"time"
	"database/sql"
	_"github.com/lib/pq"
)

type Post struct {
	User string
	Threads []string
}

func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}
func headers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	h := r.Header["Accept-Encoding"]
	fmt.Fprintln(w, h)
}

func writeHeaderExample(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "そのようなサービスはありません。")
}
func headerExample(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(302)
}

func jsonExample(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User: "Test_Name",
		Threads: []string{"1番目", "2番目", "3番目"},
	    }
		json,_:= json.Marshal(post)
		w.Write(json)
}

func setCookie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c1 := http.Cookie{
		Name: "first_cookie",
		Value: "Go Web Programming",
		HttpOnly: true,
	}
	c2 := http.Cookie{
		Name: "second_cookie",
		Value: "Manning Publications Co",
		HttpOnly: true,
	}
	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
}
func getCookie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// h := r.Header["Cookie"]
	// fmt.Fprintln(w, h)
	c1, err := r.Cookie("first_cookie")
	if err != nil {
		fmt.Fprintln(w, "Cannot get the first cookie")
	}
	cs := r.Cookies()
	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)
}

func setMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	msg := []byte("Hello World!")
	c := http.Cookie{
		Name: "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}
	http.SetCookie(w, &c)
}
func showMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "メッセージはありません。")
		}
	}else {
		rc := http.Cookie{
			Name: "flash",
			MaxAge: -1,
			Expires: time.Unix(1, 0),
		}
		http.SetCookie(w, &rc)
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	}
}

type Posts struct {
	Id int
	Content string
	Author  string
}

var PostsById map[int] *Posts
var PostsByAuthor map[string] []*Posts

func store(posts Posts) {
	PostsById[posts.Id] = &posts
	PostsByAuthor[posts.Author] = append(PostsByAuthor[posts.Author], &posts)
} 

// DB作成
var Db *sql.DB
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}
func Postss(limit int) (posts []Posts, err error) {
	rows, err := Db.Query("select id, content, author from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Posts{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts,post)
	}
	rows.Close()
	return
}

func (post *Posts) Create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	if err != nil {
		return
	}
	return
}

func GetPost(id int) (post Posts, err error) {
	post = Posts{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

func main() {

	PostsById = make(map[int] *Posts)
	PostsByAuthor = make(map[string] []*Posts)

	posts1 := Posts{Id: 1, Content: "Hello Wold", Author: "Sau Sheong"}
	posts2 := Posts{Id: 2, Content: "Bonjour Wold", Author: "Pierre"}
	posts3 := Posts{Id: 3, Content: "Hola Wold", Author: "Pedro"}
	posts4 := Posts{Id: 4, Content: "Greentings Wold", Author: "Sheong"}

	store(posts1)
	store(posts2)
	store(posts3)
	store(posts4)

	fmt.Println(PostsById[1])
	fmt.Println(PostsById[2])

	post := Posts{Content: "Hello", Author: "Sau"}
	fmt.Println(post)
	post.Create()
	fmt.Println(post)

	readPost, _ := GetPost(post.Id)
    fmt.Println(readPost)

	for _, post := range PostsByAuthor["Sau Sheong"] {
		fmt.Println(post)
	}
	for _, post := range PostsByAuthor["Pedro"] {
		fmt.Println(post)
	}

	mux := httprouter.New()
	mux.GET("/hello/:name",hello)
	mux.GET("/headers",headers)
	// 
	mux.GET("/writeheader", writeHeaderExample)
	mux.GET("/headerexample", headerExample)
	mux.GET("/json", jsonExample)

	mux.GET("/setcookie", setCookie)
	mux.GET("/getcookie", getCookie)

	mux.GET("/setmessage",setMessage)
	mux.GET("/showmessage",showMessage)

	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
