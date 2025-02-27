package db_driver

import (
	"database/sql"
	"fmt"
	"strconv"
	"types"

	_ "github.com/go-sql-driver/mysql"
)

func GetDb(connString string) *sql.DB {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(fmt.Errorf("can't connect to database"))
	}

	fmt.Println("Connected to MySql DB")
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("can't connect to database"))
	}
	return db
}

type NotFoundError struct {
	thing string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s was not found", e.thing)
}

func readOneRow(db *sql.DB, id string, query string) ([]string, []sql.RawBytes, error) {
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

func addCard(db *sql.DB, position int, cardData *types.CardJson) error {
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
	lastOrder, err := getLastDrawOrder(db, "Cards", fmt.Sprintf("column_id = \"%s\"", cardData.ColumnId))
	if err != nil {
		return nil
	}
	res, err := stmt.Exec(cardData.Id, cardData.ColumnId, cardData.Name, cardData.Description, lastOrder+1)
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
	for _, tagId := range cardData.TagIds {
		stmtCT.Exec(cardData.Id, tagId)
	}
	return nil
}

func createColumn(db *sql.DB, projectId string, position int, colData *types.ColumnJson) error {
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
	var order int
	if position == -1 {
		order, err = getLastDrawOrder(db, "Columns", fmt.Sprintf("project_id = \"%s\"", projectId))
		if err != nil {
			return err
		}
		order++
	} else {
		order = position
		err = createDrawSpace(db, position, "Columns", fmt.Sprintf("project_id = \"%s\"", projectId))
		if err != nil {
			return err
		}
	}

	res, err := stmt.Exec(colData.Id, projectId, colData.Name, order)
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
		addCard(db, card.Order, &card)
	}
	return nil
}

func getLastDrawOrder(db *sql.DB, table string, condition string) (int, error) {
	stmt, err := db.Prepare(fmt.Sprintf(` 	SELECT MAX(draw_order) FROM %s WHERE %s;`, table, condition))
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return -1, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return -1, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	if !rows.Next() {
		return -1, fmt.Errorf("draw order in table: %s was not found", table)
	}
	err = rows.Scan(scanArgs...)
	if err != nil {
		return -1, err
	}
	val, err := strconv.Atoi(string(values[0]))
	if err != nil {
		return 0, nil
	}
	return val, nil
}

func createDrawSpace(db *sql.DB, position int, table string, condition string) error {
	stmt, err := db.Prepare(fmt.Sprintf(`UPDATE %s SET draw_order = draw_order + 1 WHERE draw_order >= ? AND %s;`, table, condition))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(position)
	if err != nil {
		return nil
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
	for _, col := range projectData.Columns {
		err = createColumn(db, id, col.Order, &col)
		if err != nil {
			return err
		}
	}
	for _, tag := range projectData.Tags {
		err = createTag(db, id, &tag)
		if err != nil {
			return err
		}
	}
	return nil
}
