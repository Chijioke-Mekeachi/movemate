package utils

import (
	"os"
	"token/models"
)

func AdminFetchAll() []models.Shipment {
	shipments, err := GetAllShipments()
	if err != nil {
		return []models.Shipment{}
	}
	return shipments
}

func AdminLogin(username, password string) (*models.AdminUser, error) {
	// Demo credentials - in production, use proper authentication
	validCredentials := map[string]map[string]string{
		"admin": {
			"password": "admin123",
			"role":     "admin",
			"name":     "System Administrator",
		},
		"staff": {
			"password": "staff123",
			"role":     "staff",
			"name":     "Staff Member",
		},
	}

	credentials, exists := validCredentials[username]
	if !exists || credentials["password"] != password {
		return nil, nil // Return nil to indicate authentication failed
	}

	return &models.AdminUser{
		Id:       username,
		Username: username,
		Role:     credentials["role"],
		Name:     credentials["name"],
	}, nil
}

func ValidateAdminToken(token string) bool {
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "admin-secret"
	}
	return token == adminToken
}
