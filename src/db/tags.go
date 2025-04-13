package db_driver

import (
	"types"
	"utils"
)

func CreateTags(agent *Agent, projectId string, tags *[]types.TagJson) ([]types.TagJson, error) {
	stmt, err := agent.Prepare(`CALL create_tag(?, ?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var tagErr error
	createdTags := make([]types.TagJson, len(*tags))
	for idx, tag := range *tags {
		newTag := tag
		newTag.Id = utils.GetUUID()
		createdTags[idx] = newTag
		_, err := stmt.Exec(projectId, newTag.Id, tag.Name, tag.Color, "placeholder")
		if err != nil {
			tagErr = err
		}
	}
	if tagErr != nil {
		return nil, tagErr
	}
	return createdTags, nil
}

type NoEffect struct{}

func (NoEffect) Error() string {
	return "no rows were affected by "
}
