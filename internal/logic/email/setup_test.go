package email

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDefaultEmailTemplate(t *testing.T) {
	t.Run("Test Get Default Email Template From Cloud Api", func(t *testing.T) {
		list := FetchDefaultEmailTemplateFromCloudApi()
		require.NotNil(t, list)
		require.Equal(t, true, len(list) > 0)
	})
}
