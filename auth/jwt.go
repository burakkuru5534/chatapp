// auth/jwt.go
package auth

import (
	"example.com/m/models"
	"fmt"
	"time"


	"github.com/dgrijalva/jwt-go"
)

const hmacSecret = "SecretValueReplaceThis"
const defaulExpireTime = 604800 // 1 week

type Claims struct {
	ID string `json:"id"`
	Name string `json:"name"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

func (c *Claims) GetId() string {
	return c.ID
}

func (c *Claims) GetName() string {
	return c.Name
}

func (c *Claims) GetUserName() string {
	return c.UserName
}


// CreateJWTToken generates a JWT signed token for for the given user
func CreateJWTToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id": user.GetId(),
		"Name": user.GetName(),
		"ExpiresAt": time.Now().Unix() + defaulExpireTime,
	})
	tokenString, err := token.SignedString([]byte(hmacSecret))

	return tokenString, err
}

func ValidateToken(tokenString string) (models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}