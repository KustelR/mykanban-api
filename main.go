package main

import (
	"db_driver"
	"fmt"
	"handlers"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type HttpServer struct {
	handler http.HandlerFunc
}

func (srv *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := cors(srv.handler)
	handler(w, r)
}

func newHandler(handlerFunc http.HandlerFunc) *HttpServer {
	srv := &HttpServer{handlerFunc}

	return srv
}

func serve(port string, wg *sync.WaitGroup) {
	fmt.Printf("Server is running on %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
	wg.Done()
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("MYSQL_CONNECTION_STRING")

	db := db_driver.GetDb(connectionString)
	port := "7501"
	getHandler := newHandler(handlers.GetProjectGetter(db))
	postHandler := newHandler(handlers.GetProjectCreator(db))
	deleteHandler := newHandler(handlers.GetProjectDeleter(db))
	updateHandler := newHandler(handlers.GetProjectUpdater(db))
	http.Handle("/read", getHandler)
	http.Handle("/create", postHandler)
	http.Handle("/delete", deleteHandler)
	http.Handle("/update", updateHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go serve(port, &wg)
	wg.Wait()
	defer db.Close()
}

func cors(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		handler(w, r)
	}
}
