package usecase

import (
	"errors"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	domain "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type FlagInteractor struct {
	flagRepository        repository.FlagRepository
	environmentRepository repository.EnvironmentRepository
	keyService            *services.KeyService
	projectRepository     repository.ProjectRepository
}

func NewFlagInteractor(
	flagRepository repository.FlagRepository,
	environmentRepository repository.EnvironmentRepository,
	keyService *services.KeyService,
	projectRepository repository.ProjectRepository,
) *FlagInteractor {
	return &FlagInteractor{
		flagRepository,
		environmentRepository,
		keyService,
		projectRepository,
	}
}

func (i *FlagInteractor) CreateFlag(projectId int64, userId int, flagName string) error {
	projectWithMembers, err := i.projectRepository.GetProjectDetails(projectId)
	if err != nil || !projectWithMembers.HasMember(userId) {
		return errors.New("user is not allowed to create this flag")
	}

	if flagName == "" {
		return errors.New("flag name is required")
	}

	flag := &domain.Flag{
		Name:      flagName,
		ProjectID: projectId,
		Enabled:   false,
	}

	if err := i.flagRepository.CreateFlag(flag); err != nil {
		return err
	}

	return nil
}

// func (i *FlagInteractor) UpdateFlagValue(key string, flagId int, value bool) error {
// 	if i.keyService.IsPublicKey(key) {
// 		return errors.New("unauthorized")
// 	}

// 	env, err := i.environmentRepository.GetEnvironmentBySecretKey(key)
// 	if err != nil {
// 		fmt.Println("cannot fetch environment by secret key")
// 		return err
// 	}

// 	flag, err := i.flagRepository.GetFlagInEnvironmentById(env.ID, flagId)
// 	if err != nil {
// 		return err
// 	}

// 	updatedFlag := flag.UpdateEnabled(value)

// 	if err := i.flagRepository.UpdateEnvironmentFlagValue(env.ID, updatedFlag); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (i *FlagInteractor) GetAllFlagsByEnvironmentKey(key string) ([]entity.Flag, error) {
// 	if !i.keyService.IsPublicKey(key) {
// 		return nil, errors.New("secret key is not supported yet")
// 	}

// 	env, err := i.environmentRepository.GetEnvironmentByPublicKey(key)
// 	if err != nil {
// 		return nil, errors.New("error while fetching flags")
// 	}

// 	if flags, err := i.flagRepository.GetAllFlagsByEnvironmentID(env.ID); err != nil {
// 		return nil, errors.New("error while fetching flags")
// 	} else {
// 		return flags, nil
// 	}
// }
