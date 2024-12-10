package errs

import (
	"github.com/sunshine-dev-code/rest/response"
)

type ErrHandler = func(error) *response.RestResponse

var errHandlerList = []ErrHandler{getErrResponse, getPgErrResponse}

func GetResponseFromError(err error) *response.RestResponse {
	for _, errHandler := range errHandlerList {
		if response := errHandler(err); response != nil {
			return response
		}
	}
	return nil
}
