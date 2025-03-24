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
	stmt, err := tx.Prepare(`
	INSERT Projects (
	name, 
	id
	) VALUES (
	 ?, # project name 
	 ? # project id
	);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(project.Name, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = AddTags(agent, id, &project.Tags)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = AddColumns(agent, id, &project.Columns)
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

func PostProject(db *sql.DB, id string, projectData *types.KanbanJson) error {
	agent := CreateAgentDB(db)
	stmt, err := agent.Prepare(`
	insert Projects (
    name,
    id    
) values (
 ?, # name
 ? # id
);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(projectData.Name, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return NoEffect{}
	}
	err = AddTags(agent, id, &projectData.Tags)
	if err != nil {
		return err
	}
	err = AddColumns(agent, id, &projectData.Columns)
	if err != nil {
		return err
	}

	return nil
}
