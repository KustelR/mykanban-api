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

func GetCardCreator(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			badMethod(w, r, []string{"post"})
			return
		}
		log.Printf("[POST] Received a update card request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.CardJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		cards := []types.CardJson{reqData}
		newCards, err := db_driver.CreateCards(db_driver.CreateAgentDB(db), reqData.ColumnId, &cards)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(newCards)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		fmt.Fprint(w, string(data))
		log.Printf("[POST] Created succesfully\n")
	}
	return handler
}

func GetCardUpdater(db *sql.DB) http.HandlerFunc {

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			badMethod(w, r, []string{"put"})
			return
		}
		log.Printf("[PUT] Received a update card request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData types.CardJson
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		res, err := db_driver.UpdateCard(db, &reqData)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		marshRes, err := json.Marshal(res)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		fmt.Fprint(w, string(marshRes))
		log.Printf("[PUT] Updated succesfully\n")
	}
	return handler
}

func GetCardDeleter(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			badMethod(w, r, []string{"delete"})
			return
		}
		log.Printf("[%s] [DELETE] Received a delete card request\n", r.Host)
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
		err = db_driver.DeleteCard(db, reqData.Id)
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

func GetCardTagAdder(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			badMethod(w, r, []string{"put"})
			return
		}
		log.Printf("[PUT] Received a link tag to card request from %s\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData struct {
			CardId string `json:"cardId"`
			TagId  string `json:"tagId"`
		}
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		err = db_driver.CreateCardTags(db_driver.CreateAgentDB(db), reqData.CardId, reqData.TagId)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Tag linked succesfully")
		log.Printf("Tag linked succesfully")
	}
	return handler
}

func GetCardTagRemover(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			badMethod(w, r, []string{"delete"})
			return
		}
		params, _ := url.ParseQuery(r.URL.RawQuery)
		id := params.Get("id")
		log.Printf("[%s] [DELETE] Received a unlink tag from card request from %s\n", id, r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData struct {
			CardId string `json:"cardId"`
			TagId  string `json:"tagId"`
		}
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		err = db_driver.RemoveCardTags(db_driver.CreateAgentDB(db), reqData.CardId, reqData.TagId)
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
