package sunbeam

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/canonical/microcluster/state"

	"github.com/openstack-snaps/sunbeam-microcluster/database"
)

func GetConfig(s *state.State, key string) (string, error) {
	var value string

	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		record, err := database.GetConfigItem(ctx, tx, key)
		if err != nil {
			return err
		}
		value = record.Value
		return nil
	})

	if err != nil {
		return "", err
	}

	return value, nil
}

func CreateConfig(s *state.State, key string, value string) error {

	return s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		_, err := database.CreateConfigItem(ctx, tx, database.ConfigItem{Key: key, Value: value})
		if err != nil {
			return fmt.Errorf("Failed to record config item: %w", err)
		}
		return nil
	})
}

func UpdateConfig(s *state.State, key string, value string) error {
	configItem := database.ConfigItem{Key: key, Value: value}

	return s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		err := database.UpdateConfigItem(ctx, tx, key, configItem)
		if err != nil && strings.Contains(err.Error(), "ConfigItem not found") {
			_, err = database.CreateConfigItem(ctx, tx, configItem)
		}
		if err != nil {
			return fmt.Errorf("Failed to record config item: %w", err)
		}

		return nil
	})
}

func DeleteConfig(s *state.State, key string) error {
	return s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		return database.DeleteConfigItem(ctx, tx, key)
	})
}
