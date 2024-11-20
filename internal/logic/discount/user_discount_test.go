package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDiscountQuantityRollback(t *testing.T) {
	ctx := context.Background()
	t.Run("Rollback", func(t *testing.T) {
		err := g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath("/manifest/config")
		require.Nil(t, err)
		fmt.Println(userDiscountRollback(ctx, 15804))
	})
}
