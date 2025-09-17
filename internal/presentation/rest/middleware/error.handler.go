package middleware

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case *error.NotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"message": e.Error()})
			case *error.ForbiddenError:
				c.JSON(http.StatusForbidden, gin.H{"message": e.Error()})
			case *error.BadRequestError:
				c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
			case *error.InternalError:
				c.JSON(http.StatusInternalServerError, gin.H{"message": e.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occurred"})
			}
		}
	}
}
