package controllers

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var res []models.User
	var err error

	client, _ := utils.NewFirebaseClient(c)
	session, _ := utils.GetSessionFromContext(c)
	uid, _ := utils.GetIDFromSession(client, c, session)

	res, err = models.GetUsers(uid)
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
