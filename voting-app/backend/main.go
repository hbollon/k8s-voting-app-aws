package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

const (
	host      = "db"
	port      = 5432
	user      = "postgres"
	password  = "password"
	dbname    = "voting-app-db"
	tablename = "votes"
)

type worker struct {
	db    *sql.DB
	redis *redis.Client
}

type vote struct {
	VoterID string `json:"voter_id"`
	Value   int    `json:"vote"`
}

func (w *worker) initConnections() (err error) {
	// Init Postgresql connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	w.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}

	// Create default votes table if not exists
	_, err = w.db.Exec("CREATE TABLE IF NOT EXISTS " +
		tablename + " (id TEXT PRIMARY KEY, voter_id VARCHAR(45) NOT NULL, vote INTEGER NOT NULL)")
	if err != nil {
		return
	}

	// Init Redis connection
	w.redis = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       0, // use default DB
	})

	return
}

func (w *worker) closeConnections() {
	w.db.Close()
	w.redis.Close()
}

func main() {
	// Init worker
	w := worker{}
	err := w.initConnections()
	if err != nil {
		panic(err)
	}
	defer w.closeConnections()

	// Start worker's listening routine
	fmt.Println("Waiting for votes...")
	for {
		// Waiting for new votes in redis
		results, err := w.redis.BLPop(context.Background(), 0, "votes").Result()
		if err != nil {
			panic(err)
		}

		// Parse first vote to struct
		var voteObj vote
		err = json.Unmarshal([]byte(results[0]), &voteObj)
		if err != nil {
			panic(err)
		}
		fmt.Println(voteObj)

		// Process vote in db
		err = w.processVote(voteObj)
		if err != nil {
			panic(err)
		}
	}
}

func (w *worker) processVote(v vote) (err error) {
	stmt, err := w.db.Prepare("UPDATE votes SET vote = ? WHERE id = ?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(v.Value, v.VoterID)
	return
}
