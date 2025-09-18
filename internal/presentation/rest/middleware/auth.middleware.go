package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RoleEmployee = "employee"
	RoleInvestor = "investor"
)

// Note: very basic authentication middleware
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range allowedRoles {
			var userIDHeader string

			switch role {
			case RoleEmployee:
				userIDHeader = c.GetHeader("x-employee-id")
				if userIDHeader != "" {
					userID, err := uuid.Parse(userIDHeader)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid x-employee-id header"})
						return
					}
					c.Set("employeeID", userID)
					c.Next()
					return
				}
			case RoleInvestor:
				userIDHeader = c.GetHeader("x-investor-id")
				if userIDHeader != "" {
					userID, err := uuid.Parse(userIDHeader)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid x-investor-id header"})
						return
					}
					c.Set("investorID", userID)
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "missing required role header"})
	}
}
