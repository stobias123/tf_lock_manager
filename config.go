package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"
)

type LockManagerConfig struct {
	TableName  string `toml:"table"`
	Region     string `toml:"region"`
	AWSProfile string `toml:"profile"`
}

func GetLockManagerConfig() (*LockManagerConfig, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configPath := filepath.Join(usr.HomeDir, ".terraform.d", "lock_manager.toml")
	var config LockManagerConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &config, nil
}

func CreateLockManagerConfig() (*LockManagerConfig, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configPath := filepath.Join(usr.HomeDir, ".terraform.d", "lock_manager.toml")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(filepath.Dir(configPath), 0755)
			if errDir != nil {
				return nil, fmt.Errorf("failed to create directory: %w", errDir)
			}
		} else {
			return nil, fmt.Errorf("config file already exists at %s", configPath)
		}
	}

	f, err := os.Create(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	prompt := promptui.Prompt{
		Label: "TableName",
		Validate: func(input string) error {
			if len(input) < 1 {
				return errors.New("TableName must have more than 0 characters")
			}
			return nil
		},
	}

	TableName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	prompt = promptui.Prompt{
		Label: "Region",
		Validate: func(input string) error {
			if len(input) < 1 {
				return errors.New("Region must have more than 0 characters")
			}
			return nil
		},
	}

	Region, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	prompt = promptui.Prompt{
		Label: "AWSProfile",
		Validate: func(input string) error {
			if len(input) < 1 {
				return errors.New("AWSProfile must have more than 0 characters")
			}
			return nil
		},
	}

	AWSProfile, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	config := LockManagerConfig{
		TableName:  TableName,
		Region:     Region,
		AWSProfile: AWSProfile,
	}

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return nil, fmt.Errorf("failed to encode config: %w", err)
	}

	return GetLockManagerConfig()
}
