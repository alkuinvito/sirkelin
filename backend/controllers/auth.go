package controllers

import (
	"net/http"
	"os"

	"github.com/alkuinvito/malakh-api/models"
	"github.com/alkuinvito/malakh-api/utils"
	"github.com/gin-gonic/gin"
)

func UserHandler(rg *gin.RouterGroup) {
	register := rg.Group("/register")
	register.POST("/", userRegister)

	login := rg.Group("/login")
	login.POST("/", userLogin)
}

func userRegister(c *gin.Context) {
	var user models.User

	c.Bind(&user)

	if _, err := user.SaveUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Username,
			"email": user.Email,
		},
	})
}

func userLogin(c *gin.Context) {
	var user models.User

	c.Bind(&user)

	if user, err := user.ValidateLogin(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "Username and/or password incorrect",
			},
		})
	} else {
		token, err := utils.CreateToken(user.Username, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": "Failed to generate token",
				},
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt", token, 3600, "/", os.Getenv("HOST"), true, true)

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":    user.ID,
				"name":  user.Username,
				"email": user.Email,
				"picture": user.Picture,
			},
		})
	}
}