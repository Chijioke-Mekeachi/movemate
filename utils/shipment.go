package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"token/db"
	"token/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func generateShipmentID(seed string) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%d", strings.TrimSpace(seed), time.Now().UnixNano())))
	return hex.EncodeToString(hash[:])
}

// generateFriendlyTrackingID creates a tracking id in the format MM-LX-XXXXX
// where XXXXX is a 5 character alphanumeric string. It is returned in lower-case
// for consistent storage/lookup in the database.
func generateFriendlyTrackingID() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	for i := range b {
		b[i] = chars[int(time.Now().UnixNano()+int64(i))%len(chars)]
	}
	// mix with a small random-ish component using md5 of timestamp to reduce collisions
	hash := md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	hexpart := hex.EncodeToString(hash[:])
	// take alphanumeric characters from hexpart until we have 5 chars if needed
	var extra []byte
	for i := 0; i < len(hexpart) && len(extra) < 5; i++ {
		c := hexpart[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			extra = append(extra, c)
		}
	}
	// fallback: if extra is short, fill with '0'
	for len(extra) < 5 {
		extra = append(extra, '0')
	}

	// produce final 5-char code mixing b and extra then uppercase and lowercase for storage
	code := fmt.Sprintf("%c%c%c%c%c", b[0], b[1], extra[0], extra[1], extra[2])
	return strings.ToLower(fmt.Sprintf("MM-LX-%s", code))
}

func CreateShipment(shipForm models.ShipmentFormData) (*models.Shipment, error) {
	if db.ShipmentCollection == nil {
		return nil, mongo.ErrNilDocument
	}

	shipmentID := strings.ToLower(generateShipmentID(shipForm.Sender.SenderName))
	now := time.Now().UTC().Format(time.RFC3339)
	shipment := models.Shipment{
		Id:                shipmentID,
		TrackingId:        generateFriendlyTrackingID(),
		Status:            "pending",
		Sender:            shipForm.Sender,
		Receiver:          shipForm.Receiver,
		Package:           shipForm.Package,
		CreatedAt:         now,
		EstimatedDelivery: now,
		DeliveryAt:        "",
		TimeLine:          "created",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ShipmentCollection.InsertOne(ctx, shipment)
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func GetShipmentByID(id string) (*models.Shipment, error) {
	if db.ShipmentCollection == nil {
		return nil, mongo.ErrNilDocument
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var shipment models.Shipment
	err := db.ShipmentCollection.FindOne(ctx, bson.M{"id": id}).Decode(&shipment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &shipment, nil
}

func GetAllShipments() ([]models.Shipment, error) {
	if db.ShipmentCollection == nil {
		return nil, mongo.ErrNilDocument
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.ShipmentCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shipments []models.Shipment
	for cursor.Next(ctx) {
		var shipment models.Shipment
		if err := cursor.Decode(&shipment); err != nil {
			return nil, err
		}
		shipments = append(shipments, shipment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return shipments, nil
}

func UpdateShipment(id string, update models.ShipmentUpdate) (*models.Shipment, error) {
	if db.ShipmentCollection == nil {
		return nil, mongo.ErrNilDocument
	}

	updateFields := bson.M{}
	if update.Status != nil {
		updateFields["status"] = *update.Status
	}
	if update.Sender != nil {
		updateFields["sender"] = update.Sender
	}
	if update.Receiver != nil {
		updateFields["receiver"] = update.Receiver
	}
	if update.Package != nil {
		updateFields["package"] = update.Package
	}
	if update.EstimatedDelivery != nil {
		updateFields["estimatedDelivery"] = *update.EstimatedDelivery
	}
	if update.DeliveryAt != nil {
		updateFields["deliveryAt"] = *update.DeliveryAt
	}
	if update.TimeLine != nil {
		updateFields["timeline"] = *update.TimeLine
	}

	if len(updateFields) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedShipment models.Shipment
	err := db.ShipmentCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).Decode(&updatedShipment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &updatedShipment, nil
}

func DeleteShipment(id string) error {
	if db.ShipmentCollection == nil {
		return mongo.ErrNilDocument
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ShipmentCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func GetShipmentByTrackingID(trackingId string) (*models.Shipment, error) {
	if db.ShipmentCollection == nil {
		return nil, mongo.ErrNilDocument
	}

	trackingId = strings.ToLower(strings.TrimSpace(trackingId))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var shipment models.Shipment
	err := db.ShipmentCollection.FindOne(ctx, bson.M{"trackingId": trackingId}).Decode(&shipment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &shipment, nil
}

func GetAnalytics() (*models.Analytics, error) {
	shipments, err := GetAllShipments()
	if err != nil {
		return nil, err
	}

	analytics := &models.Analytics{
		Total:              len(shipments),
		StatusDistribution: []models.StatusCount{},
	}

	statusCounts := make(map[string]int)

	for _, shipment := range shipments {
		statusCounts[shipment.Status]++

		switch shipment.Status {
		case "delivered":
			analytics.Delivered++
		case "in_transit", "out_for_delivery":
			analytics.InTransit++
		case "pending", "processing":
			analytics.Pending++
		}
	}

	for status, count := range statusCounts {
		analytics.StatusDistribution = append(analytics.StatusDistribution, models.StatusCount{
			Status: status,
			Count:  count,
		})
	}

	return analytics, nil
}
