package config

import (
	"fmt"
	"github.com/TwiN/go-color"

	"os"

	"gopkg.in/yaml.v2"
)

/*
Environment variables:
- ZTP_CONFIGFILE
- KUBECONFIG
*/

var (
	// ConfigFile is the global configuration
	Ztp ZTPConfig //global variable to reference the config
)

const (
	DEFAULT_CONFIG_FILE = "./config.yaml"
)

//ZTPConfig is the global configuration data model
type ZTPConfig struct {
	Config struct {
		ConfigFile       string
		Clusterimageset  string `yaml:"clusterimageset"`
		KubeconfigHUB    string
		KubeframeNS      string
		OC_OCP_VERSION   string `yaml:"OC_OCP_VERSION"`
		OC_OCP_TAG       string `yaml:"OC_OCP_TAG"`
		OC_RHCOS_RELEASE string `yaml:"OC_RHCOS_RELEASE"`
		OC_ACM_VERSION   string `yaml:"OC_ACM_VERSION"`
		OC_OCS_VERSION   string `yaml:"OC_OCS_VERSION"`
	} `yaml:"config"`
	Spokes []struct {
		Name            string `yaml:"name"`
		KubeconfigSPOKE string
		Master0         struct {
			NicExtDhcp   string   `yaml:"nic_ext_dhcp"`
			NicIntStatic string   `yaml:"nic_int_static"`
			MacIntStatic string   `yaml:"mac_int_static"`
			MacExtDhcp   string   `yaml:"mac_ext_dhcp"`
			BmcUrl       string   `yaml:"bmc_url"`
			BmcUser      string   `yaml:"bmc_user"`
			BmcPass      string   `yaml:"bmc_pass"`
			StorageDisk  []string `yaml:"storage_disk"`
		} `yaml:"master0"`
		Master1 struct {
			NicExtDhcp   string   `yaml:"nic_ext_dhcp"`
			NicIntStatic string   `yaml:"nic_int_static"`
			MacIntStatic string   `yaml:"mac_int_static"`
			MacExtDhcp   string   `yaml:"mac_ext_dhcp"`
			BmcUrl       string   `yaml:"bmc_url"`
			BmcUser      string   `yaml:"bmc_user"`
			BmcPass      string   `yaml:"bmc_pass"`
			StorageDisk  []string `yaml:"storage_disk"`
		} `yaml:"master1"`
		Master2 struct {
			NicExtDhcp   string   `yaml:"nic_ext_dhcp"`
			NicIntStatic string   `yaml:"nic_int_static"`
			MacIntStatic string   `yaml:"mac_int_static"`
			MacExtDhcp   string   `yaml:"mac_ext_dhcp"`
			BmcUrl       string   `yaml:"bmc_url"`
			BmcUser      string   `yaml:"bmc_user"`
			BmcPass      string   `yaml:"bmc_pass"`
			StorageDisk  []string `yaml:"storage_disk"`
		} `yaml:"master2"`
		Worker0 struct {
			NicExtDhcp   string   `yaml:"nic_ext_dhcp"`
			NicIntStatic string   `yaml:"nic_int_static"`
			MacIntStatic string   `yaml:"mac_int_static"`
			MacExtDhcp   string   `yaml:"mac_ext_dhcp"`
			BmcUrl       string   `yaml:"bmc_url"`
			BmcUser      string   `yaml:"bmc_user"`
			BmcPass      string   `yaml:"bmc_pass"`
			StorageDisk  []string `yaml:"storage_disk"`
		} `yaml:"worker0"`
	} `yaml:"spokes"`
}

//fmt.Println(e.Spokes[0].Name, e.Spokes[0].Master0.NicExtDhcp)

//Constructor new config file from file
func NewConfig() (ZTPConfig, error) {
	err := Ztp.ReadFromConfigFile()
	if err != nil {
		return Ztp, err
	}
	return Ztp, nil
}

//ReadFromConfigFile reads the config file
func (c *ZTPConfig) ReadFromConfigFile() error {
	if getEnv("ZTP_CONFIGFILE") == "" {
		fmt.Errorf(color.InRed("ZTP_CONFIGFILE not set"), "ZTP_CONFIGFILE not set")
	}

	if getEnv("ZTP_CONFIGFILE") != "" {
		fmt.Println(color.InYellow(">>>> ConfigFile env is not empty. Reading file from this env"))
		c.Config.ConfigFile = getEnv("ZTP_CONFIGFILE")
	} else {
		fmt.Println(color.InYellow(">>>> ZTP_CONFIGFILE env var is empty. Using default path: " + DEFAULT_CONFIG_FILE))
		c.Config.ConfigFile = DEFAULT_CONFIG_FILE
	}

	f, err := os.Open(c.Config.ConfigFile)
	if err != nil {
		return fmt.Errorf(color.InRed("opening config file %s: %v"), c.Config.ConfigFile, err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf(color.InRed("decoding config file %s: %v"), c.Config.ConfigFile, err)
	}
	if getEnv("KUBECONFIG") == "" {
		return fmt.Errorf(color.InRed("Kubeconfig env empty"), "")
	}
	c.Config.KubeconfigHUB = getEnv("KUBECONFIG")

	return nil
}

//getEnv returns the value of the environment variable named by the key.
func getEnv(key string) string {
	return os.Getenv(key)
}
