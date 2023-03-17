// Package types provides shared types and structs.
package types

type JujuUsers []JujuUser
type JujuUser struct {
        Username  string `json:"username" yaml:"username"`
        Token string `json:"token" yaml:"token"`
}
