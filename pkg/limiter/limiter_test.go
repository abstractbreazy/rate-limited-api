package limiter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	correctTestIPs = []string{
		"123.45.67.89",
		"123.45.67.57",
		"123.45.67.1",
	}

	invalidTestIPs = []string{
		"123.45",
		"123",
		"correct_ip_i_promise",
	}
)

func TestLimiter(t *testing.T) {
	config := new(Config)
	config.SubnetMask = "24"

	for i := range correctTestIPs {
		subnet, err := config.Subnet(correctTestIPs[i])
		require.NoError(t, err)
		require.Equal(t, "123.45.67.0/24", subnet)
	}

	for i := range invalidTestIPs {
		_, err := config.Subnet(invalidTestIPs[i])
		require.Error(t, err)
	}
}
