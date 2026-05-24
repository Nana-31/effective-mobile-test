package handler

type CreateSubscriptionInput struct {
	UserId      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" format:"uuid"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	StartDate   string `json:"start_date" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"07-2026"`
}

type UpdateSubscriptionInput struct {
	ID          string `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d" format:"uuid"`
	UserId      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" format:"uuid"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"450"`
	StartDate   string `json:"start_date" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"08-2026"`
}

type TotalCostResponse struct {
	Total int `json:"total" example:"1200"`
}
