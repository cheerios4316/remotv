package input

import (
	"encoding/json"
	"flag"
	"os"
)

type Flags struct {
	DeviceName     string
	URI            string
	Port           int
	BrowserBinName string
	ConfigPath     string
}

func ParseFlags() Flags {
	var deviceName, uri, configPath string
	var port int

	flag.StringVar(&deviceName, "device-name", "my-device", "Device display name")
	flag.StringVar(&deviceName, "D", "my-device", "Device display name (shorthand)")

	flag.IntVar(&port, "port", 9000, "TCP Port")
	flag.IntVar(&port, "P", 9000, "TCP Port")

	flag.StringVar(&uri, "uri", "localhost", "Server URI")
	flag.StringVar(&uri, "U", "localhost", "Server URI (shorthand)")

	flag.StringVar(&configPath, "config", "config.json", "Config JSON file path")
	flag.StringVar(&configPath, "C", "config.json", "Config JSON file path (shorthand)")

	flag.Parse()

	return Flags{
		DeviceName: deviceName,
		URI:        uri,
		Port:       port,
		ConfigPath: configPath,
	}
}

type Config struct {
	Browser  string
	Shutdown string
}

func DefaultConfig() Config {
	return Config{
		Browser:  "firefox",
		Shutdown: "shutdown now",
	}
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
