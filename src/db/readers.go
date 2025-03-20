package db_driver

import (
	"database/sql"
	"fmt"
	"strconv"
	"types"
)

func readOneRow(agent *Agent, id string, query string) ([]string, []sql.RawBytes, error) {
	stmt, err := agent.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	if !rows.Next() {
		return nil, nil, NotFoundError{fmt.Sprintf("item with id %s", id), &query}
	}
	err = rows.Scan(scanArgs...)
	if err != nil {
		return nil, nil, err
	}
	return columns, values, nil
}

func readMultiRow(db *sql.DB, id string, query string) ([]string, [][]sql.RawBytes, error) {
	var err error
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	output := make([][]sql.RawBytes, len(columns))

	for rows.Next() {
		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		output = append(output, values)
	}
	if err != nil {
		return nil, nil, err
	}
	if len(columns) == 0 {
		return nil, nil, fmt.Errorf(" id: %s was not found", id)
	}

	return columns, output, nil
}

func GetCard(agent *Agent, id string) (*types.Card, error) {
	columns, values, err := readOneRow(agent, id, `SELECT * FROM Cards WHERE id = ?;`)
	if err != nil {
		return nil, err
	}
	var card types.Card
	for i, col := range values {
		switch columns[i] {
		case "id":
			card.Id = string(col)
		case "column_id":
			card.ColumnId = string(col)
		case "name":
			card.Name = string(col)
		case "description":
			card.Name = string(col)
		case "draw_order":
			val, err := strconv.Atoi(string(col))
			if err != nil {
				return nil, err
			}
			card.Order = val
		}
	}
	return &card, nil
}

func GetCards(db *sql.DB, id string) ([]types.Card, error) {
	var outputCards []types.Card
	columns, values, err := readMultiRow(db, id, `select Cards.* from Columns join
Cards on Cards.column_id = Columns.id
where Columns.id=?;`)
	for i := range values {
		rowLength := len(values[i])
		var newCard types.Card
		if values[i] == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := values[i][j]
			switch columns[j] {
			case "id":
				newCard.Id = string(col)
			case "column_id":
				newCard.ColumnId = string(col)
			case "draw_order":
				val, err := strconv.Atoi(string(col))
				if err != nil {
					return nil, err
				}
				newCard.Order = val
			case "name":
				newCard.Name = string(col)
			case "description":
				newCard.Description = string(col)
			}
		}
		outputCards = append(outputCards, newCard)
	}
	return outputCards, err
}

func GetProject(db *sql.DB, id string) (*types.KanbanJson, error) {
	var output types.KanbanJson
	project, err := readProject(db, id)
	if err != nil {
		return nil, err
	}
	output.Name = project.Name
	projectTags, err := GetTagsByProject(db, id)
	if err != nil {
		return nil, err
	}
	for _, tag := range projectTags {
		outputTag := tag.Json()
		output.Tags = append(output.Tags, *outputTag)
	}
	columns, err := readColumns(db, id)
	if err != nil {
		return nil, err
	}
	for _, col := range columns {
		outputCol := col.Json()
		var outputCards []types.CardJson
		cards, err := GetCards(db, col.Id)

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

func readProject(db *sql.DB, id string) (*types.Kanban, error) {
	agent := CreateAgentDB(db)
	var project types.Kanban
	columns, values, err := readOneRow(agent, id, "select * from Projects where id=?;")
	for i, col := range values {
		switch columns[i] {
		case "id":
			project.Id = string(col)
		case "name":
			project.Name = string(col)
		}
	}
	return &project, err
}
func readColumns(db *sql.DB, id string) ([]types.Column, error) {
	var outputColumns []types.Column
	columns, values, err := readMultiRow(db, id, `select columns.* from projects join
columns on columns.project_id = Projects.id
where Projects.id=?;`)
	for i := range values {
		rowLength := len(values[i])
		var newColumn types.Column
		if values[i] == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := values[i][j]
			switch columns[j] {
			case "id":
				newColumn.Id = string(col)
			case "project_id":
				newColumn.ProjectId = string(col)
			case "name":
				newColumn.Name = string(col)
			case "draw_order":
				val, err := strconv.Atoi(string(col))
				if err != nil {
					return nil, err
				}
				newColumn.Order = val
			}
		}
		outputColumns = append(outputColumns, newColumn)
	}
	return outputColumns, err
}
