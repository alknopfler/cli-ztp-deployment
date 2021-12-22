package config

import (
	"fmt"

	"os"

	"gopkg.in/yaml.v2"
)

/*
Environment variables:
- ZTP_CONFIGFILE
- KUBECONFIG


*/

const (
	defaultUserConfigFile = "./config.yaml"
)

var (
	// ConfigFile is the global configuration
	ConfigFile string
	Ztp        ZTPConfig //global variable to reference the config
)

//ZTPConfig is the global configuration data model
type ZTPConfig struct {
	Config struct {
		Clusterimageset  string `yaml:"clusterimageset"`
		OC_OCP_VERSION   string `yaml:"OC_OCP_VERSION"`
		OC_OCP_TAG       string `yaml:"OC_OCP_TAG"`
		OC_RHCOS_RELEASE string `yaml:"OC_RHCOS_RELEASE"`
		OC_ACM_VERSION   string `yaml:"OC_ACM_VERSION"`
		OC_OCS_VERSION   string `yaml:"OC_OCS_VERSION"`
		KubeconfigHUB    string
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
func NewConfig() error {
	return Ztp.ReadFromConfigFile()
}

//ReadFromConfigFile reads the config file
func (c *ZTPConfig) ReadFromConfigFile() error {
	if getEnv("ZTP_CONFIGFILE") == "" {
		fmt.Errorf("ZTP_CONFIGFILE not set", "")
	}

	if getEnv("ZTP_CONFIGFILE") != "" {
		fmt.Println(">>>> ConfigFile env is not empty. Reading file from this env")
		ConfigFile = getEnv("ZTP_CONFIGFILE")
	} else {
		fmt.Println(">>>> ZTP_CONFIGFILE env var is empty. Using default path: " + defaultUserConfigFile)
		ConfigFile = defaultUserConfigFile
	}

	f, err := os.Open(ConfigFile)
	if err != nil {
		return fmt.Errorf("opening config file %s: %v", ConfigFile, err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("decoding config file %s: %v", ConfigFile, err)
	}
	if getEnv("KUBECONFIG") == "" {
		return fmt.Errorf("Kubeconfig env empty", "")
	}
	c.Config.KubeconfigHUB = getEnv("KUBECONFIG")

	return nil
}

//getEnv returns the value of the environment variable named by the key.
func getEnv(key string) string {
	return os.Getenv(key)
}
