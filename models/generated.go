// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Node interface {
	IsNode()
}

type CreateGameInvitationInput struct {
	Type      GameType `json:"type"`
	BoardSize int      `json:"boardSize"`
}

type CreateIdentityInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}