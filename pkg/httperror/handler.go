package httperror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleError(err error, c echo.Context) {
	var he *echo.HTTPError

	if errors.As(err, &he) {
		if he.Internal != nil {
			c.Logger().Errorf("HTTP %d - %v (due to: %v)", he.Code, he.Message, he.Internal)
		} else {
			c.Logger().Errorf("HTTP %d - %v", he.Code, echo.Map{
				"message": he.Message,
			})
		}

		c.JSON(he.Code, he.Message)
		return
	}

	c.Logger().Errorf("HTTP 500 - InternalServerError (due to: %v)", err)
	c.JSON(http.StatusInternalServerError, echo.Map{
		"message": "Internal server error",
	})
	return
}
