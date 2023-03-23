package sunbeam

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/canonical/microcluster/state"

	"github.com/openstack-snaps/sunbeam-microcluster/database"
)

// GetTerraformState returns the terraform state from the database
func GetTerraformState(s *state.State) (string, error) {
	state := "{}"

	// Get the state from the database.
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		record, err := database.GetConfigItem(ctx, tx, "TerraformState")
		if err != nil {
			if strings.Contains(err.Error(), "ConfigItem not found") {
				return nil
			}
			return fmt.Errorf("Failed to fetch terraform lock: %w", err)
		}

		state = record.Value
		return nil
	})
	if err != nil {
		return "", err
	}

	return state, nil
}

// UpdateTerraformState updates the terraform state record in the database
func UpdateTerraformState(s *state.State, state string) error {
	c := database.ConfigItem{Key: "TerraformState", Value: state}

	// Add state to the database.
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		err := database.UpdateConfigItem(ctx, tx, "TerraformState", c)
		if err != nil && strings.Contains(err.Error(), "ConfigItem not found") {
			_, err = database.CreateConfigItem(ctx, tx, c)
		}
		if err != nil {
			return fmt.Errorf("Failed to record terraform state: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// GetTerraformLock returns the terraform lock from the database
func GetTerraformLock(s *state.State) (string, error) {
	lock := "{}"

	// Get the lock from the database.
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		record, err := database.GetConfigItem(ctx, tx, "TerraformLock")
		if err != nil {
			if strings.Contains(err.Error(), "ConfigItem not found") {
				return nil
			}
			return fmt.Errorf("Failed to fetch terraform lock: %w", err)
		}

		lock = record.Value
		return nil
	})
	if err != nil {
		return "", err
	}

	return lock, nil
}

// UpdateTerraformLock updates the terraform lock record in the database
func UpdateTerraformLock(s *state.State, lock string) error {
	c := database.ConfigItem{Key: "TerraformLock", Value: lock}

	// Add lock to the database.
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		err := database.UpdateConfigItem(ctx, tx, "TerraformLock", c)
		if err != nil && strings.Contains(err.Error(), "ConfigItem not found") {
			_, err = database.CreateConfigItem(ctx, tx, c)
		}
		if err != nil {
			return fmt.Errorf("Failed to record terraform lock: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
