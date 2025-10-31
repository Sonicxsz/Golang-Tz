package main

import (
	_ "awesomeProject1/docs"
	"awesomeProject1/internal/server"
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.yaml", "Path to API server config")
}

// @title           Subscription API
// @version         1.0
// @description     API для управления подписками пользователей.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	flag.Parse()
	config := server.NewConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		println("Cannot read config file, using default values")
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			println("Cannot parse YAML config, using default values")
		}
	}

	println(config.Storage.DbConnString)
	if err := config.Storage.RunMigrations(); err != nil {
		log.Fatalf("Migration error: %v", err.Error())
	}

	api := server.New(config)

	println("Server starting")
	log.Fatal(api.Start())
}
