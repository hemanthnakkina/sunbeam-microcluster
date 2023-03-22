package sunbeam

import (
        "context"
        "database/sql"
        "fmt"

        "github.com/canonical/microcluster/state"

        "github.com/openstack-snaps/sunbeam-microcluster/api/types"
        "github.com/openstack-snaps/sunbeam-microcluster/database"
)

func ListNodes(s *state.State) (types.Nodes, error) {
        nodes := types.Nodes{}

        // Get the nodes from the database.
        err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
                records, err := database.GetNodes(ctx, tx)
                if err != nil {
                        return fmt.Errorf("Failed to fetch nodes: %w", err)
                }

                for _, node := range records {
                        nodes = append(nodes, types.Node{
                                Name: node.Name,
                                Role:  node.Role,
				MachineID: node.MachineID,
                        })
                }

                return nil
        })
        if err != nil {
                return nil, err
        }

        return nodes, nil
}

func GetNode(s *state.State, name string) (types.Node, error) {
	node := types.Node{}
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		record, err := database.GetNode(ctx, tx, name)
		if err != nil {
			return err
		}

		node.Name = record.Name
		node.Role = record.Role
		node.MachineID = record.MachineID

		return nil
	})

	return node, err
}

func AddNode(s *state.State, name string, role string, machineid int) error {
        // Add node to the database.
        err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		_, err := database.CreateNode(ctx, tx, database.Node{Member: s.Name(), Name: name, Role: role, MachineID: machineid})
                if err != nil {
                        return fmt.Errorf("Failed to record node: %w", err)
                }

                return nil
        })
        if err != nil {
                return err
        }

        return nil
}

func UpdateNode(s *state.State, name string, role string, machineid int) error {
        // Update node to the database.
        err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		node, err := database.GetNode(ctx, tx, name)
		if err != nil {
			return fmt.Errorf("Failed to retrieve node details: %w", err)
		}

		if role == "" {
			role = node.Role
		}
		if machineid == 0 {
			machineid = node.MachineID
		}

		err = database.UpdateNode(ctx, tx, name, database.Node{Member: s.Name(), Name: name, Role: role, MachineID: machineid})
                if err != nil {
                        return fmt.Errorf("Failed to update record node: %w", err)
                }

		node, err = database.GetNode(ctx, tx, name)
                return nil
        })
        if err != nil {
                return err
        }

        return nil
}

func DeleteNode(s *state.State, name string) error {
        // Delete node from the database.
        err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		err := database.DeleteNode(ctx, tx, name)
                if err != nil {
                        return fmt.Errorf("Failed to delete node: %w", err)
                }

                return nil
        })
        if err != nil {
                return err
        }

        return nil
}
