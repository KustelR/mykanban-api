package types

type Tag struct {
	Id        string
	ProjectId string
	Name      string
	Color     string
}
type TagJson struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (t *Tag) Json() *TagJson {
	return &TagJson{t.Id, t.Name, t.Color}
}

type Card struct {
	Id          string
	ColumnId    string
	Name        string
	Order       int
	Description string
}
type CardJson struct {
	Id          string   `json:"id"`
	ColumnId    string   `json:"columnId"`
	Name        string   `json:"name"`
	Order       int      `json:"order"`
	Description string   `json:"description"`
	TagIds      []string `json:"tagIds"`
}

func (c *Card) Json() *CardJson {
	var tagIds [0]string
	return &CardJson{c.Id, c.ColumnId, c.Name, c.Order, c.Description, tagIds[:]}
}

type Column struct {
	Id        string
	Name      string
	Order     int
	ProjectId string
}
type ColumnJson struct {
	Id    string     `json:"id"`
	Name  string     `json:"name"`
	Order int        `json:"order"`
	Cards []CardJson `json:"cards"`
}

func (c *Column) Json() *ColumnJson {
	var cards [0]CardJson
	return &ColumnJson{c.Id, c.Name, c.Order, cards[:]}
}

type Kanban struct {
	Name string
	Id   string
}
type KanbanJson struct {
	Name    string       `json:"name"`
	Columns []ColumnJson `json:"columns"`
	Tags    []TagJson    `json:"tags"`
}

func (c *Kanban) Json() *KanbanJson {
	var columns [0]ColumnJson
	var tags [0]TagJson
	return &KanbanJson{c.Name, columns[:], tags[:]}
}
