package models

type PackageData struct {
	PackageDescription string `json:"packageDescription" bson:"packageDescription"`
	PackageWeigth      int    `json:"packageWeight" bson:"packageWeight"`
	PackageCategories  string `json:"packageCategories" bson:"packageCategories"`
}

type ShipmentLocation struct {
	City       string `json:"city" bson:"city"`
	Country    string `json:"country" bson:"country"`
	Coordinate Coord  `json:"coordinate" bson:"coordinate"`
}

type ShipmentTimeline struct {
	Status      string `json:"status" bson:"status"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	TimeStamp   string `json:"timestamp" bson:"timestamp"`
	Location    string `json:"location" bson:"location"`
	Completed   bool   `json:"completed" bson:"completed"`
}

type Shipment struct {
	Id                string       `json:"id" bson:"id"`
	TrackingId        string       `json:"trackingId" bson:"trackingId"`
	Status            string       `json:"status" bson:"status"`
	Sender            SenderData   `json:"sender" bson:"sender"`
	Receiver          ReceiverData `json:"receiver" bson:"receiver"`
	Package           PackageData  `json:"package" bson:"package"`
	CreatedAt         string       `json:"createdAt" bson:"createdAt"`
	EstimatedDelivery string       `json:"estimatedDelivery" bson:"estimatedDelivery"`
	DeliveryAt        string       `json:"deliveryAt" bson:"deliveryAt"`
	TimeLine          string       `json:"timeline" bson:"timeline"`
}

type ShipmentFormData struct {
	Sender   SenderData   `json:"sender" bson:"sender"`
	Receiver ReceiverData `json:"receiver" bson:"receiver"`
	Package  PackageData  `json:"package" bson:"package"`
}

type ShipmentUpdate struct {
	Status            *string       `json:"status,omitempty" bson:"status,omitempty"`
	Sender            *SenderData   `json:"sender,omitempty" bson:"sender,omitempty"`
	Receiver          *ReceiverData `json:"receiver,omitempty" bson:"receiver,omitempty"`
	Package           *PackageData  `json:"package,omitempty" bson:"package,omitempty"`
	EstimatedDelivery *string       `json:"estimatedDelivery,omitempty" bson:"estimatedDelivery,omitempty"`
	DeliveryAt        *string       `json:"deliveryAt,omitempty" bson:"deliveryAt,omitempty"`
	TimeLine          *string       `json:"timeline,omitempty" bson:"timeline,omitempty"`
}

type TrackedShipment struct {
	TrackingId string `json:"trackingId" bson:"trackingId"`
	Status     string `json:"status" bson:"status"`
	TrackedAt  string `json:"trackedAt" bson:"trackedAt"`
}
