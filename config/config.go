package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Mode string `toml:"mode" json:"mode"`
	App  struct {
		Domain string `toml:"domain" json:"domain"`
		Port   int    `toml:"port" json:"port"`
	} `toml:"app" json:"app"`
	Matrix struct {
		Homeserver string `toml:"homeserver" json:"homeserver"`
		ServerName string `toml:"server_name" json:"server_name"`
		DB         string `toml:"db" json:"db"`
	} `json:"matrix"`
	Cache struct {
		PublicRooms bool `json:"public_rooms"`
	} `json:"cache"`
	Security struct {
		AllowedOrigins []string `toml:"allowed_origins" json:"allowed_origins"`
	} `json:"security"`
	Capabilities struct {
		PublicRooms struct {
			ListRooms     bool `json:"list_rooms" toml:"list_rooms"`
			ViewHierarchy bool `json:"view_hierarchy" toml:"view_hierarchy"`
			ReadMessages  bool `json:"read_messages" toml:"read_messages"`
		} `toml:"public_rooms" json:"public_rooms"`
	} `json:"capabilities"`
}

var conf Config

// Read reads the config file and returns the Values struct
func Read(s string) (*Config, error) {
	file, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(b), &conf); err != nil {
		panic(err)
	}

	return &conf, err
}
