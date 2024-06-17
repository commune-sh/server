package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type App struct {
	Domain string `toml:"domain"`
	Port   int    `toml:"port"`
}

type Matrix struct {
	Homeserver string `toml:"homeserver"`
	ServerName string `toml:"server_name"`
	DB         string `toml:"db"`
}

type Security struct {
	AllowedOrigins []string `toml:"allowed_origins"`
}

type Redis struct {
	Address         string `toml:"address"`
	Password        string `toml:"password"`
	SessionsDB      int    `toml:"sessions_db"`
	PostsDB         int    `toml:"posts_db"`
	SystemDB        int    `toml:"system_db"`
	NotificationsDB int    `toml:"notifications_db"`
}

type Cache struct {
	IndexEvents  bool `toml:"index_events"`
	SpaceEvents  bool `toml:"space_events"`
	EventReplies bool `toml:"event_replies"`
}

type Capabilities struct {
	PublicRooms struct {
		ListRooms     bool `json:"list_rooms" toml:"list_rooms"`
		ViewHierarchy bool `json:"view_hierarchy" toml:"view_hierarchy"`
		ReadMessages  bool `json:"read_messages" toml:"read_messages"`
	} `toml:"public_rooms" json:"public_rooms"`
}

type Config struct {
	Name         string       `toml:"name"`
	Mode         string       `toml:"mode"`
	App          App          `toml:"app"`
	Matrix       Matrix       `toml:"matrix"`
	Redis        Redis        `toml:"redis"`
	Cache        Cache        `toml:"cache"`
	Security     Security     `toml:"security"`
	Capabilities Capabilities `toml:"capabilities"`
}

var conf Config

// Read reads the config file and returns the Values struct
func Read(s string) (*Config, error) {
	file, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(b), &conf); err != nil {
		panic(err)
	}

	return &conf, err
}
