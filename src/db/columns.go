package db_driver

import (
	"context"
	"database/sql"
	"types"
	"utils"
)

func UpdateColumnData(db *sql.DB, column *types.Column) error {
	tx, err := db.BeginTx(context.Background(), nil)
	agent := CreateAgentTX(tx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := agent.Prepare("CALL update_column_data(?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(column.Id, column.Name, column.Order)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func DeleteColumn(db *sql.DB, id string) error {
	tx, err := db.BeginTx(context.Background(), nil)
	agent := CreateAgentTX(tx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := agent.Prepare("DELETE FROM ProjectColumns WHERE id = ?;")
	if err != nil {
		tx.Rollback()
		return err
	}
	oldCol, err := GetColumn(agent, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	stmtP, err := agent.Prepare("CALL pop_column_reorder(?, ?)")
	if err != nil {
		return nil
	}
	defer stmtP.Close()
	stmtP.Exec(oldCol.ProjectId, oldCol.Order)
	stmt.Exec(id)

	err = tx.Commit()
	if err != nil {
		return nil
	}
	return nil
}

func CreateColumns(agent *Agent, projectId string, columns []types.ColumnJson) ([]types.ColumnJson, error) {
	stmt, err := agent.Prepare(`CALL create_column(?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var colErr error

	newCols := make([]types.ColumnJson, len(columns))

out:
	for idx, col := range columns {
		id := utils.GetUUID()
		changedCol := col
		changedCol.Id = id

		dbColNames, data, err := readOneRow(agent, projectId, "SELECT max(draw_order) FROM ProjectColumns WHERE project_id= ?;")
		if err != nil {
			colErr = err
			break out
		}
		drawOrder, err := GetMaxDrawOrder(dbColNames, data)
		if err != nil {
			colErr = err
			break out
		}
		changedCol.Order = drawOrder + 1
		_, err = stmt.Exec(projectId, id, changedCol.Name, changedCol.Order, "placeholder")
		if err != nil {
			colErr = err
			break out
		}
		cards, err := CreateCards(agent, id, &col.Cards)
		if err != nil {
			colErr = err
			break out
		}
		changedCol.Cards = cards
		newCols[idx] = changedCol

	}
	if colErr != nil {
		return nil, colErr
	}
	return newCols, nil
}
