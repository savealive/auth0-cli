package manager

import (
	"github.com/spf13/viper"
	"gopkg.in/auth0.v3/management"
)

func New() (*management.Management, error) {
  m, err := management.New(
  	viper.GetString("domain"),
  	viper.GetString("clientID"),
  	viper.GetString("secret"))
  return m, err
}
