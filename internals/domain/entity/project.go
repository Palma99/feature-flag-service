package domain

import "errors"

type Project struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Environments []Environment `json:"environments"`
	OwnerId      int           `json:"owner_id"`
}

func NewProject(name string, ownerId int) *Project {
	return &Project{
		Name:    name,
		OwnerId: ownerId,
	}
}

func (p *Project) CanCreateEnvironment(env Environment) error {
	for _, e := range p.Environments {
		if e.Name == env.Name {
			return errors.New("environment with the same name already exists")
		}
	}
	return nil
}

type ProjectWithMembers struct {
	Project
	Members []User `json:"members"`
}
