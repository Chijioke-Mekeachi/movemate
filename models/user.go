package models

type SenderData struct {
	SenderName     string `json:"senderName"`
	SenderPhone    string `json:"senderPhone"`
	PickupLocation string `json:"pickupLocation"`
	PickupCity     string `json:"pickupCity"`
	PickupCountry  string `json:"pickupCountry"`
}

type ReceiverData struct {
	ReceiverName     string `json:"receiverName"`
	ReceiverPhone    string `json:"receiverPhone"`
	DeliveryLocation string `json:"deliveryLocation"`
	DeliveryCity     string `json:"deliveryCity"`
	DeliveryCountry  string `json:"deliveryCountry"`
}

type Admin struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}
