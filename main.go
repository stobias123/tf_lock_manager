package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lockmanager [command]",
		Short: "LockManager is a tool for managing DynamoDB locks.",
	}

	config, err := GetLockManagerConfig()
	if err != nil {
		log.Errorf("No Config Found - Creating one")
		config, err = CreateLockManagerConfig()
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
	}

	var unlockCmd = &cobra.Command{
		Use:   "unlock [pattern]",
		Short: "Unlock deletes DynamoDB locks matching a pattern",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pattern := args[0]
			lm, err := NewLockManager(config)
			if err != nil {
				log.Fatalf("Error creating LockManager: %v", err)
			}

			err = lm.Unlock(pattern)
			if err != nil {
				log.Fatalf("Error unlocking: %v", err)
			}
		},
	}

	rootCmd.AddCommand(unlockCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
