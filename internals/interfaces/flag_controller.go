package interfaces

import (
	"github.com/Palma99/feature-flag-service/internals/application/usecase"
)

type FlagController struct {
	flagInteractor *usecase.FlagInteractor
}

func NewFlagController(flagInteractor *usecase.FlagInteractor) *FlagController {
	return &FlagController{
		flagInteractor: flagInteractor,
	}
}

// func (flagController *FlagController) GetAllFlagsByEnvironmentKey(key string) ([]entity.Flag, error) {
// 	return flagController.flagInteractor.GetAllFlagsByEnvironmentKey(key)
// }

// func (flagController *FlagController) UpdateFlagValue(key string, flagId int, value bool) error {
// 	return flagController.flagInteractor.UpdateFlagValue(key, flagId, value)
// }
