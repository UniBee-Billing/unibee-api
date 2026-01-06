package message

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignature(t *testing.T) {
	body := "ddddd"
	secret := "ub_test_fH8IlSDGqiv30ruR4p38j4tsu9t9O31x"
	sign, arg := SignHMACWebhook(body, secret)
	t.Run("test hmac signature", func(t *testing.T) {
		require.Equal(t, "hmac", arg)
		require.Equal(t, nil, VerifyHMACSignature(body, secret, sign))
	})
}
