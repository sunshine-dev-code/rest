package errs

import (
	"errors"

	"github.com/sunshine-dev-code/rest/response"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// errors
var (
	ErrUnknow           = errors.New("unknow_error")
	ErrParameter        = errors.New("parameter_error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrExistentResponse = errors.New("existent_parameter")
	ErrStageTransform   = errors.New("stage_transform")
	ErrPermissionDenied = errors.New("permission_denied")
)

// errors information
var ErrorMap = map[error]*response.RestResponse{
	gorm.ErrRecordNotFound:              &response.ErrRecordNotFoundResponse,
	gorm.ErrDuplicatedKey:               &response.ErrExistentResponse,
	ErrUnauthorized:                     &response.ErrUnauthorizedResponse,
	ErrParameter:                        &response.ErrParamsResponse,
	ErrPermissionDenied:                 &response.ErrPermissionDenied,
	ErrExistentResponse:                 &response.ErrParamsResponse,
	bcrypt.ErrMismatchedHashAndPassword: &response.ErrParamsResponse,
}

func getErrResponse(err error) *response.RestResponse {
	if response, ok := ErrorMap[err]; ok {
		return response
	}
	return nil
}
