package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"
	"unibee-api/api/user/auth"
	"unibee-api/internal/logic/email"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func IsEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	validEmail := regexp.MustCompile(emailRegex)
	return validEmail.MatchString(email)
}

func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(IsEmailValid(req.Email), "Invalid Email")

	redisKey := fmt.Sprintf("UserAuth-Regist-Email:%s", req.Email)
	//isDuplicatedInvoke := false
	//defer func() {
	//	if !isDuplicatedInvoke {
	//		utility.ReleaseLock(ctx, redisKey)
	//	}
	//}()

	if !utility.TryLock(ctx, redisKey, 10) {
		//isDuplicatedInvoke = true
		utility.Assert(false, "click too fast, please wait for second")
	}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, req.Email) //Id(ctx, user.Id)
	utility.Assert(newOne == nil, "Email already existed")
	//if newOne != nil {
	//	return nil, gerror.NewCode(gcode.New(400, "Email already existed", nil))
	//}

	userStr, err := json.Marshal(
		struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName, CountryCode, CountryName string
		}{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			CountryCode: req.CountryCode,
			CountryName: req.CountryName,
			Password:    hashAndSalt([]byte(req.Password)),
			Phone:       req.Phone,
			Address:     req.Address,
			UserName:    req.UserName,
		},
	)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	_, err = g.Redis().Set(ctx, req.Email, userStr)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	_, err = g.Redis().Expire(ctx, req.Email, 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	verificationCode := generateRandomString(6)
	fmt.Printf("verification ", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-verify", verificationCode)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	//email.SendEmailToUser(req.Email, "Verification Code from UniBee", verificationCode)
	err = email.SendTemplateEmail(ctx, 0, req.Email, "", email.TemplateUserRegistrationCodeVerify, "", &email.TemplateVariable{
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}

	return &auth.RegisterRes{}, nil
}
