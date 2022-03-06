package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type article struct {
	Id          int
	Name        string
	Description string
	Url         string
}

func dbConn() (db *sql.DB) {
	dbDriver := ""
	dbUser := ""
	dbPass := ""
	dbName := ""
	hostDb := ""
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+hostDb+")/"+dbName)
	// db, err := sql.Open("mysql", "<username>:<pw>@tcp(<HOST>:<port>)/<dbname>")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "GET" {
		selectDb, err := db.Query("SELECT * FROM articles ORDER BY id DESC")
		if err != nil {
			panic(err)
		}
		articles := article{}
		res := []article{}
		for selectDb.Next() {
			var id int
			var name, description, url string
			err = selectDb.Scan(&id, &name, &description, &url)
			if err != nil {
				panic(err)
			}
			articles.Id = id
			articles.Name = name
			articles.Description = description
			articles.Url = url
			res = append(res, articles)
		}
		fileJson, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		w.Write(fileJson)

	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Allowed")
	}
	defer db.Close()

}

func getArticleByID(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "GET" {
		nId := r.URL.Query().Get("id")
		selectDb, err := db.Query("SELECT * FROM articles WHERE id=?", nId)
		if err != nil {
			panic(err)
		}
		emp := article{}
		for selectDb.Next() {
			var id int
			var name, description, url string
			err = selectDb.Scan(&id, &name, &description, &url)
			if err != nil {
				panic(err)
			}
			emp.Id = id
			emp.Name = name
			emp.Description = description
			emp.Url = url
		}
		fileJson, err := json.Marshal(emp)
		if err != nil {
			panic(err)
		}
		w.Write(fileJson)
	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Allowed")
	}
	defer db.Close()
}

func newArticle(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		// name := r.FormValue("name")
		// description := r.FormValue("description")
		// url := r.FormValue("url")
		decoder := json.NewDecoder(r.Body)
		var articles article
		err := decoder.Decode(&articles)
		if err != nil {
			panic(err)
		}
		id := articles.Id
		name := articles.Name
		description := articles.Description
		url := articles.Url
		insForm, err := db.Prepare("insert into articles (id,name, description, url) values (?,?,?,?)")
		if err != nil {
			panic(err)
		}
		insForm.Exec(id, name, description, url)
		log.Println("INSERT: Name: " + name)
		w.WriteHeader(201)

	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Allowed")
	}
	defer db.Close()
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "DELETE" {
		decoder := json.NewDecoder(r.Body)
		var articles article
		err := decoder.Decode(&articles)
		if err != nil {
			panic(err)
		}
		id := articles.Id
		delForm, err := db.Prepare("DELETE FROM articles where id=?")
		if err != nil {
			panic(err)
		}
		delForm.Exec(id)
		log.Printf("DELETED: ArticleID: ", id)
		w.WriteHeader(201)
	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Allowed")
	}
	defer db.Close()
}

func main() {
	http.HandleFunc("/articles", getArticles)
	http.HandleFunc("/article", getArticleByID)
	http.HandleFunc("/create", newArticle)
	http.HandleFunc("/delete", deleteArticle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
