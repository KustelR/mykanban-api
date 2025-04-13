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
	"utils"
)

func readProject(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params.Get("id")
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
			badResponse(w, r, err)
			return
		}
	}
	w.Write(data)
	log.Printf("[%s] Readed project to %s\n", id, r.Host)
}

func createProject(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Printf("[NEW] Received a post request from %s\n", r.Host)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var reqData types.KanbanJson
	err := decoder.Decode(&reqData)
	if err != nil {
		if err != io.EOF {
			badRequest(w, r, err)
			return
		}
	}
	id := utils.GetUUID()
	err = db_driver.CreateProject(db, id, &reqData)
	if err != nil {
		badResponse(w, r, err)
		return
	}
	fmt.Fprint(w, id)
	log.Printf("[%s] Created project\n", id)
}

func deleteProject(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params.Get("id")
	log.Printf("[%s] Received a delete request from %s\n", id, r.Host)
	res, err := db.Exec("DELETE FROM Projects WHERE id = ?", id)
	if err != nil {
		badResponse(w, r, err)
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

func updateProject(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params.Get("id")
	log.Printf("[%s] Received a put request from %s\n", id, r.Host)
	decoder := json.NewDecoder(r.Body)
	var reqData types.KanbanJson
	err := decoder.Decode(&reqData)
	if err != nil {
		if err != io.EOF {
			badRequest(w, r, err)
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
		badResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Updated succesfully")
	log.Printf("[%s] Updated succesfully", id)
}

func GetProjectDataUpdater(db *sql.DB) http.HandlerFunc {

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			badMethod(w, r, []string{"patch"})
			return
		}
		id := getProjectId(w, r)
		if id == nil {
			return
		}
		log.Printf("[PUT] Received a update project data request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData struct {
			Name string `json:"name"`
		}
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		_, err = db.Exec("CALL update_project_data(?, ?)", id, reqData.Name, "placeholder")
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Updated succesfully")
		log.Printf("[PUT] Updated succesfully\n")
	}
	return handler
}

func GetProjectRequestHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createProject(db, w, r)
			return
		case http.MethodPut:
			updateProject(db, w, r)
			return
		case http.MethodGet:
			readProject(db, w, r)
			return
		case http.MethodDelete:
			deleteProject(db, w, r)
			return
		default:
			badMethod(w, r, []string{"get", "put", "delete", "post"})
		}
	}
}
