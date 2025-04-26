package handlers

import (
	"database/sql"
	"db_driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"types"
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

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Can't read payload: %s\n", err)
	log.Printf("Request from %s not fulfilled, bad request: %s\n", r.Host, err)
}
func badResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("Bad Request: %s\n", err)
	fmt.Fprintf(w, "[%s] Request not fulfilled, contact api developers for more data\n", r.Host)
}

func badMethod(w http.ResponseWriter, r *http.Request, methods []string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Bad Method, allowed methods: %s\n", methods)
	log.Printf("[%s] Request not fulfilled, bad method: %s, allowed: %s\n", r.Host, r.Method, methods)
}

func getProjectId(w http.ResponseWriter, r *http.Request) *string {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params.Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad project id provided \n")
		log.Printf("[%s] Request failed, bad project id\n", r.Host)
		return nil
	}
	return &id
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

func GetTagCreator(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			badMethod(w, r, []string{"post"})
			return
		}
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		log.Printf("[%s] [POST] Received a post tag request from %s\n", id, r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.TagJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		tags := make([]types.TagJson, 0)
		tags = append(tags, reqData)
		newTags, err := db_driver.CreateTags(db_driver.CreateAgentDB(db), id, &tags)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(newTags)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		fmt.Fprint(w, string(data))
		log.Printf("[%s] Updated succesfully", id)
	}
	return handler
}

func GetTagDeleter(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			badMethod(w, r, []string{"delete"})
			return
		}
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		log.Printf("[%s] [DELETE] Received a delete tag request from %s\n", id, r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData struct {
			Id string `json:"id"`
		}
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		_, err = db.Exec("DELETE FROM Tags WHERE id = ?", reqData.Id)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Deleted succesfully")
		log.Printf("[%s] Deleted succesfully", id)
	}
	return handler
}
