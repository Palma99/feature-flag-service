package usecase

import (
	"errors"
	"fmt"

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

	if !projectWithMembers.UserHasPermission(userId, domain.PermissionCreateFlag) {
		return errors.New("user is not allowed to create flags on this project")
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

func (i *FlagInteractor) DeleteFlag(userId, flagId int) error {
	projectWithMembers, err := i.projectRepository.GetProjectDetailsByFlagId(flagId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !projectWithMembers.UserHasPermission(userId, domain.PermissionDeleteFlag) {
		return errors.New("not allowed to delete flag on this project")
	}

	err = i.flagRepository.DeleteFlag(flagId)
	if err != nil {
		return err
	}

	return nil
}

func (i *FlagInteractor) GetEnvironmentFlagsByPublicKey(key string) ([]string, error) {
	activeFlagsNames, err := i.flagRepository.GetEnvironmentActiveFlagsByPublicKey(key)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("cannot get active environment flags")
	}

	return activeFlagsNames, nil
}
