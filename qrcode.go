package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

func getQRCodeController(a *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("Name")

		var png []byte
		url := fmt.Sprintf("%s/%s", a.Config.Domain, name)
		png, err := qrcode.Encode(url, qrcode.Medium, 512)

		if err != nil {
			log.Errorf("unable to parse request, error: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": fmt.Sprintf("Unable to parse request, error: %s", err),
			})
			return
		}
		// Return QR code
		c.Data(http.StatusOK, "image/png", png)

	}
}
