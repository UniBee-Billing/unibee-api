package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"go-oversea-pay/api/merchant/auth"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/utility"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"

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
	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantAccountByEmail(ctx, req.Email)
	utility.Assert(newOne == nil, "Email already existed")
	//if newOne != nil {
	//	return nil, gerror.NewCode(gcode.New(400, "Email already existed", nil))
	//}

	redisKey := fmt.Sprintf("MerchantAuth-Regist-Email:%s", req.Email)
	isDuplicatedInvoke := false
	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		return nil, gerror.Newf(`click too fast`)
	}

	userStr, err := json.Marshal(
		struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName string
			MerchantId                                                     uint64
		}{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Password:   hashAndSalt([]byte(req.Password)),
			Phone:      req.Phone,
			MerchantId: req.MerchantId,
			// Address:   req.Address,
			UserName: req.UserName,
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
	fmt.Println("verification ", verificationCode)
	// add merchant-verify, user-verify
	_, err = g.Redis().Set(ctx, req.Email+"-verify", verificationCode)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	//email.SendEmailToUser(req.Email, "Verification Code from UniBee", verificationCode)
	err = email.SendTemplateEmail(ctx, 0, req.Email, email.TemplateUserRegistrationCodeVerify, "", &email.TemplateVariable{
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}

	return &auth.RegisterRes{}, nil

	/*
		result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitNil().Insert(user)
		if err != nil {
			err = gerror.Newf(`record insert failure %s`, err)
			return
		}
		id, _ := result.LastInsertId()
		user.Id = uint64(id)
		var newOne *entity.UserAccount
		newOne = query.GetUserAccountById(ctx, user.Id)
		if newOne == nil {
			return nil, gerror.New("internal err:user query")
		}

		email.SendEmailToUser(newOne)
	*/

	// return &auth.RegisterRes{User: newOne}, nil
}
