package db_driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"types"
)

func UpdateCard(db *sql.DB, card *types.CardJson) error {
	tx, err := db.BeginTx(context.Background(), nil)
	agent := CreateAgentTX(tx)
	if err != nil {
		return err
	}
	stmt, err := agent.Prepare("CALL update_card(?, ?, ?, ?, ?);")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	oldCard, err := GetCard(agent, card.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if oldCard.ColumnId != card.ColumnId {
		stmtPop, err := agent.Prepare("CALL pop_card_reorder(?, ?);")
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		fmt.Println(oldCard, card.ColumnId)
		_, err = stmtPop.Exec(oldCard.ColumnId, oldCard.Order)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = stmt.Exec(card.Id, card.ColumnId, card.Name, card.Description, card.Order)
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

func AddCardTags(agent *Agent, cardId string, tagId string) error {
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
	stmt.Exec(cardId, tagId)
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
	fmt.Println(card)

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func addCards(agent *Agent, columnId string, cards *[]types.CardJson) error {
	stmt, err := agent.Prepare(`
	CALL add_card(?, ?, ?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmtRT, err := agent.Prepare(`select * from CardsTags 
	where 
	card_id = ? AND # card id
	tag_id = ? # tag id`)
	if err != nil {
		return err
	}
	defer stmtRT.Close()
	var cardErr error
out:
	for _, card := range *cards {
		_, err := stmt.Exec(columnId, card.Id, card.Name, card.Description, card.Order)
		if err != nil {
			cardErr = err
			break out
		}
		for _, tagId := range card.TagIds {
			row := stmtRT.QueryRow(card.Id, tagId)
			err := row.Scan()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					AddCardTags(agent, card.Id, tagId)
				} else {
					cardErr = err
					break out
				}

			}
		}
	}
	if cardErr != nil {
		return cardErr
	}
	return nil
}
