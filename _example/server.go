package main

import (
	"github.com/yohcop/openid-go"
	"html/template"
	"log"
	"net/http"
)

const dataDir = "_example/"

// For the demo, we use in-memory infinite storage nonce and discovery
// cache. In your app, do not use this as it will eat up memory and never
// free it. Use your own implementation, on a better database system.
// If you have multiple servers for example, you may need to share at least
// the nonceStore between them.
var nonceStore = &openid.SimpleNonceStore{
	Store: make(map[string][]*openid.Nonce)}
var discoveryCache = &openid.SimpleDiscoveryCache{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]string)
	if t, err := template.ParseFiles(dataDir + "index.html"); err == nil {
		t.Execute(w, p)
	} else {
		log.Print(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]string)
	if t, err := template.ParseFiles(dataDir + "login.html"); err == nil {
		t.Execute(w, p)
	} else {
		log.Print(err)
	}
}

func discoverHandler(w http.ResponseWriter, r *http.Request) {
	if url, err := openid.RedirectURL(r.FormValue("id"),
		"http://localhost:8080/openidcallback",
		""); err == nil {
		http.Redirect(w, r, url, 303)
	} else {
		log.Print(err)
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fullUrl := "http://localhost:8080" + r.URL.String()
	log.Print(fullUrl)
	id, err := openid.Verify(
		fullUrl,
		discoveryCache, nonceStore)
	if err == nil {
		p := make(map[string]string)
		p["user"] = id
		if t, err := template.ParseFiles(dataDir + "index.html"); err == nil {
			t.Execute(w, p)
		} else {
			log.Println("WTF")
			log.Print(err)
		}
	} else {
		log.Println("WTF2")
		log.Print(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/discover", discoverHandler)
	http.HandleFunc("/openidcallback", callbackHandler)
	http.ListenAndServe(":8080", nil)
}
