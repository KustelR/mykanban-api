package db_driver

import (
	"database/sql"
	"types"
)

func AddTags(agent *Agent, projectId string, tags *[]types.TagJson) error {
	stmt, err := agent.Prepare(`
	insert Tags (
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

func GetTagsByCard(db *sql.DB, id string) ([]types.Tag, error) {
	var outputTags []types.Tag
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `select Tags.* from Tags join
CardsTags on Tags.id = CardsTags.tag_id join
Cards on Cards.id = CardsTags.card_id
where Cards.id=?;`)
	for i := range values {
		rowLength := len(values[i])
		var newTag types.Tag
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
	columns, values, err := readMultiRow(CreateAgentDB(db), id, `select * from Tags where project_id=?;`)
	if err != nil {
		return nil, err
	}
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
