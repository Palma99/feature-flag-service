package domain

type Flag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	ProjectID int64  `json:"project_id"`
}

func (f *Flag) UpdateEnabled(enabled bool) *Flag {
	return &Flag{
		ID:        f.ID,
		Name:      f.Name,
		ProjectID: f.ProjectID,
	}
}
