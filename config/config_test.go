package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/teabough/docker-startup-script/config"
	"os"
)

var _ = Describe("Config", func() {
	var (
		err error
	)

	AfterEach(func() {
		os.Unsetenv("CONSUL_URL")
		os.Unsetenv("ENVCONSUL_PATH")
		os.Unsetenv("VAULT_URL")
		os.Unsetenv("CONSUL_TOKEN")

	})
	JustBeforeEach(func() {
		err = config.ReadConfig()
	})
	Context("When setting all the env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CONSUL_URL", "1")
			os.Setenv("ENVCONSUL_PATH", "2")
			os.Setenv("VAULT_URL", "8")
			os.Setenv("CONSUL_TOKEN", "9")
		})
		It("should return nil", func() {

			Expect(err).NotTo(HaveOccurred())
			Expect(viper.GetString("Consul.Url")).To(Equal("1"))
			Expect(viper.GetString("EnvConsul.Path")).To(Equal("2"))
			Expect(viper.GetString("Vault.Url")).To(Equal("8"))
			Expect(viper.GetString("Consul.Token")).To(Equal("9"))
		})

	})

	Context("When the vault url is missing", func() {
		BeforeEach(func() {
			os.Setenv("CONSUL_URL", "1")
			os.Setenv("ENVCONSUL_PATH", "1")
			os.Setenv("CONSUL_TOKEN", "1")
		})
		It("should eror", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VAULT_URL is not defined"))
		})
	})

	Context("When the envconsul path is missing", func() {
		BeforeEach(func() {
			os.Setenv("CONSUL_URL", "1")
			os.Setenv("VAULT_URL", "1")
			os.Setenv("CONSUL_TOKEN", "1")
		})
		It("should eror", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("ENVCONSUL_PATH is not defined"))
		})
	})

	Context("When the vault url is missing", func() {
		BeforeEach(func() {
			os.Setenv("ENVCONSUL_PATH", "1")
			os.Setenv("CONSUL_URL", "1")
			os.Setenv("CONSUL_TOKEN", "1")
		})
		It("should eror", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VAULT_URL is not defined"))
		})
	})

	Context("When the consul token is missing", func() {
		BeforeEach(func() {
			os.Setenv("ENVCONSUL_PATH", "1")
			os.Setenv("CONSUL_URL", "1")
			os.Setenv("VAULT_URL", "1")
		})
		It("should eror", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("CONSUL_TOKEN is not defined"))
		})
	})
})
