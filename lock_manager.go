package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

type LockManager struct {
	dbClient *dynamodb.Client
	config   *LockManagerConfig
}

func NewLockManager(conf *LockManagerConfig) (*LockManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(conf.Region), config.WithSharedConfigProfile(conf.AWSProfile))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	return &LockManager{dbClient: svc, config: conf}, nil
}

func (lm *LockManager) Unlock(pattern string) error {
	// Scan the table to get all items
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(lm.config.TableName),
	}
	log.Infof("Scanning table %s\n", lm.config.TableName)

	scanResult, err := lm.dbClient.Scan(context.TODO(), scanInput)
	if err != nil {
		return fmt.Errorf("failed to scan table: %w", err)
	}

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	// Find and delete matching items
	for _, item := range scanResult.Items {
		var lockItem map[string]interface{}
		err = attributevalue.UnmarshalMap(item, &lockItem)
		if err != nil {
			return fmt.Errorf("failed to unmarshal DynamoDB item: %w", err)
		}

		lockID, ok := lockItem["LockID"].(string)
		if !ok {
			continue // skip if ID is not a string
		}

		if re.MatchString(lockID) {
			// Delete the matching item
			_, err = lm.dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(lm.config.TableName),
				Key: map[string]types.AttributeValue{
					"LockID": &types.AttributeValueMemberS{Value: lockID},
				},
			})

			if err != nil {
				return fmt.Errorf("failed to delete item with ID %s: %w", lockID, err)
			}

			log.Infof("Deleted lock with ID: %s\n", lockID)
		}
	}

	return nil
}
