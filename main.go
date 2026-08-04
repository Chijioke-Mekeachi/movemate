package main

import (
	"net/http"
	"token/models"
	"token/routes"

	"github.com/gin-gonic/gin"
)

func HomeApi(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.Homestruct{Msg: "Server is UP"})
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Admin-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Apply CORS middleware
	router.Use(CORSMiddleware())

	// Health check
	router.GET("/", HomeApi)

	// Shipment endpoints
	router.POST("/shipment", routes.PostShipment)
	router.POST("/shipments", routes.PostShipment)
	router.GET("/shipments", routes.GetAllShipments)
	router.GET("/shipments/:id", routes.GetShipmentByID)
	router.GET("/track/:trackingId", routes.TrackShipment)
	router.PUT("/shipments/:id", routes.UpdateShipment)
	router.DELETE("/shipments/:id", routes.DeleteShipment)
	router.PUT("/shipment/status/:trackingId", routes.UpdateShipmentStatus)
	router.PUT("/shipment/eta/:trackingId", routes.UpdateShipmentETA)

	// Admin endpoints
	router.GET("/admin/shipment", routes.FetchAll)
	router.GET("/admin/shipments", routes.FetchAll)
	router.GET("/admin/analytics", routes.GetAnalytics)
	router.POST("/admin/login", routes.AdminLogin)
	router.POST("/admin/logout", routes.AdminLogout)

	router.Run("localhost:8081")
}
