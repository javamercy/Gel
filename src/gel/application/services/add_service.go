package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/context"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/utilities"
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
	validationResult := validator.Validate(request)

	if !validationResult.IsValid() {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, validationResult.Error())
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
