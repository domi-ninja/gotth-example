package webapp

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	toml "github.com/pelletier/go-toml/v2"
)

// this is a nice way to parse nested configs from a toml file IMO
type AppConfig struct {
	Server struct {
		Port        int    `toml:"port"`
		BindAddress string `toml:"bind_address"`
	} `toml:"server"`

	Site SiteConfig `toml:"site"`

	Assets struct {
		Path string `toml:"path"`
	} `toml:"assets"`

	Secrets struct {
		JWT_SECRET string
	}
}

// this struct is out here (instead of nested in AppConfig) because we reference it separate from AppConfig
type SiteConfig struct {
	AppPath string `toml:"app_path"`

	Title        string `toml:"title"`
	Description  string `toml:"description"`
	DefaultImage string `toml:"default_image"`
	Keywords     string `toml:"keywords"`
}

func MustLoadConfig(path string) *AppConfig {
	config := &AppConfig{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = toml.Unmarshal(bytes, config)
	if err != nil {
		log.Fatal(err)
	}
	debugPrint(config)

	err = godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
	// TODO: make sure we have a jwt secret (instead of empty string or something)
	// config.Secrets.JWT_SECRET = os.Getenv("JWT_SECRET")
	// if len(config.Secrets.JWT_SECRET) < 32 {
	// 	log.Fatal("JWT_SECRET with len >= 32 is required, auth is not secure without it")
	// }
	if config.Site.Title == "" {
		log.Fatal("site.title is required")
	}
	if config.Site.Description == "" {
		log.Fatal("site.description is required")
	}
	if config.Site.DefaultImage == "" {
		log.Fatal("site.default_image is required")
	}
	if config.Site.Keywords == "" {
		log.Fatal("site.keywords is required")
	}

	return config
}

func debugPrint(config *AppConfig) {
	config2 := *config // dont overwrite the real config in memory please
	if config2.Secrets.JWT_SECRET != "" {
		config2.Secrets.JWT_SECRET = "********" // dont print the jwt secret
	}
	json, err := json.MarshalIndent(config2, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Loaded config:\n", string(json))
}
