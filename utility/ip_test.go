package utility

import (
	"fmt"
	"testing"
)

func TestIpExtract(t *testing.T) {
	fmt.Println(ExtractFirstIPAddresses("57.180.13.66_Chrome_136.0.0.0(Macintosh)_9ea86cc8fdda2895b9e38b9ab50c9c66"))
}
