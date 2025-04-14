package types

type Tag struct {
	Id        string
	ProjectId string
	Name      string
	Color     string
	CreatedAt int
	UpdatedAt int
	CreatedBy string
	UpdatedBy string
}

type TagJson struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

func (t *Tag) Json() *TagJson {
	return &TagJson{t.Id, t.Name, t.Color, t.CreatedAt, t.UpdatedAt, t.CreatedBy, t.UpdatedBy}
}

type Card struct {
	Id          string
	ColumnId    string
	Name        string
	Order       int
	Description string
	CreatedAt   int
	UpdatedAt   int
	CreatedBy   string
	UpdatedBy   string
}
type CardJson struct {
	Id          string   `json:"id"`
	ColumnId    string   `json:"columnId"`
	Name        string   `json:"name"`
	Order       int      `json:"order"`
	Description string   `json:"description"`
	TagIds      []string `json:"tagIds"`
	CreatedAt   int      `json:"createdAt"`
	UpdatedAt   int      `json:"updatedAt"`
	CreatedBy   string   `json:"createdBy"`
	UpdatedBy   string   `json:"updatedBy"`
}

func (c *Card) Json() *CardJson {
	var tagIds [0]string
	return &CardJson{c.Id, c.ColumnId, c.Name, c.Order, c.Description, tagIds[:], c.CreatedAt, c.UpdatedAt, c.CreatedBy, c.UpdatedBy}
}

type Column struct {
	Id        string
	Name      string
	Order     int
	ProjectId string
	CreatedAt int
	UpdatedAt int
	CreatedBy string
	UpdatedBy string
}
type ColumnJson struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Order     int        `json:"order"`
	Cards     []CardJson `json:"cards"`
	CreatedAt int        `json:"createdAt"`
	UpdatedAt int        `json:"updatedAt"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy string     `json:"updatedBy"`
}

func (c *Column) Json() *ColumnJson {
	var cards [0]CardJson
	return &ColumnJson{c.Id, c.Name, c.Order, cards[:], c.CreatedAt, c.UpdatedAt, c.CreatedBy, c.UpdatedBy}
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
	Name      string       `json:"name"`
	Columns   []ColumnJson `json:"columns"`
	Tags      []TagJson    `json:"tags"`
	CreatedAt int          `json:"createdAt"`
	UpdatedAt int          `json:"updatedAt"`
	CreatedBy string       `json:"createdBy"`
	UpdatedBy string       `json:"updatedBy"`
}

func (k *Kanban) Json() *KanbanJson {
	var columns [0]ColumnJson
	var tags [0]TagJson
	return &KanbanJson{k.Name, columns[:], tags[:], k.Created_At, k.Updated_At, k.Created_By, k.Updated_By}
}
