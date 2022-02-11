package response

import uuid "github.com/satori/go.uuid"

// 响应模型

type Admin struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"admin"`
}

type Account struct {
	ID   int       `json:"id" example:"1" format:"int64"`
	Name string    `json:"name" example:"account"`
	UUID uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

type Bottle struct {
	ID      int     `json:"id" example:"1"`
	Name    string  `json:"name" example:"bottle"`
	Account Account `json:"account"`
}
