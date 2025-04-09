package db_driver

import "database/sql"

func CreateAgentDB(db *sql.DB) *Agent {
	agent := Agent{db, nil}
	return &agent
}
func CreateAgentTX(tx *sql.Tx) *Agent {
	agent := Agent{nil, tx}
	return &agent
}

type Agent struct {
	db *sql.DB
	tx *sql.Tx
}

func (a *Agent) Prepare(query string) (*sql.Stmt, error) {
	if a.db != nil {
		return a.db.Prepare(query)
	} else {
		return a.tx.Prepare(query)
	}
}

func (a *Agent) Exec(query string, args ...any) (sql.Result, error) {
	if a.db != nil {
		return a.db.Exec(query, args)
	} else {
		return a.tx.Exec(query, args)
	}
}
