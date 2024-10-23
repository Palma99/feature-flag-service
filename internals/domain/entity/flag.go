package domain

type Flag struct {
	ID          int         `json:"id"`
	Name        string      `json:"value"`
	Enabled     bool        `json:"enabled"`
	Environment Environment `json:"environment"`
}

func (f *Flag) UpdateEnabled(enabled bool) *Flag {
	return &Flag{
		ID:          f.ID,
		Name:        f.Name,
		Enabled:     enabled,
		Environment: f.Environment,
	}
}
