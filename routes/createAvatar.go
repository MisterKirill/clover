package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAvatarPayload struct {
	Success bool `json:"success"`
	Error string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

func CreateAvatar(c *gin.Context) {
	header, err := c.FormFile("avatar")
	if err != nil {
		log.Panic(err)
	}
	
	file, err := header.Open()
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)
	image, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Panic(err)
	}

	defer buf.Reset()

	if image.Height != image.Width {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error": "Bad image aspect ratio (should be 1:1)",
			"error_code": "BAD_ASPECT_RATIO",
		})
		return
	}

	if format != "jpg" && format != "jpeg" && format != "png" {
		contentType := header.Header.Get("Content-Type")

		if contentType != "image/png" && contentType != "image/jpeg" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"error": "Bad image format (only png and jpeg allowed)",
				"error_code": "BAD_FORMAT",
			})
			return
		}
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(buf.Bytes())); err != nil {
		log.Panic(err)
	}

	avatarID := hex.EncodeToString(hash.Sum(nil))

	c.SaveUploadedFile(header, "avatars/" + avatarID + "." + format)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"avatarID": avatarID,
	})
}
