package db_driver

import (
	"context"
	"database/sql"
	"types"
)

func UpdateProject(db *sql.DB, ctx context.Context, id string, project *types.KanbanJson) error {
	tx, err := db.BeginTx(ctx, nil)
	agent := CreateAgentTX(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("DELETE FROM Projects WHERE id = ?;", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	stmt, err := tx.Prepare(`CALL create_project(?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project.Name, id, "placeholder")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = CreateTag(agent, id, &project.Tags)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = CreateColumns(agent, id, &project.Columns)
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
	_, err = CreateTag(agent, id, &projectData.Tags)
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = CreateColumns(agent, id, &projectData.Columns)
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
