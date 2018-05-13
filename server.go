package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type server struct {
	//db     *someDatabase
	router *http.ServeMux
	//email  *emailSender
}

func (s *server) routes() {
	s.router.HandleFunc("/api/", s.handleAPI())
	s.router.HandleFunc("/about", s.handleAbout())
	s.router.HandleFunc("/", s.handleIndex("tsumutsumu"))
	s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
}

func (s *server) handleAPI() http.HandlerFunc {
	doPrepare()
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greetinig"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if name := params.Get("name"); name != "" {
			req := request{Name: name}
			greeting := nameToGreeting(req.Name)
			resp := response{Greeting: greeting + ", " + req.Name}
			json.NewEncoder(w).Encode(resp)
		} else {
			json.NewEncoder(w).Encode(response{Greeting: "Hello, This is API."})
		}
	}
}

func (s *server) handleAbout() http.HandlerFunc {
	var (
		init sync.Once
		// tpl  *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() { doHeavy() })
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "I am about!")
	}

}

func (s *server) handleIndex(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is index! Passed data is %s.", name)
	}
}

func (s *server) handleAdminIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, admin!")
	}
}

// middleware functioin
/// decide whether call original handller or not
func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if false {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

// mock function
/// do some task
func doPrepare() {
	time.Sleep(100 * time.Millisecond)
}

/// do heavy task
func doHeavy() {
	time.Sleep(2 * time.Second)
}

/// do real task
func nameToGreeting(name string) string {
	if name == "P" {
		return "bakananodesuka"
	}
	return "Hello"
}

func main() {
	// net Listener
	l, _ := net.Listen("tcp", fmt.Sprintf(":8080"))
	// net/http Router
	router := http.NewServeMux()
	// Server instance
	s := server{router: router}
	s.routes()
	http.Serve(l, s.router)
}
