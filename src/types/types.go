package types

type Tag struct {
	Id         string
	ProjectId  string
	Name       string
	Color      string
	Created_At int
	Updated_At int
	Created_By string
	Updated_By string
}

type TagJson struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Created_At int    `json:"created_at"`
	Updated_At int    `json:"updated_at"`
	Created_By string `json:"created_by"`
	Updated_By string `json:"updated_by"`
}

func (t *Tag) Json() *TagJson {
	return &TagJson{t.Id, t.Name, t.Color, t.Created_At, t.Updated_At, t.Created_By, t.Updated_By}
}

type Card struct {
	Id          string
	ColumnId    string
	Name        string
	Order       int
	Description string
	Created_At  int
	Updated_At  int
	Created_By  string
	Updated_By  string
}
type CardJson struct {
	Id          string   `json:"id"`
	ColumnId    string   `json:"columnId"`
	Name        string   `json:"name"`
	Order       int      `json:"order"`
	Description string   `json:"description"`
	TagIds      []string `json:"tagIds"`
	Created_At  int      `json:"created_at"`
	Updated_At  int      `json:"updated_at"`
	Created_By  string   `json:"created_by"`
	Updated_By  string   `json:"updated_by"`
}

func (c *Card) Json() *CardJson {
	var tagIds [0]string
	return &CardJson{c.Id, c.ColumnId, c.Name, c.Order, c.Description, tagIds[:], c.Created_At, c.Updated_At, c.Created_By, c.Updated_By}
}

type Column struct {
	Id         string
	Name       string
	Order      int
	ProjectId  string
	Created_At int
	Updated_At int
	Created_By string
	Updated_By string
}
type ColumnJson struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Order      int        `json:"order"`
	Cards      []CardJson `json:"cards"`
	Created_At int        `json:"created_at"`
	Updated_At int        `json:"updated_at"`
	Created_By string     `json:"created_by"`
	Updated_By string     `json:"updated_by"`
}

func (c *Column) Json() *ColumnJson {
	var cards [0]CardJson
	return &ColumnJson{c.Id, c.Name, c.Order, cards[:], c.Created_At, c.Updated_At, c.Created_By, c.Updated_By}
}

type Kanban struct {
	Name       string
	Id         string
	Created_At int
	Updated_At int
	Created_By string
	Updated_By string
}
type KanbanJson struct {
	Name       string       `json:"name"`
	Columns    []ColumnJson `json:"columns"`
	Tags       []TagJson    `json:"tags"`
	Created_At int          `json:"created_at"`
	Updated_At int          `json:"updated_at"`
	Created_By string       `json:"created_by"`
	Updated_By string       `json:"updated_by"`
}

func (k *Kanban) Json() *KanbanJson {
	var columns [0]ColumnJson
	var tags [0]TagJson
	return &KanbanJson{k.Name, columns[:], tags[:], k.Created_At, k.Updated_At, k.Created_By, k.Updated_By}
}
