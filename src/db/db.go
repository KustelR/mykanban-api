package db_driver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"types"

	_ "github.com/go-sql-driver/mysql"
)

func GetDb(connString string) *sql.DB {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(fmt.Errorf("can't connect to database"))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("can't connect to database"))
	}
	fmt.Println("Connected to MySql DB")
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	return db
}

type NotFoundError struct {
	thing string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s was not found", e.thing)
}

func readOneRow(db *sql.DB, id string, query string) ([]string, []sql.RawBytes, error) {
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
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	if !rows.Next() {
		return nil, nil, NotFoundError{fmt.Sprintf("project with id %s", id)}
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

func GetCard(db *sql.DB, id string) (*types.Card, error) {
	columns, values, err := readOneRow(db, id, `SELECT * FROM Cards WHERE id = ?`)
	if err != nil {
		return nil, err
	}
	var card types.Card
	for i, col := range values {
		switch columns[i] {
		case "id":
			card.Id = string(col)
		case "column_id":
			card.Id = string(col)
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

func readProject(db *sql.DB, id string) (*types.Kanban, error) {
	var project types.Kanban
	columns, values, err := readOneRow(db, id, "select * from Projects where id=?;")
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

func GetTagsByCard(db *sql.DB, id string) ([]types.Tag, error) {
	var outputTags []types.Tag
	columns, values, err := readMultiRow(db, id, `select Tags.* from Tags join
CardsTags on Tags.id = CardsTags.tag_id join
Cards on Cards.id = CardsTags.card_id
where Cards.id=?;`)
	for i := range values {
		rowLength := len(values[i])
		var newTag types.Tag
		if values[i] == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := values[i][j]
			switch columns[j] {
			case "id":
				newTag.Id = string(col)
			case "name":
				newTag.Name = string(col)
			case "color":
				newTag.Color = string(col)
			}
		}
		outputTags = append(outputTags, newTag)
	}
	return outputTags, err
}

func GetTagsByProject(db *sql.DB, id string) ([]types.Tag, error) {
	var outputTags []types.Tag
	columns, values, err := readMultiRow(db, id, `select * from tags where project_id=?`)
	for i := range values {
		rowLength := len(values[i])
		var newTag types.Tag
		if values[i] == nil {
			continue
		}
		for j := 0; j < rowLength; j++ {
			col := values[i][j]
			switch columns[j] {
			case "id":
				newTag.Id = string(col)
			case "name":
				newTag.Name = string(col)
			case "color":
				newTag.Color = string(col)
			}
		}
		outputTags = append(outputTags, newTag)
	}
	return outputTags, err
}

type NoEffect struct{}

func (NoEffect) Error() string {
	return "no rows were affected by "
}

func createTag(db *sql.DB, projectId string, tagData *types.TagJson) error {
	stmt, err := db.Prepare(`
	insert tags (
    id,
    project_id,
    name,
    color
) values (
    ?, # id
    ?, # project id
    ?, # name
    ? # color
);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(tagData.Id, projectId, tagData.Name, tagData.Color)
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
	return nil
}

func addCard(db *sql.DB, cardData *types.CardJson) error {
	stmt, err := db.Prepare(`
insert cards (
    id,
    column_id,
    name,
    description,
	draw_order
) values (
    ?, # id
    ?, # associated column id
    ?, # name
    ?, # card description
	? # draw order
);
`)
	if err != nil {
		return err
	}
	stmtCT, err := db.Prepare(`
	insert CardsTags (card_id, tag_id) values (

    ?, # card_id
    ? # tag_id
);`)
	defer stmt.Close()
	if err != nil {
		return err
	}
	res, err := stmt.Exec(cardData.Id, cardData.ColumnId, cardData.Name, cardData.Description, cardData.Order)
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
	var cycleErr error
	for _, tagId := range cardData.TagIds {
		_, cycleErr = stmtCT.Exec(cardData.Id, tagId)
	}
	if cycleErr != nil {
		return cycleErr
	}
	return nil
}

func createColumn(db *sql.DB, projectId string, colData *types.ColumnJson) error {
	stmt, err := db.Prepare(`
insert columns (
    id,
    project_id,
    name,
	draw_order
) values (
    ?, # id
    ?, # associated project id
    ?, # name
	? # draw order
);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(colData.Id, projectId, colData.Name, colData.Order)
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
	for _, card := range colData.Cards {
		addCard(db, &card)
	}
	return nil
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
				outputCard.TagIds = append(outputCard.TagIds, tag.Id)
			}

			outputCards = append(outputCards, *outputCard)
		}
		outputCol.Cards = outputCards
		output.Columns = append(output.Columns, *outputCol)
	}
	return &output, nil
}

func PostProject(db *sql.DB, id string, projectData *types.KanbanJson) error {
	stmt, err := db.Prepare(`
	insert projects (
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
	for _, tag := range projectData.Tags {
		err = createTag(db, id, &tag)
		if err != nil {
			return err
		}
	}
	for _, col := range projectData.Columns {
		err = createColumn(db, id, &col)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateProject(db *sql.DB, ctx context.Context, id string, project *types.KanbanJson) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	UPDATE Projects
    SET 
	name = ? # project name 
	WHERE 
	id = ? # project id;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(project.Name, id)
	err = commitTags(tx, id, &project.Tags)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = commitColumns(tx, id, &project.Columns)
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

func commitTags(tx *sql.Tx, projectId string, tags *[]types.TagJson) error {
	stmt, err := tx.Prepare(`
	insert tags (
    id,
    project_id,
    name,
    color
) values (
    ?, # id
    ?, # project id
    ?, # name
    ? # color
) ON DUPLICATE KEY UPDATE
 id = id,
 project_id = project_id,
 name=name,
 color=color;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var tagErr error
	for _, tag := range *tags {
		_, err := stmt.Exec(tag.Id, projectId, tag.Name, tag.Color)
		if err != nil {
			tagErr = err
		}
	}
	if tagErr != nil {
		return tagErr
	}
	return nil
}

func commitCards(tx *sql.Tx, columnId string, cards *[]types.CardJson) error {
	stmt, err := tx.Prepare(`
insert cards (
    id,
    column_id,
    name,
    description,
	draw_order
) values (
    ?, # id
    ?, # associated column id
    ?, # name
    ?, # card description
	? # draw order
) ON DUPLICATE KEY UPDATE
 id = id,
 column_id = column_id,
 name = name,
 description = description,
 draw_order = draw_order;
`)
	if err != nil {
		return err
	}
	stmtCT, err := tx.Prepare(`
	insert CardsTags (card_id, tag_id) values (

    ?, # card_id
    ? # tag_id
);`)
	if err != nil {
		return err
	}
	var cardErr error
	for _, card := range *cards {
		_, err := stmt.Exec(card.Id, columnId, card.Name, card.Description, card.Order)
		if err != nil {
			cardErr = err
		}
		for _, tagId := range card.TagIds {
			stmtCT.Exec(card.Id, tagId)
		}
	}
	if cardErr != nil {
		return cardErr
	}
	return nil
}

func commitColumns(tx *sql.Tx, projectId string, columns *[]types.ColumnJson) error {
	stmt, err := tx.Prepare(`
insert columns (
    id,
    project_id,
    name,
	draw_order
) values (
    ?, # id
    ?, # associated project id
    ?, # name
	? # draw order
) ON DUPLICATE KEY UPDATE id=id, project_id=project_id, name=name, draw_order=draw_order;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var colErr error
	for _, col := range *columns {
		_, err := stmt.Exec(col.Id, projectId, col.Name, col.Order)
		if err != nil {
			colErr = err
		}
		commitCards(tx, col.Id, &col.Cards)
	}
	if colErr != nil {
		return colErr
	}
	return nil
}
