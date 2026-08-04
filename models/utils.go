package models

type Homestruct struct {
	Msg string `json:"msg"`
}

type ErrorStruct struct {
	Msg string `json:"msg"`
}

type Coord struct {
	Lat int `json:"lat"`
	Lng int `json:"lng"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type Analytics struct {
	Total              int           `json:"total"`
	Delivered          int           `json:"delivered"`
	InTransit          int           `json:"inTransit"`
	Pending            int           `json:"pending"`
	StatusDistribution []StatusCount `json:"statusDistribution"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}
