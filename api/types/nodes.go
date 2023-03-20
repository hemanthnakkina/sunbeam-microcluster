// Package types provides shared types and structs.
package types

type Nodes []Node
type Node struct {
        Name  string `json:"name" yaml:"name"`
        Role string `json:"role" yaml:"role"`
	MachineID int `json:"machineid" yaml:"machineid"`
}
