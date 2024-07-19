package utility

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type GatewayBank struct {
	AccountHolder string `json:"accountHolder"   dc:"The AccountHolder of wire transfer " v:"required" `
	BIC           string `json:"bic"   dc:"The BIC of wire transfer " v:"required" `
	IBAN          string `json:"iban"   dc:"The IBAN of wire transfer " v:"required" `
	Address       string `json:"address"   dc:"The address of wire transfer " v:"required" `
}

func TestJsonMarshal(t *testing.T) {
	fmt.Println(ToFirstCharLowerCase("firstName"))
	t.Run("test for json marshal", func(t *testing.T) {
		one := &GatewayBank{
			AccountHolder: "1",
			BIC:           "2",
			IBAN:          "3",
			Address:       "4",
		}
		json := MarshalToJsonString(one)
		require.NotNil(t, json)
		var next *GatewayBank
		err := UnmarshalFromJsonString(json, &next)
		require.Nil(t, err)
		require.NotNil(t, next)
	})

}
