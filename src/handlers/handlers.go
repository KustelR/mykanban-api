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
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Bad Request: %s\n", err)
	log.Printf("[%s] Request not fulfilled, bad request: %s\n", r.Host, err)
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
		err = db_driver.AddCardTags(db_driver.CreateAgentDB(db), reqData.CardId, reqData.TagId)
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
		newTags, err := db_driver.AddTags(db_driver.CreateAgentDB(db), id, &tags)
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
		newCards, err := db_driver.AddCards(db_driver.CreateAgentDB(db), reqData.ColumnId, &cards)
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
		db_driver.UpdateCard(db, &reqData)
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

func GetCardForcePopOrder(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			badMethod(w, r, []string{"patch"})
			return
		}
		log.Printf("[%s] [PATCH] Received a pop card order request\n", r.Host)
		decoder := json.NewDecoder(r.Body)
		var reqData struct {
			Order    int    `json:"order"`
			ColumnId string `json:"columnId"`
		}
		err := decoder.Decode(&reqData)
		if err != nil {
			if err != io.EOF {
				badRequest(w, r, err)
				return
			}
		}
		_, err = db.Exec("CALL pop_card_reorder(?, ?)", reqData.ColumnId, reqData.Order)
		if err != nil {
			badResponse(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Updated succesfully")
		log.Printf("Update succesfully")
	}
	return handler
}

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
		db_driver.UpdateColumnData(db, &colData)
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
		newColumns, err := db_driver.AddColumns(db_driver.CreateAgentDB(db), *id, &columns)
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
