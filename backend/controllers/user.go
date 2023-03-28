package controllers

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var req models.GetUsersParam
	var res []models.User
	var err error

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "query must between 3 and 16 characters",
			},
		})
		return
	}

	res, err = models.GetUsers(req.Fullname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "failed retrieving users",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": res,
		},
	})
}
