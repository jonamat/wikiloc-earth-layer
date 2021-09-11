package _init

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	// Build paths based on GOENV
	var basepath string
	bin, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if os.Getenv("GOENV") == "production" {
		basepath = path.Dir(bin)
	} else {
		var err error
		if basepath, err = os.Getwd(); err != nil {
			panic(err)
		}
	}

	// Load envs
	if err := godotenv.Load(path.Join(basepath, "./.env")); err != nil {
		panic(err)
	}
	protocol := os.Getenv("PROTOCOL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Use 80 as fallback port
	if len(port) == 0 {
		port = "80"
	}

	// Check envs
	if len(protocol) == 0 || len(host) == 0 {
		panic("Protocol or host not defined")
	}

	// Define server URL
	var serverURL string
	if port != "80" || len(port) == 0 {
		serverURL = fmt.Sprintf("%s://%s:%s", protocol, host, port)
	} else {
		serverURL = fmt.Sprintf("%s://%s", protocol, host)
	}

	// Get config.yml
	config, err := os.ReadFile(path.Join(basepath, "config.yml"))
	if err != nil {
		panic(err)
	}

	// Get user configuration
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(config))

	// Setup global configuration
	viper.Set("basepath", basepath)
	viper.Set("protocol", protocol)
	viper.Set("host", host)
	viper.Set("port", port)
	viper.Set("serverURL", serverURL)
}

func Init() {}