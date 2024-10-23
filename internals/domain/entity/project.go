package domain

type Project struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Environments []Environment `json:"environments"`
}
