package user

import (
	"context"
	"fmt"
	"go-oversea-pay/api/user/auth"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"log"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

var secretKey = []byte("3^&secret-key-for-UniBee*1!8*")

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			// "userId": userId,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	if req.Email == "" {
		// return nil, gerror.New("email empty")
		return nil, gerror.NewCode(gcode.New(400, "email cannot be empty", nil))
	}

	if req.Password == "" {
		return nil, gerror.NewCode(gcode.New(400, "password cannot be empty", nil))
	}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, req.Email)
	if newOne == nil {
		// return nil, gerror.New("internal err: user not found")
		return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	}
	if !comparePasswords(newOne.Password, []byte(req.Password)) { // wrong password
		return nil, gerror.NewCode(gcode.New(400, "Login failed", nil))
	}

	token, err := createToken(req.Email)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	return &auth.LoginRes{User: newOne, Token: token}, nil
}
