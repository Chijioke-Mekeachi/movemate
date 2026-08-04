package db

import "token/models"


var DummySender = models.SenderData{
	SenderName: "loi yung",
	SenderPhone: "+23459857956",
	PickupLocation: "Iwofe road",
	PickupCity: "Port harcourt",
	PickupCountry: "Nigeria",
}

var DummyClient = models.ReceiverData{
	ReceiverName: "John Harmond",
	ReceiverPhone: "+1348983443",
	DeliveryLocation: "somewhere arround UK",
	DeliveryCity: "London",
	DeliveryCountry: "England",
}

var DummyPackage  = models.PackageData{
	PackageDescription: "Mack book laptop",
	PackageWeigth: 5, 
	PackageCategories: "Tech",
}


var Dummydb = []models.Shipment{
	{Id: "00000000", TrackingId: "00000000", Status: "pending", Sender: DummySender, Receiver: DummyClient, Package: DummyPackage, CreatedAt: "31 july", EstimatedDelivery: "1st August", DeliveryAt: "", TimeLine: ""},
	{Id: "00000001", TrackingId: "00000001", Status: "pending", Sender: DummySender, Receiver: DummyClient, Package: DummyPackage, CreatedAt: "31 july", EstimatedDelivery: "1st August", DeliveryAt: "", TimeLine: ""},
}