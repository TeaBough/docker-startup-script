package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

func ReadConfig() error {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if !viper.IsSet("Consul.Url") {
		return errors.New("CONSUL_URL is not defined")
	}

	if !viper.IsSet("EnvConsul.Path") {
		return errors.New("ENVCONSUL_PATH is not defined")
	}

	if !viper.IsSet("Vault.Url") {
		return errors.New("VAULT_URL is not defined")
	}

	return nil
}
