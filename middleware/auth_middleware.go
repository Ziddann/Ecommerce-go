package middleware

import (
    "net/http"
    "strings"

    "ecommerce/services"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {

        authHeader := c.GetHeader("Authorization")

        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Authorization token is required",
            })
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)

        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Invalid authorization format",
            })
            c.Abort()
            return
        }

        tokenString := parts[1]

        claims := &services.JWTClaims{}

        token, err := jwt.ParseWithClaims(
            tokenString,
            claims,
            func(token *jwt.Token) (interface{}, error) {
                if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
                    return nil, jwt.ErrSignatureInvalid
                }

                return []byte(jwtSecret), nil
            },
            jwt.WithIssuer("ecommerce-api"),
        )

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)

        c.Next()
    }
}