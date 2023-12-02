package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/icza/gox/imagex/colorx"
	log "github.com/sirupsen/logrus"
)

func getQRCodeController(a *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("Name")
		format := c.Query("Format")
		color_hex := c.Query("Color")

		size := c.DefaultQuery("Size", "256")
		size_int, err := strconv.Atoi(size)

		if err != nil {
			log.Errorf("unable to parse request, error: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": fmt.Sprintf("Unable to parse request, error: %s", err),
			})
			return
		}

		// url that will be encoded in QR code
		url := fmt.Sprintf("http://%s/%s", a.Config.Domain, name)

		// Create the qrcode as barcode object
		qrCode, _ := qr.Encode(url, qr.M, qr.Auto)

		if format == "svg" {
			buffer := new(bytes.Buffer)

			// Make svg file of qr code
			s := svg.New(buffer)
			qs := goqrsvg.NewQrSVG(qrCode, 5)
			qs.StartQrSVG(s)
			qs.WriteQrSVG(s)
			s.End()

			// []byte representing svg file
			buffer_bytes := buffer.Bytes()

			// Change color of all black boxes by replacing the color black everywhere
			// (if it looks stupid but it works, it's not stupid)
			buffer_bytes = bytes.ReplaceAll(buffer_bytes, []byte("style=\"fill:black;stroke:none\""), []byte(fmt.Sprintf("style=\"fill:%s;stroke:none\"", color_hex)))

			// Send response
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"QR_%s.svg\"", name))
			c.Data(http.StatusOK, "image/svg-xml", buffer_bytes)

		} else if format == "png" {

			// Scale the barcode
			qrCode, _ = barcode.Scale(qrCode, size_int, size_int)

			qr_color, err := colorx.ParseHexColor(color_hex)

			if err != nil {
				log.Errorf("unable to parse request, error: %s", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"errorMessage": fmt.Sprintf("Unable to parse request, error: %s", err),
				})
				return
			}

			// For png images the QR code is colored using a ColoredQrCode,
			// which overwrites the At and ColorModel methods of Barcode
			coloredQrCode := ColoredBarcode{Barcode: qrCode, mainColor: qr_color, backgroundColor: color.White}

			buffer := new(bytes.Buffer)
			// encode the barcode as png
			png.Encode(buffer, coloredQrCode)

			// Return QR code
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"QR_%s.png\"", name))
			c.Data(http.StatusOK, "image/png", buffer.Bytes())
		}
	}
}
