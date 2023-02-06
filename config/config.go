package config

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents a configuration parse
type Config struct {
	AutoPlay      string `yaml:"AutoPlay" desc:"Automatically Play"`
	AutoPatch     string `yaml:"AutoPatch" desc:"Automatically patch"`
	ClientVersion string `yaml:"ClientVersion" desc:"Version of client patched"`
	PatcherHash   string `yaml:"PatcherHash" desc:"Hash of current patcher"`
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
		enc := yaml.NewEncoder(f)
		cfg = getDefaultConfig()
		err = enc.Encode(cfg)
		if err != nil {
			return nil, fmt.Errorf("encode default: %w", err)
		}
		return &cfg, nil
	}

	err = yaml.NewDecoder(f).Decode(&cfg)
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
	cfg := Config{
		AutoPlay:  "false",
		AutoPatch: "false",
	}

	if len(os.Args) == 0 {
		return cfg
	}

	f, err := os.Open(os.Args[0])
	if err != nil {
		fmt.Println("open:", err)
		return cfg
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		fmt.Println("copy:", err)
		return cfg
	}
	cfg.PatcherHash = fmt.Sprintf("%x", h.Sum(nil))

	return cfg
}

// Save writes the config to disk
func (c *Config) Save() error {
	w, err := os.Create("eqemupatcher.yml")
	if err != nil {
		return fmt.Errorf("create eqemupatcher.yml: %w", err)
	}
	defer w.Close()

	enc := yaml.NewEncoder(w)
	err = enc.Encode(c)
	if err != nil {
		return fmt.Errorf("encode default: %w", err)
	}
	return nil
}
