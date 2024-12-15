package domain

type Flag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	ProjectID int64  `json:"project_id"`
}

type FlagWithoutEnabled struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ProjectID int64  `json:"project_id"`
}
