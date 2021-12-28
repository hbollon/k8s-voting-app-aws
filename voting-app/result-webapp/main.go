package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "voting-app-db"
)

type DB struct {
	*sql.DB
}

type vote struct {
	VoterID string `json:"voter_id"`
	Value   int    `json:"vote"`
}

// Pass the DB to API handlers
// This takes a callback and returns a HandlerFunc which calls the callback with the DB
func handleWithDB(apiF func(w http.ResponseWriter, r *http.Request,
	d *DB), d *DB) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, d)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		next.ServeHTTP(w, r)
	})
}

// JSONError is a wrapper function for errors
// which prints them to the http.ResponseWriter as a JSON response
func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

func getVotesHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	votesStats := db.getVotesStats()

	j, err := json.Marshal(votesStats)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

func (db *DB) getVotesStats() []int {
	rows, err := db.Query("SELECT vote, COUNT(id) AS count FROM votes GROUP BY vote")
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	defer rows.Close()

	votes := []int{0, 0}
	for rows.Next() {
		var v vote
		var nbVotes int
		if err := rows.Scan(&v.Value, &nbVotes); err != nil {
			log.Error(err.Error())
			return nil
		}
		votes[v.Value] = nbVotes
	}
	return votes
}

func main() {
	// Set up the DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// Instantiate gorilla/mux router instance
	r := mux.NewRouter()

	// Handle API endpoints
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/votes", handleWithDB(getVotesHandler, &DB{database}))

	// Serve static files (CSS, JS, images) from dir
	fs := http.FileServer(http.Dir("/static"))
	r.PathPrefix("/").Handler(fs)

	// Add CORS Middleware to mux router
	r.Use(corsMiddleware)

	// Start server
	port := 3000
	log.Infof("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
