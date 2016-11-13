package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type templatesMap map[string]*template.Template

type wikiPage struct {
	Title string
	Body  string
}

// Returning a closure with the handler and access to the DB
func viewHandler(db *sql.DB, tpls templatesMap) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the article from the DB
		var body string
		title := r.URL.Path[len("/view/"):]
		row := db.QueryRow("SELECT body FROM articles WHERE title = %", title)
		if err := row.Scan(&body); err != nil {
			var httpError int
			switch err {
			case sql.ErrNoRows:
				httpError = http.StatusNotFound
			default:
				httpError = http.StatusInternalServerError
			}
			http.Error(w, err.Error(), httpError)
			return
		}
		tpls["/view/"].Execute(w, &wikiPage{title, body})
	})
}

func withLogger(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Println("Handling a request for", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sql.Open("mysql", "gowiki:gowiki@/gowiki")
	if err != nil {
		log.Fatal(err.Error())
	}

	tpls := make(templatesMap)
	tpls["/view/"], err = template.ParseFiles("templates/view.html")
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := log.New(os.Stdout, "", 0)

	http.Handle("/view/", withLogger(logger, viewHandler(db, tpls)))
	http.ListenAndServe(":8080", nil)
}
