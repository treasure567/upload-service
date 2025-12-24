package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type R2Config struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	BucketName      string `yaml:"bucket_name"`
	Endpoint        string `yaml:"endpoint"`
	Region          string `yaml:"region"`
	CustomDomain    string `yaml:"custom_domain"`
}

type UploadConfig struct {
	MaxFileSize int64  `yaml:"max_file_size"`
	TempDir     string `yaml:"temp_dir"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	R2     R2Config     `yaml:"r2"`
	Upload UploadConfig `yaml:"upload"`
}

var AppConfig *Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	AppConfig = &cfg
	return nil
}
