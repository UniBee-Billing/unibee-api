package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"go-oversea-pay/api/merchant/auth"
	"log"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"go-oversea-pay/internal/logic/email"
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

// ???????????????????? 8entitty.MerchantAccount
func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	var newOne *entity.MerchantUserAccount //  .UserAccount
	newOne = query.GetMerchantAccountByEmail(ctx, req.Email) //Id(ctx, user.Id)
	if newOne != nil {
		return nil, gerror.NewCode(gcode.New(400, "Email already existed", nil))
	}

	userStr, err := json.Marshal(
		struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName string
		}{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  hashAndSalt([]byte(req.Password)),
			// Phone:     req.Phone,
			Address:   req.Address,
			// UserName:  req.UserName,
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
	// add merchant-verify, user-verify
	_, err = g.Redis().Set(ctx, req.Email+"-verify", verificationCode)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	email.SendEmailToUser(req.Email, "Verification Code from Unibee", verificationCode)

	return &auth.RegisterRes{}, nil

	/*
		result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitEmpty().Insert(user)
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
