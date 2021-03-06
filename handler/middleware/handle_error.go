package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/gerrors"
)

func NewHandleErrorMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			var ge *gerrors.GlobalError

			switch errorToPrint.Err {
			case gorm.ErrRecordNotFound:
				ge = errors.NotFound(errorToPrint.Err.Error()).(*gerrors.GlobalError)
			default:
				ge = &gerrors.GlobalError{}
				if json.Unmarshal([]byte(errorToPrint.Err.Error()), ge) != nil {
					ge = errors.InternalServerError(errorToPrint.Err.Error(), errorToPrint.Err).(*gerrors.GlobalError)
				}
			}

			if ge.ServiceName == "" {
				ge.ServiceName = serviceName
			}
			c.JSON(ge.StatusCode, gin.H{
				"code":    ge.Code,
				"message": ge.Message,
			})
		}

	}
}
