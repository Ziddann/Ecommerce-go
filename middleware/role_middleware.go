package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {

        role, exists := c.Get("role")

        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "User role not found",
            })
            c.Abort()
            return
        }

        userRole, ok := role.(string)

        if !ok || userRole != requiredRole {
            c.JSON(http.StatusForbidden, gin.H{
                "success": false,
                "message": "You do not have permission to access this resource",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}