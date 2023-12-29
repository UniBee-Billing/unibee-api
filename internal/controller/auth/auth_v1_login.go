package auth

import (
	"context"

	v1 "go-oversea-pay/api/auth/v1"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"log"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	if (req.Email == "") {
		// return nil, gerror.New("email empty")
		return nil, gerror.NewCode(gcode.New(400, "email cannot be empty", nil))
	}
	
	if (req.Password == "") {
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

	return &v1.LoginRes{User: newOne}, nil
}
