package request

// 请求模型

type AddAccount struct {
	Name string `json:"name" example:"account"`
}
