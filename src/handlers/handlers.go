package handlers

import (
	"context"
	"database/sql"
	"db_driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"types"

	"github.com/google/uuid"
)

type postRequest struct {
	ActionType     string           `json:"type"`
	Position       int              `json:"position"`
	TagPayload     types.TagJson    `json:"tag"`
	CardPayload    types.CardJson   `json:"card"`
	ColumnPayload  types.ColumnJson `json:"column"`
	ProjectPayload types.KanbanJson `json:"project"`
}

func HandleDeleteRequest(db *sql.DB, id string, reader io.Reader) error {
	return nil
}

func HandleRequest(db *sql.DB, id string, reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	var reqData postRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}

func readProjectById(db *sql.DB, id string) ([]byte, error) {
	output, err := db_driver.GetProject(db, id)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetProjectGetter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		log.Printf("[%s] Received a get request from %s\n", id, r.Host)
		data, err := readProjectById(db, id)
		if err != nil {
			var nfe db_driver.NotFoundError
			if errors.As(err, &nfe) {
				w.WriteHeader(http.StatusNotFound)
				log.Printf("[%s] Get request not fulfilled, project not found\n", id)
				fmt.Fprint(w, err.Error())
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Bad Request: %s\n", err)
				log.Printf("[%s] Get request not fulfilled, bad request: %s\n", id, err)
				return
			}
		}
		w.Write(data)
		log.Printf("[%s] Readed project to %s\n", id, r.Host)
	}
}

func GetProjectCreator(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		log.Printf("[NEW] Received a post request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.KanbanJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Can't read payload: %s\n", err)
				log.Printf("Post request from %s not fulfilled, bad request: %s\n", r.Host, err)
				return
			}
		}
		id := uuid.New()
		err = db_driver.PostProject(db, id.String()[:30], &reqData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong on the server side: %s\n", err)
			log.Printf("[%s] Post request not fulfilled, internal error: %s\n", id, err)
			return
		}
		fmt.Fprint(w, id.String()[:30])
		log.Printf("[%s] Created project\n", id)
	}
}
func GetProjectDeleter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.NotFound(w, r)
			return
		}
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		log.Printf("[%s] Received a delete request from %s\n", id, r.Host)
		res, err := db.Exec("DELETE FROM Projects WHERE id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong on the serverside: %s\n", err)
			log.Printf("[%s] Delete request not fulfilled, internal error\n", id)
			return
		}
		affected, _ := res.RowsAffected()
		if affected <= 0 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Project %s not found\n", id)
			log.Printf("[%s] Delete request not fulfilled, project not found\n", id)
			return
		}
		log.Printf("[%s] Deleted project\n", id)
	}
}

func GetProjectUpdater(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.NotFound(w, r)
			return
		}
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		log.Printf("[%s] Received a put request from %s\n", id, r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.KanbanJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Can't read payload: %s\n", err)
				log.Printf("Put request from %s not fulfilled, bad request: %s\n", r.Host, err)
				return
			}
		}
		err = db_driver.UpdateProject(db, context.Background(), id, &reqData)
		if err != nil {
			if (err == db_driver.NoEffect{}) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Project %s not found\n", id)
				log.Printf("[%s] Put request not fulfilled, can't put with new id\n", id)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong on the server side: %s\n", err)
			log.Printf("[%s] Put request not fulfilled, internal error\n", id)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Updated succesfully")
		log.Printf("[%s] Updated succesfully", id)
	}
	return handler
}
