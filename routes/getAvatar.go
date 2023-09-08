package routes

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

var CheckAvatarID = regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString

func GetAvatar(c *gin.Context) {
	avatarID := c.Param("avatarID")

	if !CheckAvatarID(avatarID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error": "Bad avatar ID",
			"error_code": "BAD_ID",
		})
		return
	}

	if _, err := os.Stat("avatars/" + avatarID + ".png"); err == nil {
		c.File("avatars/" + avatarID + ".png")
		return
	}

	if _, err := os.Stat("avatars/" + avatarID + ".jpg"); err == nil {
		c.File("avatars/" + avatarID + ".jpg")
		return
	}

	if _, err := os.Stat("avatars/" + avatarID + ".jpeg"); err == nil {
		c.File("avatars/" + avatarID + ".jpeg")
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error": "Avatar not found",
		"error_code": "BAD_ID",
	})
}
