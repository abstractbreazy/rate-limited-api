package limiter

import (
	"net"
	"strings"
	"time"
)

// Limiting config
type Config struct {
	Limit      uint64        `yaml:"limit" mapstructure:"limit"`
	Timeout    time.Duration `yaml:"timeout" mapstructure:"timeout"`
	SubnetMask string        `yaml:"subnet_mask" mapstructure:"subnet_mask"`
}

func NewConfig() *Config {
	return &Config{}
}

// Returns the subnet string by given configurations and IP.
func (r *Config) Subnet(ip string) (string, error) {
	address := strings.Join([]string{ip, r.SubnetMask}, "/")
	_, subnet, err := net.ParseCIDR(address)
	if err != nil {
		return "", err
	}

	return subnet.String(), nil
}
