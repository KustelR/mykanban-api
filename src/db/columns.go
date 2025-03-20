package db_driver

import (
	"context"
	"database/sql"
	"fmt"
	"types"
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
	stmt, err := agent.Prepare("DELETE FROM Columns WHERE id = ?;")
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
	fmt.Println(oldCol.ProjectId, oldCol.Order)
	stmtP.Exec(oldCol.ProjectId, oldCol.Order)
	stmt.Exec(id)

	err = tx.Commit()
	if err != nil {
		return nil
	}
	return nil
}

func AddColumns(agent *Agent, projectId string, columns *[]types.ColumnJson) error {
	stmt, err := agent.Prepare(`CALL add_column(?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var colErr error
out:
	for _, col := range *columns {
		_, err := stmt.Exec(projectId, col.Id, col.Name, col.Order)
		if err != nil {
			colErr = err
			break out
		}
		colErr = addCards(agent, col.Id, &col.Cards)
		if colErr != nil {
			break out
		}
	}
	if colErr != nil {
		return colErr
	}
	return nil
}
