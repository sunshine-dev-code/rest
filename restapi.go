package rest

import (
	"github.com/sunshine-dev-code/rest/errs"
	"github.com/sunshine-dev-code/rest/response"

	"github.com/gin-gonic/gin"
)

// rest wrapper function
type restApi[T any] func(*gin.Context) (T, error)

func RestApiWrapper[T any](api restApi[T]) gin.HandlerFunc {
	wrappedApi := func(c *gin.Context) {
		ret, err := api(c)

		if err == nil {
			resp := response.NewOkRestResponse(ret)
			c.JSON(200, resp)
			return
		}

		// process error
		if response := errs.GetResponseFromError(err); response != nil {
			c.JSON(200, response)
			return
		}

		// unknown error
		c.JSON(200, response.ErrUnknowResponse)
	}

	return wrappedApi
}
