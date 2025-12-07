package services

import (
	"Gel/application/dto"
	"Gel/application/validators"
	"Gel/core/context"
	"Gel/core/crossCuttingConcerns/gelErrors"
	"Gel/core/utilities"
)

type IAddService interface {
	Add(request *dto.AddRequest) ([]string, *gelErrors.GelError)
}

type AddService struct {
	updateIndexService IUpdateIndexService
}

func NewAddService(updateIndexService IUpdateIndexService) *AddService {
	return &AddService{
		updateIndexService: updateIndexService,
	}
}

func (addService *AddService) Add(request *dto.AddRequest) ([]string, *gelErrors.GelError) {
	validator := validators.NewAddValidator()
	gelError := validator.Validate(request)

	if gelError != nil {
		return nil, gelError
	}

	ctx := context.GetContext()
	pathResolver := utilities.NewPathResolver(ctx.RepositoryDir)
	normalizedPaths, err := pathResolver.Resolve(request.Pathspecs)
	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	if request.DryRun {
		return normalizedPaths, nil
	}

	addPathErr := addService.addPath(normalizedPaths)

	if addPathErr != nil {
		return nil, addPathErr
	}

	return normalizedPaths, nil
}

func (addService *AddService) addPath(paths []string) *gelErrors.GelError {

	updateIndexRequest := dto.NewUpdateIndexRequest(paths, true, false)
	err := addService.updateIndexService.UpdateIndex(updateIndexRequest)
	return err
}
