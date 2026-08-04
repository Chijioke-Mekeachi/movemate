package routes

import (
	"net/http"
	"token/models"
	"token/utils"

	"github.com/gin-gonic/gin"
)

func FetchAll(c *gin.Context) {
	var allTheDbData = utils.AdminFetchAll()
	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: allTheDbData})
}

func AdminLogin(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ApiResponse{Success: false, Error: "Invalid request body"})
		return
	}

	user, err := utils.AdminLogin(loginReq.Username, loginReq.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Error: "Internal server error"})
		return
	}

	if user == nil {
		c.IndentedJSON(http.StatusUnauthorized, models.ApiResponse{Success: false, Error: "Invalid credentials"})
		return
	}

	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true, Data: user})
}

func AdminLogout(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.ApiResponse{Success: true})
}
