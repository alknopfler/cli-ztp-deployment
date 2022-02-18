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
	MODE_HUB            = "hub"
	MODE_SPOKE          = "spoke"
)

//ZTPConfig is the global configuration data model
type ZTPConfig struct {
	Config struct {
		ConfigFile      string
		KubeconfigHUB   string
		Clusterimageset string `yaml:"clusterimageset"`
		OcOCPVersion    string `yaml:"OC_OCP_VERSION"`
		OcOCPTag        string `yaml:"OC_OCP_TAG"`
		OcRHCOSRelease  string `yaml:"OC_RHCOS_RELEASE"`
		OcACMVersion    string `yaml:"OC_ACM_VERSION"`
		OcOCSVersion    string `yaml:"OC_OCS_VERSION"`
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
	//Read main config from the config file
	err := Ztp.ReadFromConfigFile()
	if err != nil {
		return Ztp, err
	}
	// Set the rest of config from env
	if getEnv("KUBECONFIG") == "" {
		return Ztp, fmt.Errorf(color.InRed("Kubeconfig env empty"), "")
	}
	Ztp.Config.KubeconfigHUB = getEnv("KUBECONFIG")
	return Ztp, nil
}

//ReadFromConfigFile reads the config file
func (c *ZTPConfig) ReadFromConfigFile() error {
	if getEnv("ZTP_CONFIGFILE") == "" {
		fmt.Println(color.InRed("ZTP_CONFIGFILE not set"))
		return fmt.Errorf("ZTP_CONFIGFILE not set")
	}

	if getEnv("ZTP_CONFIGFILE") != "" {
		fmt.Println(color.InYellow(">>>> [INFO] ConfigFile env is not empty. Reading file from this env"))
		c.Config.ConfigFile = getEnv("ZTP_CONFIGFILE")
	} else {
		fmt.Println(color.InYellow(">>>> [INFO] ZTP_CONFIGFILE env var is empty. Using default path: " + DEFAULT_CONFIG_FILE))
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
	return nil
}

//getEnv returns the value of the environment variable named by the key.
func getEnv(key string) string {
	return os.Getenv(key)
}

func GetKubeconfigFromMode(mode string) string {
	if isHub(mode) {
		return Ztp.Config.KubeconfigHUB
	}
	return Ztp.Spokes[0].KubeconfigSPOKE
}

func isHub(mode string) bool {
	return mode == MODE_HUB
}

func isSpoke(mode string) bool {
	return mode == MODE_SPOKE
}
