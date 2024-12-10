package domain

type Environment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ProjectID int64  `json:"project_id"`
	PublicKey string `json:"public_key"`
}

type EnvironmentWithFlags struct {
	Environment
	Flags []Flag `json:"flags"`
}
