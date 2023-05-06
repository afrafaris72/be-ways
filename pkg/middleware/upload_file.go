package middleware

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func UploadFile(newfile echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var method = c.Request().Method
		file, err := c.FormFile("image")

		if err != nil {
			if method == "PATCH" && err.Error() == "http: no such file" {
				c.Set("dataFile", "")
				return newfile(c)
			}
		}
		if err != nil {
			if method == "POST" && err.Error() == "http: no such file" {
				c.Set("datafile", "")
				return newfile(c)
			}
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		extension := filepath.Ext(file.Filename)
		if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".webp" {
			src, err := file.Open()
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

			data := tempFile.Name()

			c.Set("dataFile", data)
			return newfile(c)
		} else {
			return c.JSON(http.StatusBadRequest, "File Extension Must Be (.png, .jpg, .jpeg, .webp)")
		}
	}
}
