package handlers

import (
	"database/sql"
	"db_driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"types"
)

func GetColumnDataUpdater(db *sql.DB) http.HandlerFunc {

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			badMethod(w, r, []string{"put"})
			return
		}
		id := getProjectId(w, r)
		if id == nil {
			return
		}
		log.Printf("[PUT] Received a update column data request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.ColumnJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		colData := types.Column{Id: reqData.Id, Name: reqData.Name, Order: reqData.Order, ProjectId: *id}
		err = db_driver.UpdateColumnData(db, &colData)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		newCol, err := db_driver.GetColumn(db_driver.CreateAgentDB(db), reqData.Id)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		marshRes, err := json.Marshal(newCol.Json())
		if err != nil {
			badResponse(w, r, err)
			return
		}
		fmt.Fprint(w, string(marshRes))
		w.WriteHeader(http.StatusOK)
		log.Printf("[PUT] Updated succesfully\n")
	}
	return handler
}

func GetColumnDeleter(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			badMethod(w, r, []string{"delete"})
			return
		}
		log.Printf("[%s] [DELETE] Received a delete column request\n", r.Host)
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
		err = db_driver.DeleteColumn(db, reqData.Id)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Deleted succesfully")
		log.Printf("Deleted succesfully")
	}
	return handler
}

func GetColumnCreator(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			badMethod(w, r, []string{"post"})
			return
		}
		id := getProjectId(w, r)
		if id == nil {
			return
		}
		log.Printf("[%s] [POST] Received a create column request from %s\n", *id, r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.ColumnJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		columns := make([]types.ColumnJson, 0)
		columns = append(columns, reqData)
		newColumns, err := db_driver.CreateColumns(db_driver.CreateAgentDB(db), *id, columns)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		data, err := json.Marshal(newColumns)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(data))
		log.Printf("[%s] Updated succesfully", *id)
	}
	return handler
}
