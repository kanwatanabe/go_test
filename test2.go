package main

// import (
// 	"fmt"
// 	"net/http"
// )

// type HelloHandler struct{}

// func (h *HelloHandler) ServerHTTP (w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello!")
// }
// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello!")
// }

// func main() {
// 	// hello := HelloHandler{}

// 	server := http.Server{
// 		Addr: "127.0.0.1:8080",
// 	}

// 	http.HandleFunc("/hello", hello)

// 	server.ListenAndServe()
// }