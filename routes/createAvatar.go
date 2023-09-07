package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
)

type Payload struct {
	Success bool `json:"success"`
	Error string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

func CreateAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(5 << 20)

	file, header, err := r.FormFile("avatar")
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
		json.NewEncoder(w).Encode(Payload{
			Success:   false,
			Error:     "Bad image aspect ratio (should be 1:1)",
			ErrorCode: "BAD_ASPECT_RATIO",
		})
		return
	}

	if format != "jpg" && format != "jpeg" && format != "png" {
		contentType := header.Header.Get("Content-Type")

		if contentType != "image/png" && contentType != "image/jpeg" {
			json.NewEncoder(w).Encode(Payload{
				Success:   false,
				Error:     "Bad image format (only png and jpeg allowed)",
				ErrorCode: "BAD_FORMAT",
			})
			return
		}
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(buf.Bytes())); err != nil {
		log.Panic(err)
	}

	fileName := hex.EncodeToString(hash.Sum(nil))

	dst, err := os.Create("avatars/" + fileName + "." + format)
	if err != nil {
		log.Panic(err)
	}

	defer dst.Close()

	if _, err := io.Copy(dst, bytes.NewReader(buf.Bytes())); err != nil {
		log.Panic(err)
		return
	}

	json.NewEncoder(w).Encode(Payload{
		Success: true,
	})
}
