package model

func init() {
	Register(&UserBind{})
}

// UserBind represents user_bind domain model
type UserBind struct {
	Base
	ID       int    `json:"id" pg:"id"`
	UserID   int    `json:"user_id" pg:"user_id"`
	BindType string `json:"bind_type" pg:"bind_type"`
	BindID   int    `json:"bind_id" pg:"bind_id"`
}
