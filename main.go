package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"github.com/teabough/docker-startup-script/config"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	Upcase   bool         `json:"upcase,omitempty"`
	Vault    *VaultConfig `json:"vault,omitempty"`
}

type VaultConfig struct {
	Address string `json:"address,omitempty"`
	Renew   bool   `json:"renew"`
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
	content, err := ioutil.ReadFile(pathContainingTempToken + "temp_" + os.Getenv("HOSTNAME"))
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
		Renew:   false,
		Token:   permSecret.Data["token"].(string),
	}

	envconsulConfig := &EnvconsulConfig{
		Consul:   consulURL,
		Token:    consulToken,
		Sanitize: true,
		Upcase:   true,
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

	consulConfig := consulApi.DefaultConfig()
	consulConfig.Address = consulURL
	consulConfig.Token = consulToken

	client, err := consulApi.NewClient(consulConfig)

	pair, _, err := client.KV().Get("secrets-permissions/"+os.Getenv("APP_NAME"), nil)
	if err != nil {
		log.Fatal(err)
	}

	secrets := strings.Split(string(pair.Value), ",")
	secretsCmd := ""
	for _, secret := range secrets {
		secretsCmd = fmt.Sprintf("%s-secret secret/%s", secretsCmd, secret)
	}

	cmd := "/envconsul -config /envconsul_config.json -once " + secretsCmd + " env"

	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()
	fmt.Printf("%s", out)
	if err != nil {
		log.WithFields(log.Fields{
			"envconsulPath": envconsulPath,
			"args":          cmd,
		}).Warn("Something went wrong with envconsul command")
		log.Fatal(err)
	}

	log.Info("DONNNEEEE")
}
