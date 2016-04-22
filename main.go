package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"github.com/teabough/docker-startup-script/config"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

const (
	pathContainingTempToken = "/tmp/"
)

func init() {

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

type EnvconsulConfig struct {
	Consul   string       `json:"consul,omitempty"`
	Token    string       `json:"token,omitempty"`
	Sanitize bool         `json:"sanitize,omitempty"`
	Vault    *VaultConfig `json:"vault,omitempty"`
}

type VaultConfig struct {
	Address string `json:"address,omitempty"`
	Token   string `json:"token"`
}

func main() {

	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
	envconsulPath := viper.GetString("Envconsul.Path")
	consulURL := viper.GetString("Consul.Url")
	vaultURL := viper.GetString("Vault.Url")
	consulToken := viper.GetString("Consul.Token")

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = vaultURL

	vault, err := api.NewClient(vaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	//GET TEMP TOKEN
	content, err := ioutil.ReadFile(pathContainingTempToken + "temp_" + strings.Replace(os.Getenv("MARATHON_APP_DOCKER_IMAGE"), "/", "_", -1))
	if err != nil {
		log.Fatal(err)
	}

	//GET PERM TOKEN
	vault.SetToken(string(content))
	log.WithFields(log.Fields{
		"token": string(content),
	}).Info("Set temp token for request")
	permSecret, err := vault.Logical().Read("cubbyhole/perm")

	log.WithFields(log.Fields{
		"perm token": permSecret,
	}).Info("Perm token from result")

	if err != nil {
		log.Fatal(err)
	}

	vaultConfigStruct := &VaultConfig{
		Address: vaultURL,
		Token:   permSecret.Data["token"].(string),
	}

	envconsulConfig := &EnvconsulConfig{
		Consul:   consulURL,
		Token:    consulToken,
		Sanitize: true,
		Vault:    vaultConfigStruct,
	}

	res, err := json.Marshal(envconsulConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("/envconsul_config.json", res, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("envconsul config written")

	args := []string{"envconsul", "-config", "/envconsul_config.json", "env"}

	env := os.Environ()

	execErr := syscall.Exec(envconsulPath, args, env)

	if execErr != nil {
		log.Warn("Something went wrong with envconsul command")
		log.Fatal(execErr)
	}

	log.Info("DONNNEEEE")
}
