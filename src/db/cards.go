package db_driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"types"
	"utils"

	"github.com/KustelR/jsondiff"
)

func UpdateCard(db *sql.DB, card *types.CardJson) (*types.CardJson, error) {
	tx, err := db.BeginTx(context.Background(), nil)
	agent := CreateAgentTX(tx)
	if err != nil {
		return nil, err
	}
	stmt, err := agent.Prepare("CALL update_card(?, ?, ?, ?, ?, ?);")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	oldCard, err := GetCard(agent, card.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	oldJson, err := json.Marshal(oldCard.Json())

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newJson, err := json.Marshal(*card)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	diff1, diff2 := jsondiff.Diff(oldJson, newJson)
	stmtUR, err := agent.Prepare("CALL create_card_update_record(?, ?, ?);")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = stmtUR.Exec(card.Id, string(diff1), string(diff2))
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newCard := *card
	if oldCard.ColumnId != newCard.ColumnId {
		stmtPop, err := agent.Prepare("CALL pop_card_reorder(?, ?);")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer stmt.Close()
		_, err = stmtPop.Exec(oldCard.ColumnId, oldCard.Order)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		dbColNames, data, err := readOneRow(agent, card.ColumnId, "SELECT max(draw_order) FROM Cards WHERE column_id = ?;")
		if err != nil {
			return nil, err
		}
		maxDrawOrder, err := GetMaxDrawOrder(dbColNames, data)
		if err != nil {
			return nil, err
		}
		newCard.Order = maxDrawOrder + 1
	}
	_, err = stmt.Exec(newCard.Id, newCard.ColumnId, newCard.Name, newCard.Description, "placeholder", newCard.Order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &newCard, nil
}

func CreateCardTags(agent *Agent, cardId string, tagId string) error {
	stmt, err := agent.Prepare(`
	INSERT CardsTags 
		(card_id, tag_id) 
	VALUES 
	(
		?, # card_id 
		? # tag_id
	);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(cardId, tagId)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func RemoveCardTags(agent *Agent, cardId string, tagId string) error {
	stmtCT, err := agent.Prepare(`
	DELETE FROM CardsTags WHERE card_id = ? AND tag_id = ?;`)
	if err != nil {
		return err
	}
	stmtCT.Exec(cardId, tagId)
	stmtCT.Close()
	return nil
}

func DeleteCard(db *sql.DB, id string) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	agent := CreateAgentTX(tx)
	stmt, err := agent.Prepare("DELETE FROM Cards WHERE id = ?;")
	if err != nil {
		tx.Rollback()
		return err
	}
	card, err := GetCard(agent, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	stmtPC, err := agent.Prepare("CALL pop_card_reorder(?, ?);")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtPC.Exec(card.ColumnId, card.Order)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func CreateCards(agent *Agent, columnId string, cards *[]types.CardJson) ([]types.CardJson, error) {
	stmt, err := agent.Prepare(`
	CALL create_card(?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	stmtRT, err := agent.Prepare(`select * from CardsTags 
	where 
	card_id = ? AND # card id
	tag_id = ? # tag id`)
	if err != nil {
		return nil, err
	}
	defer stmtRT.Close()
	var cardErr error
	newCards := make([]types.CardJson, len(*cards))
out:
	for idx, card := range *cards {
		id := utils.GetUUID()
		changedCard := card
		changedCard.Id = id

		cols, data, err := readOneRow(agent, card.ColumnId, "SELECT max(draw_order) FROM Cards WHERE column_id= ?;")
		if err != nil {
			cardErr = err
			break out
		}
		drawOrder, err := GetMaxDrawOrder(cols, data)
		if err != nil {
			cardErr = err
			break out
		}
		changedCard.Order = drawOrder + 1
		_, err = stmt.Exec(columnId, changedCard.Id, changedCard.Name, changedCard.Description, changedCard.Order, "placeholder")
		if err != nil {
			cardErr = err
			break out
		}
		for _, tagId := range card.TagIds {
			row := stmtRT.QueryRow(card.Id, tagId)
			err := row.Scan()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = CreateCardTags(agent, id, tagId)
				} else {
					cardErr = err
					break out
				}
				if err != nil {
					cardErr = err
					break out
				}

			}
		}
		newCards[idx] = changedCard
	}
	if cardErr != nil {
		return nil, cardErr
	}
	return newCards, nil
}
