package configparser

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// type Config struct {
// 	NUM8 struct {
// 		Proxy           string   `yaml:"Proxy"`
// 		TrustedDomains  []string `yaml:"trustedDomain"`
// 		TemplateSources []string `yaml:"templateSources"`
// 	} `yaml:"NUM8"`

// 	DatabaseConfig struct {
// 		Location string `yaml:"location"`
// 		Port     int    `yaml:"port"`
// 		Database string `yaml:"database"`
// 		Username string `yaml:"username"`
// 		Password string `yaml:"password"`
// 	} `yaml:"Database"`
// }

// func ParseConfig(path string) (*Config, error) {
// 	// Read file data
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		// log.Fatalf("error: %v", err)
// 		return nil, err
// 	}

// 	// Initialize configuration
// 	var config *Config

// 	// Unmarshal YAML data into Config struct
// 	err = yaml.Unmarshal(data, &config)
// 	if err != nil {
// 		// log.Fatalf("error: %v", err)
// 		return nil, err
// 	}
// 	return config, nil
// }

func InitConfigParser() (*viper.Viper, error) {
	var err error
	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigType("yaml")
	v.SetConfigName("configuration")
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file has changed:", e.Name)
	})
	v.WatchConfig()
	// If a config file is found, read it in.
	err = v.ReadInConfig()
	return v, err
}
