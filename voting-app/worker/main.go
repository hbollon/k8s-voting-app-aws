package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	hostPostgres = "postgres"
	hostRedis    = "redis"
)

const (
	port      = 5432
	user      = "postgres"
	password  = "password"
	dbname    = "voting-app-db"
	tablename = "votes"
)

func init() {
	if os.Getenv("POSTGRES_HOST") != "" {
		hostPostgres = os.Getenv("POSTGRES_HOST")
	}
	if os.Getenv("REDIS_HOST") != "" {
		hostRedis = os.Getenv("REDIS_HOST")
	}
}

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
		hostPostgres, port, user, password, dbname)
	w.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}

	// Create default votes table if not exists
	_, err = w.db.Exec("CREATE TABLE IF NOT EXISTS " +
		tablename + " (id SERIAL PRIMARY KEY, voter_id VARCHAR(45) NOT NULL UNIQUE, vote INTEGER NOT NULL)")
	if err != nil {
		return
	}

	// Init Redis connection
	w.redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", hostRedis),
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
		err = json.Unmarshal([]byte(results[1]), &voteObj)
		if err != nil {
			panic(err)
		}
		fmt.Println(voteObj)

		// Process vote in db
		err = w.processInsertVote(voteObj)
		if err != nil {
			panic(err)
		}
	}
}

func (w *worker) processInsertVote(v vote) (err error) {
	stmt, err := w.db.Prepare("INSERT INTO votes (voter_id, vote) VALUES ($1, $2)")
	if err != nil {
		return
	}

	if _, err = stmt.Exec(v.VoterID, v.Value); err != nil {
		return w.processUpdateVote(v)
	}

	return
}

func (w *worker) processUpdateVote(v vote) (err error) {
	stmt, err := w.db.Prepare("UPDATE votes SET vote = $1 WHERE voter_id = $2")
	if err != nil {
		return
	}

	_, err = stmt.Exec(v.Value, v.VoterID)
	return
}
