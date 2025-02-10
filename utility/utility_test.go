package utility

import (
	"fmt"
	"testing"
)

func TestMergeMetadata(t *testing.T) {
	t.Run("TestMergeMetadata", func(t *testing.T) {
		var oldOne string = "{\"data\":1}"
		var newOne map[string]interface{} = map[string]interface{}{"test": 1}
		fmt.Println(MergeMetadata(oldOne, &newOne))
	})
}

func TestConvertCentToDollarFloat(t *testing.T) {
	ConvertDollarStrToCent("2000.00", "RUB")
}
