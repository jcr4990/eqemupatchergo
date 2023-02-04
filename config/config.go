package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jbsmith7741/toml"
)

// Config represents a configuration parse
type Config struct {
	AutoPlay           string `toml:"AutoPlay" desc:"Blender Path to start Blender from"`
	AutoPatch          string `toml:"AutoPatch" desc:"EverQuest Path to copy converted zones to"`
	ClientVersion      string `toml:"eq_copy" desc:"copy eqgzi output to eq path"`
	LastPatchedVersion string `toml:"last_zone" desc:"Last zone selected"`
}

// New creates a new configuration
func New(ctx context.Context) (*Config, error) {
	var f *os.File
	cfg := Config{}
	path := "eqemupatcher.yml"

	isNewConfig := false
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config info: %w", err)
		}
		f, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create eqemupatcher.yml: %w", err)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("new config info: %w", err)
		}
		isNewConfig = true
	}
	if !isNewConfig {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open config: %w", err)
		}
	}

	defer f.Close()
	if fi.IsDir() {
		return nil, fmt.Errorf("eqemupatcher.yml is a directory, should be a file")
	}

	if isNewConfig {
		enc := toml.NewEncoder(f)
		cfg = getDefaultConfig()
		err = enc.Encode(cfg)
		if err != nil {
			return nil, fmt.Errorf("encode default: %w", err)
		}
		return &cfg, nil
	}

	_, err = toml.DecodeReader(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decode eqemupatcher.yml: %w", err)
	}

	return &cfg, nil
}

// Verify returns an error if configuration appears off
func (c *Config) Verify() error {

	return nil
}

func getDefaultConfig() Config {
	cfg := Config{}

	return cfg
}

// Save writes the config to disk
func (c *Config) Save() error {
	w, err := os.Create("eqemupatcher.yml")
	if err != nil {
		return fmt.Errorf("create eqemupatcher.yml: %w", err)
	}
	defer w.Close()

	enc := toml.NewEncoder(w)
	err = enc.Encode(c)
	if err != nil {
		return fmt.Errorf("encode default: %w", err)
	}
	return nil
}
