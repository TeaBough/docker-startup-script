package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
	"os"
	"encoding/json"
	"syscall"
)

func init() {

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

type EnvconsulConfig struct {
	Consul   string            `json:"consul,omitempty"`
	Token    string            `json:"token,omitempty"`
	Sanitize bool               `json:"sanitize,omitempty"`
	Vault    *VaultConfig       `json:"vault,omitempty"`
}

type VaultConfig struct {
	Address string            `json:"address,omitempty"`
	Token   string            `json:"token"`
}

func main() {

	config := api.DefaultConfig()
	config.Address = "http://localhost:8200"

	vault, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile("/tmp/temp_alpine")
	if err != nil {
		log.Fatal(err)
	}

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

	vaultConfig := &VaultConfig{
		Address: "http://127.0.0.1:8200",
		Token: permSecret.Data["token"].(string),
	}

	envconsulConfig := &EnvconsulConfig{
		Consul: "http://127.0.0.1:8500",
		Token: "toto",
		Sanitize: true,
		Vault: vaultConfig,
	}

	res, err := json.Marshal(envconsulConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("/tmp/envconsul_config.json", res, 0644)
	if err != nil {
		log.Fatal(err)
	}

	args := []string{"envconsul","-config", "/tmp/envconsul_config.json", "env"}

	env := os.Environ()

	execErr := syscall.Exec("/home/tibo/go/src/github.com/teabough/docker-startup-script/", args, env)

	if execErr != nil {
		log.Fatal(execErr)
	}

	log.Info("DONNNEEEE")
}
