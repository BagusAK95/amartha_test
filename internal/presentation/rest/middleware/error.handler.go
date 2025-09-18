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
				res := gin.H{"message": e.Error()}
				if len(e.Errors) > 0 {
					res["errors"] = e.Errors
				}
				c.JSON(http.StatusNotFound, res)
			case *error.ForbiddenError:
				res := gin.H{"message": e.Error()}
				if len(e.Errors) > 0 {
					res["errors"] = e.Errors
				}
				c.JSON(http.StatusForbidden, res)
			case *error.BadRequestError:
				res := gin.H{"message": e.Error()}
				if len(e.Errors) > 0 {
					res["errors"] = e.Errors
				}
				c.JSON(http.StatusBadRequest, res)
			case *error.InternalServerError:
				res := gin.H{"message": e.Error()}
				if len(e.Errors) > 0 {
					res["errors"] = e.Errors
				}
				c.JSON(http.StatusInternalServerError, res)
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occurred"})
			}
		}
	}
}
