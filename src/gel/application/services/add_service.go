package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/context"
	"Gel/src/gel/core/utilities"
	"errors"
)

type IAddService interface {
	Add(request *dto.AddRequest) *dto.AddResponse
}

type AddService struct {
	updateIndexService IUpdateIndexService
}

func NewAddService(updateIndexService IUpdateIndexService) *AddService {
	return &AddService{
		updateIndexService: updateIndexService,
	}
}

func (addService *AddService) Add(request *dto.AddRequest) *dto.AddResponse {
	validator := validators.NewAddValidator()
	validationResult := validator.Validate(request)

	if !validationResult.IsValid() {
		return dto.NewAddResponse(nil, errors.New(validationResult.Error()))
	}

	ctx := context.GetContext()
	pathResolver := utilities.NewPathResolver(ctx.RepositoryDir)
	normalizedPaths, err := pathResolver.Resolve(request.Pathspecs)
	if err != nil {
		return dto.NewAddResponse(nil, err)
	}

	if request.DryRun {
		return dto.NewAddResponse(normalizedPaths, nil)
	}

	err = addService.addPath(normalizedPaths)

	if err != nil {
		return dto.NewAddResponse(nil, err)
	}

	return dto.NewAddResponse(normalizedPaths, nil)
}

func (addService *AddService) addPath(paths []string) error {

	updateIndexRequest := dto.NewUpdateIndexRequest(paths, true, false)
	err := addService.updateIndexService.UpdateIndex(updateIndexRequest)
	return err
}
