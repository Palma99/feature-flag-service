package domain

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string
	Nickname string `json:"nickname"`
}

type LoggedUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

func (lu LoggedUser) CanCreateProjectEnvironment(project Project) bool {
	return lu.ID == project.OwnerId
}
