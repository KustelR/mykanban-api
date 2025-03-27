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
		log.Printf("Error when loading .env file: %s\n", err)
	}
	connectionString := os.Getenv("MYSQL_CONNECTION_STRING")
	if connectionString == "" {
		panic(fmt.Errorf("provide connection string to a database via MYSQL_CONNECTION_STRING enviroment variable"))
	}
	port := os.Getenv("PORT")
	if port == "" {
		panic(fmt.Errorf("provide port via PORT enviroment variable"))
	}

	db := db_driver.GetDb(connectionString)
	getHandler := newHandler(handlers.GetProjectGetter(db))
	postHandler := newHandler(handlers.GetProjectCreator(db))
	deleteHandler := newHandler(handlers.GetProjectDeleter(db))
	updateHandler := newHandler(handlers.GetProjectUpdater(db))
	updateDataHandler := newHandler(handlers.GetProjectDataUpdater(db))

	http.Handle("/read", getHandler)
	http.Handle("/create", postHandler)
	http.Handle("/delete", deleteHandler)
	http.Handle("/update", updateHandler)
	http.Handle("/data/update", updateDataHandler)

	cardCreateHandler := newHandler(handlers.GetCardCreator(db))
	cardUpdateHandler := newHandler(handlers.GetCardUpdater(db))
	cardDeleteHandler := newHandler(handlers.GetCardDeleter(db))

	http.Handle("/cards/create", cardCreateHandler)
	http.Handle("/cards/update", cardUpdateHandler)
	http.Handle("/cards/delete", cardDeleteHandler)

	addTagToCardHandler := newHandler(handlers.GetCardTagAdder(db))
	removeTagFromCardHandler := newHandler(handlers.GetCardTagRemover(db))
	postTagHandler := newHandler(handlers.GetTagCreator(db))
	deleteTagHandler := newHandler(handlers.GetTagDeleter(db))

	http.Handle("/tags/create", postTagHandler)
	http.Handle("/tags/delete", deleteTagHandler)
	http.Handle("/tags/link", addTagToCardHandler)
	http.Handle("/tags/unlink", removeTagFromCardHandler)

	forceReorderHandler := newHandler(handlers.GetCardForcePopOrder(db))
	columnDataUpdateHandler := newHandler(handlers.GetColumnDataUpdater(db))
	columnDeleteHandler := newHandler(handlers.GetColumnDeleter(db))
	columnCreateHandler := newHandler(handlers.GetColumnCreator(db))

	http.Handle("/columns/create", columnCreateHandler)
	http.Handle("/columns/update", columnDataUpdateHandler)
	http.Handle("/columns/delete", columnDeleteHandler)
	http.Handle("/columns/force_reorder", forceReorderHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go serve(port, &wg)
	wg.Wait()
	defer db.Close()
}

func cors(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		handler(w, r)
	}
}
