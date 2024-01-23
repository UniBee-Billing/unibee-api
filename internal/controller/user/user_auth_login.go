package user

import (
	"context"
	"fmt"
	"go-oversea-pay/api/user/auth"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
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
		fmt.Printf("comparePasswords err:%s\n", err.Error())
		return false
	}
	return true
}

var secretKey = []byte("3^&secret-key-for-UniBee*1!8*")

func createToken(email string, userId uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"id":    userId,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
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

	utility.Assert(token.Valid, "Invalid Code")
	//if !token.Valid {
	//	return fmt.Errorf("invalid token")
	//}

	return nil
}

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "email cannot be empty")
	utility.Assert(req.Password != "", "password cannot be empty")
	//if req.Email == "" {
	//	// return nil, gerror.New("email empty")
	//	return nil, gerror.NewCode(gcode.New(400, "email cannot be empty", nil))
	//}
	//
	//if req.Password == "" {
	//	return nil, gerror.NewCode(gcode.New(400, "password cannot be empty", nil))
	//}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "Login Failed")
	//if newOne == nil {
	//	// return nil, gerror.New("internal err: user not found")
	//	return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	//}
	utility.Assert(comparePasswords(newOne.Password, []byte(req.Password)), "Login Failed, Password Not Match")
	//if !comparePasswords(newOne.Password, []byte(req.Password)) { // wrong password
	//	return nil, gerror.NewCode(gcode.New(400, "Login failed", nil))
	//}

	token, err := createToken(req.Email, newOne.Id)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", newOne.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	return &auth.LoginRes{User: newOne, Token: token}, nil
}
