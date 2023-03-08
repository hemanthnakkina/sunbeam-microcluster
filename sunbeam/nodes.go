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
                        return fmt.Errorf("Failed to fetch node: %w", err)
                }

                for _, node := range records {
                        nodes = append(nodes, types.Node{
                                Name: node.Name,
                                Role:  node.Role,
                        })
                }

                return nil
        })
        if err != nil {
                return nil, err
        }

        return nodes, nil
}

func AddNode(s *state.State, name string, role string) error {
        // Add node to the database.
        err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
                _, err := database.CreateNode(ctx, tx, database.Node{Member: s.Name(), Name: name, Role: role})
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
