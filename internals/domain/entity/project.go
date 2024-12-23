package domain

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

func (p *Project) EnvironmentWithSameNameAlreadyExists(env Environment) bool {
	for _, e := range p.Environments {
		if e.Name == env.Name {
			return true
		}
	}
	return false
}

type ProjectWithMembers struct {
	Project
	Members []User `json:"members"`
}

func (pm *ProjectWithMembers) HasMember(userId int) bool {
	for _, member := range pm.Members {
		if member.ID == userId {
			return true
		}
	}
	return false
}

func (pm *ProjectWithMembers) GetUserPermissions(userId int) []string {
	var permissions []string

	if pm.OwnerId == userId {
		permissions = append(permissions, PermissionCreateProjectEnvironment)
		permissions = append(permissions, PermissionDeleteFlag)
	}

	for _, member := range pm.Members {
		if member.ID == userId {
			permissions = append(permissions, PermissionCreateFlag)
			break
		}
	}

	return permissions
}

func (p *ProjectWithMembers) UserHasPermission(userId int, permission string) bool {
	userPermissions := p.GetUserPermissions(userId)

	for _, userPermission := range userPermissions {
		if userPermission == permission {
			return true
		}
	}

	return false
}
