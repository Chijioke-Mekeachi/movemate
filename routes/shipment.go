package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"token/models"
	"token/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func isAdmin(c *gin.Context) bool {
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "admin-secret"
	}
	requestToken := c.GetHeader("X-Admin-Token")
	return requestToken == adminToken
}

func PostShipment(c *gin.Context) {
	var newShipment models.ShipmentFormData
	if err := c.BindJSON(&newShipment); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ApiResponse{Success: false, Error: "Invalid parameter input."})
		return
	}

	shipment, err := utils.CreateShipment(newShipment)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not create shipment."})
		return
	}

	c.IndentedJSON(http.StatusCreated, models.ApiResponse{Success: true, Data: shipment})
}

func GetShipmentByID(c *gin.Context) {
	shipmentID := c.Param("id")
	shipment, err := utils.GetShipmentByID(shipmentID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Internal server error."})
		return
	}
	if shipment == nil {
		c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: shipment})
}

func TrackShipment(c *gin.Context) {
	trackingID := c.Param("trackingId")
	shipment, err := utils.GetShipmentByTrackingID(trackingID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Internal server error."})
		return
	}
	if shipment == nil {
		c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: shipment})
}

func GetAllShipments(c *gin.Context) {
	shipments, err := utils.GetAllShipments()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not fetch shipments."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: shipments})
}

func UpdateShipment(c *gin.Context) {
	if !isAdmin(c) {
		c.IndentedJSON(http.StatusForbidden, models.ApiResponse{Success: false, Error: "Admin access required to update shipments."})
		return
	}

	shipmentID := c.Param("id")
	var updateData models.ShipmentUpdate
	if err := c.BindJSON(&updateData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ApiResponse{Success: false, Error: "Invalid update payload."})
		return
	}

	updatedShipment, err := utils.UpdateShipment(shipmentID, updateData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not update shipment."})
		return
	}
	if updatedShipment == nil {
		c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: updatedShipment})
}

func UpdateShipmentStatus(c *gin.Context) {
	if !isAdmin(c) {
		c.IndentedJSON(http.StatusForbidden, models.ApiResponse{Success: false, Error: "Admin access required."})
		return
	}

	trackingID := c.Param("trackingId")
	var statusUpdate struct {
		Status   string `json:"status" binding:"required"`
		Location string `json:"location"`
	}

	if err := c.BindJSON(&statusUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ApiResponse{Success: false, Error: "Invalid request body."})
		return
	}

	shipment, err := utils.GetShipmentByTrackingID(trackingID)
	if err != nil || shipment == nil {
		c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	updates := models.ShipmentUpdate{
		Status: &statusUpdate.Status,
	}

	// Record a simple timeline string combining status and optional location
	timelineText := statusUpdate.Status
	if statusUpdate.Location != "" {
		timelineText = fmt.Sprintf("%s - %s", statusUpdate.Status, statusUpdate.Location)
	}
	updates.TimeLine = &timelineText

	if statusUpdate.Status == "delivered" {
		updates.DeliveryAt = &now
	}

	updatedShipment, err := utils.UpdateShipment(shipment.Id, updates)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not update shipment status."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: updatedShipment})
}

func UpdateShipmentETA(c *gin.Context) {
	if !isAdmin(c) {
		c.IndentedJSON(http.StatusForbidden, models.ApiResponse{Success: false, Error: "Admin access required."})
		return
	}

	trackingID := c.Param("trackingId")
	var etaUpdate struct {
		NewEta string `json:"newEta" binding:"required"`
	}

	if err := c.BindJSON(&etaUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ApiResponse{Success: false, Error: "Invalid request body."})
		return
	}

	shipment, err := utils.GetShipmentByTrackingID(trackingID)
	if err != nil || shipment == nil {
		c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
		return
	}

	updates := models.ShipmentUpdate{
		EstimatedDelivery: &etaUpdate.NewEta,
	}

	updatedShipment, err := utils.UpdateShipment(shipment.Id, updates)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not update ETA."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: updatedShipment})
}

func DeleteShipment(c *gin.Context) {
	if !isAdmin(c) {
		c.IndentedJSON(http.StatusForbidden, models.ApiResponse{Success: false, Error: "Admin access required."})
		return
	}

	shipmentID := c.Param("id")
	if err := utils.DeleteShipment(shipmentID); err != nil {
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusNotFound, models.ApiResponse{Success: false, Error: "Shipment not found."})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not delete shipment."})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetAnalytics(c *gin.Context) {
	if !isAdmin(c) {
		c.IndentedJSON(http.StatusForbidden, models.ApiResponse{Success: false, Error: "Admin access required."})
		return
	}

	analytics, err := utils.GetAnalytics()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Could not fetch analytics."})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: analytics})
}
