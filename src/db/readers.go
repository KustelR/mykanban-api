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

func readMultiRow(agent *Agent, id string, query string) ([]string, *[][]sql.RawBytes, error) {
	var err error
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
	output := make([][]sql.RawBytes, 0)

	for rows.Next() {
		values := make([]sql.RawBytes, len(columns))
		dest := make([]any, len(values))
		for i := range values {
			dest[i] = &values[i]
		}
		err = rows.Scan(dest...)

		if err != nil {
			return nil, nil, err
		}

		output = copyAndAppend(output, values)
	}

	if len(columns) == 0 {
		return nil, nil, fmt.Errorf(" id: %s was not found", id)
	}
	return columns, &output, nil
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
			card.Description = string(col)
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
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `select * from Cards where column_id=?;`)

	for i := range *values {
		row := (*values)[i]
		rowLength := len(row)
		var newCard types.Card
		if row == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := row[j]
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
	columns, err := readColumns(CreateAgentDB(db), id)
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

func GetColumn(agent *Agent, id string) (*types.Column, error) {
	colNames, values, err := readOneRow(agent, id, "SELECT * FROM Columns WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	var column types.Column
	for i, col := range values {
		switch colNames[i] {
		case "id":
			column.Id = string(col)
		case "name":
			column.Name = string(col)
		case "project_id":
			column.ProjectId = string(col)
		case "draw_order":
			val, err := strconv.Atoi(string(col))
			if err != nil {
				return nil, err
			}
			column.Order = val
		}
	}
	return &column, nil
}
func readColumns(agent *Agent, projectId string) ([]types.Column, error) {
	var outputColumns []types.Column
	columns, values, err := readMultiRow(agent, projectId, `SELECT * FROM Columns WHERE project_id=?;`)
	for i := range *values {
		row := (*values)[i]
		rowLength := len(row)
		var newColumn types.Column
		if row == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := row[j]
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

func copyAndAppend(sl [][]sql.RawBytes, item []sql.RawBytes) [][]sql.RawBytes {
	newSlice := make([][]sql.RawBytes, len(sl)+1)
	copy(newSlice, sl)
	itemCopy := make([]sql.RawBytes, len(item))
	for i := range itemCopy {
		entryCopy := make(sql.RawBytes, len(item[i]))
		copy(entryCopy, item[i])
		itemCopy[i] = entryCopy
	}
	newSlice[len(newSlice)-1] = itemCopy
	return newSlice
}
