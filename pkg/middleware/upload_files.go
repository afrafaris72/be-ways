package middleware

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func UploadFiles(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var method = c.Request().Method
		form, err := c.MultipartForm()

		if err != nil {
			if method == "PATCH" && err.Error() == "http: no such file" {
				c.Set("dataFile", "")
				return next(c)
			}
			if method == "POST" && err.Error() == "http: no such file" {
				c.Set("dataFile", "")
				return next(c)
			}
			return c.JSON(http.StatusBadRequest, err)
		}

		var dataFiles []string

		for i := 0; i < 5; i++ {
			file := form.File["image"+strconv.Itoa(i+1)]
			if len(file) == 0 {
				dataFiles = append(dataFiles, "")
				continue
			}
			extens := filepath.Ext(file[0].Filename)
			if extens == ".png" || extens == ".jpg" || extens == ".jpeg" || extens == ".webp" {
				src, err := file[0].Open()
				if err != nil {
					return c.JSON(http.StatusBadRequest, err)
				}
				defer src.Close()

				tempFile, err := ioutil.TempFile("uploads", "image-*.png")
				if err != nil {
					return c.JSON(http.StatusBadRequest, err)
				}
				defer tempFile.Close()

				if _, err = io.Copy(tempFile, src); err != nil {
					return c.JSON(http.StatusBadRequest, err)
				}
				dataFiles = append(dataFiles, tempFile.Name())
			} else {
				return c.JSON(http.StatusBadRequest, "ONLY => (.png, .jpg, .jpeg, .webp)")
			}

		}
		c.Set("dataFiles", dataFiles)

		return next(c)
	}
}
