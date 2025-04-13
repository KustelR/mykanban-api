package db_driver

import (
	"context"
	"database/sql"
	"types"
)

func CreateProject(db *sql.DB, id string, projectData *types.KanbanJson) error {
	transaction, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	agent := CreateAgentTX(transaction)
	stmt, err := agent.Prepare(`CALL create_project(?, ?, ?);`)
	if err != nil {
		transaction.Rollback()
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id, projectData.Name, "placeholder")
	if err != nil {
		transaction.Rollback()
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		transaction.Rollback()
		return err
	}
	if rows == 0 {
		transaction.Rollback()
		return NoEffect{}
	}
	_, err = CreateTags(agent, id, &projectData.Tags)
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = CreateColumns(agent, id, projectData.Columns)
	if err != nil {
		transaction.Rollback()
		return err
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetProject(db *sql.DB, id string) (*types.KanbanJson, error) {
	var output types.KanbanJson
	project, err := ReadProject(db, id)
	if err != nil {
		return nil, err
	}
	output = *project.Json()
	projectTags, err := GetTagsByProject(db, id)
	if err != nil {
		return nil, err
	}
	for _, tag := range projectTags {
		outputTag := tag.Json()
		output.Tags = append(output.Tags, *outputTag)
	}
	columns, err := ReadColumns(CreateAgentDB(db), id)
	if err != nil {
		return nil, err
	}
	for _, col := range columns {
		outputCol := col.Json()
		var outputCards []types.CardJson
		cards, err := GetCardsByColumnId(db, col.Id)

		if err != nil {
			return nil, err
		}
		for _, card := range cards {
			outputCard := card.Json()

			tags, err := GetTagsByCard(db, card.Id)
			if err != nil {
				return nil, err
			}
			for _, tag := range tags {
				if tag.Id != "" {
					outputCard.TagIds = append(outputCard.TagIds, tag.Id)
				}
			}

			outputCards = append(outputCards, *outputCard)
		}
		outputCol.Cards = outputCards
		output.Columns = append(output.Columns, *outputCol)
	}
	return &output, nil
}
