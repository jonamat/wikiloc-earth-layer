package _setup

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	vp "github.com/spf13/viper"
)

func init() {
	// Build paths based on GOENV
	var basepath string

	var err error
	if basepath, err = os.Getwd(); err != nil {
		panic(err)
	}

	// Load envs
	if err := godotenv.Load(path.Join(basepath, "./.env")); err != nil {
		panic(err)
	}
	protocol := os.Getenv("PROTOCOL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	servicePort := os.Getenv("SERVICE_PORT")

	// Use 80 as fallback port
	if len(port) == 0 {
		port = "80"
	}
	if len(servicePort) == 0 {
		servicePort = port
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
	vp.SetConfigType("yaml")
	vp.ReadConfig(bytes.NewBuffer(config))

	// Setup global configuration
	vp.Set("basepath", basepath)
	vp.Set("protocol", protocol)
	vp.Set("host", host)
	vp.Set("port", port)
	vp.Set("servicePort", servicePort)
	vp.Set("serverURL", serverURL)
}

func Init() {}
