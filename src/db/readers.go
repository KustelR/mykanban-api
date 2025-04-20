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

func readMultiRow(agent *Agent, id string, query string) ([]string, [][]sql.RawBytes, error) {
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

func GetCardsByColumnId(db *sql.DB, id string) ([]types.Card, error) {
	var outputCards []types.Card
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `CALL read_cards_by_column_id(?);`)

	for i := range values {
		row := values[i]
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

		meta, err := readMeta(columns, row)
		if err != nil {
			return nil, err
		}
		newCard.CreatedAt = meta.Created_at
		newCard.UpdatedAt = meta.Updated_at
		newCard.CreatedBy = meta.Created_by
		newCard.UpdatedBy = meta.Updated_by
		outputCards = append(outputCards, newCard)
	}
	return outputCards, err
}

func ReadProject(db *sql.DB, id string) (*types.Kanban, error) {
	agent := CreateAgentDB(db)
	var project types.Kanban
	columns, values, err := readOneRow(agent, id, "CALL read_project(?);")
	if err != nil {
		return nil, err
	}
	for i, col := range values {
		switch columns[i] {
		case "id":
			project.Id = string(col)
		case "name":
			project.Name = string(col)
		}
	}
	meta, err := readMeta(columns, values)
	if err != nil {
		return nil, err
	}
	project.Created_At = meta.Created_at
	project.Updated_At = meta.Updated_at
	project.Created_By = meta.Created_by
	project.Updated_By = meta.Updated_by
	return &project, nil
}

func GetColumn(agent *Agent, id string) (*types.Column, error) {
	columns, values, err := readOneRow(agent, id, "SELECT * FROM ProjectColumns WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	var column types.Column
	for i, col := range values {
		switch columns[i] {
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

	meta, err := readMeta(columns, values)
	if err != nil {
		return nil, err
	}
	column.CreatedAt = meta.Created_at
	column.UpdatedAt = meta.Updated_at
	column.CreatedBy = meta.Created_by
	column.UpdatedBy = meta.Updated_by
	return &column, nil
}
func ReadColumns(agent *Agent, projectId string) ([]types.Column, error) {
	var outputColumns []types.Column
	columns, values, err := readMultiRow(agent, projectId, `CALL read_columns_by_project_id(?);`)
	for i := range values {
		row := values[i]
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
		meta, err := readMeta(columns, row)
		if err != nil {
			return nil, err
		}
		newColumn.CreatedAt = meta.Created_at
		newColumn.UpdatedAt = meta.Updated_at
		newColumn.CreatedBy = meta.Created_by
		newColumn.UpdatedBy = meta.Updated_by
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

type Metadata struct {
	Created_at int
	Updated_at int
	Created_by string
	Updated_by string
}

func readMeta(columns []string, data []sql.RawBytes) (*Metadata, error) {
	var result Metadata
	for idx, val := range data {
		switch columns[idx] {
		case "created_at":
			data, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, err
			}
			result.Created_at = data
		case "updated_at":
			data, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, err
			}
			result.Updated_at = data
		case "created_by":
			result.Created_by = string(val)
		case "updated_by":
			result.Updated_by = string(val)
		}
	}
	return &result, nil
}

func GetTagsByCard(db *sql.DB, id string) ([]types.Tag, error) {
	var outputTags []types.Tag
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `CALL read_tags_by_card_id(?);`)
	for i := range values {
		row := values[i]
		rowLength := len(row)
		var newTag types.Tag
		for j := 0; j < rowLength; j++ {
			col := row[j]
			switch columns[j] {
			case "id":
				newTag.Id = string(col)
			case "name":
				newTag.Name = string(col)
			case "color":
				newTag.Color = string(col)
			case "project_id":
				newTag.ProjectId = string(col)
			}
		}
		meta, err := readMeta(columns, row)
		if err != nil {
			return nil, err
		}
		newTag.CreatedAt = meta.Created_at
		newTag.UpdatedAt = meta.Updated_at
		newTag.CreatedBy = meta.Created_by
		newTag.UpdatedBy = meta.Updated_by
		outputTags = append(outputTags, newTag)
	}
	return outputTags, err
}

func GetTagsByProject(db *sql.DB, id string) ([]types.Tag, error) {
	var outputTags []types.Tag
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `CALL read_tags_by_project_id(?);`)
	if err != nil {
		return nil, err
	}
	for i := range values {
		row := values[i]
		rowLength := len(row)
		var newTag types.Tag
		if values[i] == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := row[j]
			switch columns[j] {
			case "id":
				newTag.Id = string(col)
			case "name":
				newTag.Name = string(col)
			case "color":
				newTag.Color = string(col)
			}
		}
		meta, err := readMeta(columns, row)
		if err != nil {
			return nil, err
		}
		newTag.CreatedAt = meta.Created_at
		newTag.UpdatedAt = meta.Updated_at
		newTag.CreatedBy = meta.Created_by
		newTag.UpdatedBy = meta.Updated_by
		outputTags = append(outputTags, newTag)
	}
	return outputTags, err
}
