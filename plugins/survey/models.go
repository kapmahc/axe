package survey

import "time"

// Form form
type Form struct {
	tableName struct{} `sql:"survey_forms"`
	ID        uint
	Deadline  time.Time
	Title     string
	Body      string
	Type      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Expire expire?
func (p *Form) Expire() bool {
	return time.Now().After(p.Deadline)
}

// Field field
type Field struct {
	tableName struct{} `sql:"survey_fields"`
	ID        uint
	Name      string
	Body      string
	Type      string
	Label     string
	Value     string
	SortOrder int
	Form      Form
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Record record
type Record struct {
	tableName struct{} `sql:"survey_records"`
	ID        uint
	Value     string
	Form      Form
	UpdatedAt time.Time
	CreatedAt time.Time
}
