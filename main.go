package main

import (
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/yageek/socialios/twitter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	SESSION_NAME = "SCIOSSESSION"
	LoggedKey    = "logged"
)

var validUsers = map[string]string{
	"foo":   "bar",
	"user2": "password2",
}

type UserLog struct {
	Login    string `json="login"`
	Password string `json="password"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_STORE_SECRET")))

func main() {
	//Router
	mux := pat.New()
	mux.Post("/login", loginHandler)
	mux.Get("/tweets", authHandlerFunc(searchTweets))
	mux.Delete("/logout", logoutHandler)

	//Negroni classic instance
	n := negroni.Classic()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseFunc(jsonValider)
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))

}

func jsonValider(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
	} else {
		next(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)

	userlog := UserLog{}
	err := json.Unmarshal(data, &userlog)

	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		log.Println("Invalid payload:", err)
		return
	}

	if passwordValue, contains := validUsers[userlog.Login]; !contains || passwordValue != userlog.Password {
		http.Error(w, "unknown user", http.StatusUnauthorized)
	} else {

		session, err := store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("Could not create session:", err)
			http.Error(w, "Could not create session", http.StatusInternalServerError)
		} else {

			session.Values[LoggedKey] = true
			session.Save(r, w)

		}
	}
}

func searchTweets(w http.ResponseWriter, r *http.Request) {

	search_id := r.URL.Query().Get("max_id")
	tweets := twitter.Search("apple", search_id)
	log.Println("Tweets:", tweets)
	data, err := json.Marshal(tweets)
	if err != nil {
		http.Error(w, "Could not parse received data", http.StatusInternalServerError)
		log.Println("Err during marshaling:", err)
	} else {
		w.Header().Set("Content-Type", "application/json;application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func authHandler(handler http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, SESSION_NAME)
		log.Println("Session logged", session.Values[LoggedKey])
		if logged, ok := session.Values[LoggedKey].(bool); ok && logged {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func authHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return authHandler(http.HandlerFunc(handlerFunc))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, SESSION_NAME)

	session.Values[LoggedKey] = false
	session.Save(r, w)

}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fileData, _ := ioutil.ReadFile("api.html")
	w.WriteHeader(http.StatusOK)
	w.Write(fileData)

}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not found")
	w.WriteHeader(http.StatusNotFound)
}
