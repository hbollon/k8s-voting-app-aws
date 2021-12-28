package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "voting-app-db"
)

var randGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

type Redis struct {
	*redis.Client
}

type vote struct {
	VoterID string `json:"voter_id"`
	Value   int    `json:"vote"`
}

// Pass the Redis client to API handlers
// This takes a callback and returns a HandlerFunc which calls the callback with the Redis client
func handleWithRedis(apiF func(w http.ResponseWriter, r *http.Request,
	d *Redis), d *Redis) func(http.ResponseWriter, *http.Request) {
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

func submitVoteHandler(w http.ResponseWriter, r *http.Request, red *Redis) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read body: %v", err)
		JSONError(w, "Failed to read body during vote submit", err)
		return
	}

	if err = red.submitVote(body); err != nil {
		log.Errorf("Failed to submit vote to Redis: %v", err)
		JSONError(w, "Failed to submit vote to Redis", err)
		return
	}
}

func (red *Redis) submitVote(voteData []byte) (err error) {
	var voteObj vote
	if err = json.Unmarshal(voteData, &voteObj); err != nil {
		return
	}
	voteObj.VoterID = strconv.Itoa(randGenerator.Intn(10000))

	// Marshal voteObj to JSON
	var voteJson []byte
	if voteJson, err = json.Marshal(voteObj); err != nil {
		return
	}

	// Push the vote to Redis
	if err := red.LPush(context.Background(), "votes", voteJson).Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Set up the Redis connection
	red := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       0, // use default DB
	})

	// Instantiate gorilla/mux router instance
	r := mux.NewRouter()

	// Handle API endpoints
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/vote", handleWithRedis(submitVoteHandler, &Redis{red}))

	// Serve static files (CSS, JS, images) from dir
	fs := http.FileServer(http.Dir("/static"))
	r.PathPrefix("/").Handler(fs)

	// Add CORS Middleware to mux router
	r.Use(corsMiddleware)

	// Start server
	port := 80
	log.Infof("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
