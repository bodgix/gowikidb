package main

import (
	"database/sql"
	"fmt"
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

func getTitle(uri string, r *http.Request) string {
	return r.URL.Path[len(uri):]
}

func getArticle(db *sql.DB, title string) (string, error) {
	var body string
	row := db.QueryRow("SELECT body FROM articles WHERE title = ?", title)
	err := row.Scan(&body)
	return body, err
}

// Returning a closure with the handler and access to the DB
func viewHandler(db *sql.DB, tpls templatesMap) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the article from the DB
		var body string
		var err error
		title := getTitle("/view/", r)

		if body, err = getArticle(db, title); err != nil {
			var httpError int
			switch err {
			case sql.ErrNoRows:
				http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			default:
				httpError = http.StatusInternalServerError
			}
			http.Error(w, err.Error(), httpError)
			return
		}
		tpls["/view/"].Execute(w, &wikiPage{title, body})
	})
}

func editHandler(db *sql.DB, tpls templatesMap) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := getTitle("/edit/", r)
		body, _ := getArticle(db, title)
		tpls["/edit/"].Execute(w, &wikiPage{Title: title, Body: body})
	})
}

func saveHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := getTitle("/edit/", r)
		body := r.FormValue("body")
		db.Query("INSERT INTO `articles` (title, body) VALUES(?, ?) ON DUPLICATE KEY UPDATE body = VALUES(body)", title, body)
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	})
}

func withLogger(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Println("Handling a request for", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/gowiki")
	if err != nil {
		log.Fatal(err.Error())
	}

	tpls := make(templatesMap)
	tpls["/view/"], err = template.ParseFiles("templates/view.html")
	tpls["/edit/"], err = template.ParseFiles("templates/edit.html")
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := log.New(os.Stdout, "", 0)

	http.Handle("/view/", withLogger(logger, viewHandler(db, tpls)))
	http.Handle("/edit/", withLogger(logger, editHandler(db, tpls)))
	http.Handle("/save/", withLogger(logger, saveHandler(db)))
	fmt.Println("Starting the server on port 8080")
	http.ListenAndServe(":8080", nil)
}
