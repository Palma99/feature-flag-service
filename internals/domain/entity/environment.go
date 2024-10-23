package domain

type Environment struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProjectID  string `json:"project_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Flags      []Flag `json:"flags"`
}
