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

func (i *FlagInteractor) UpdateFlagEnvironment(environmentId, userId int, flagsToUpdate []domain.Flag) error {
	projectDetails, err := i.projectRepository.GetProjectDetailsByEnvironmentId(environmentId)
	if err != nil || !projectDetails.HasMember(userId) {
		return errors.New("user is not allowed to update this flag")
	}

	if err := i.flagRepository.UpdateFlagEnvironment(environmentId, flagsToUpdate); err != nil {
		return err
	}

	return nil
}

func (i *FlagInteractor) GetProjectFlags(userId int, projectId int64) ([]domain.FlagWithoutEnabled, error) {
	projectWithMembers, err := i.projectRepository.GetProjectDetails(projectId)
	if err != nil || !projectWithMembers.HasMember(userId) {
		return nil, errors.New("not allowed")
	}

	flags, err := i.flagRepository.GetProjectFlags(projectId)
	if err != nil {
		return nil, err
	}

	return flags, nil
}
