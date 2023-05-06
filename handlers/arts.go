package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"
	artsdto "waysgallery/dto/arts"
	dto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerArt struct {
	ArtRepository repositories.ArtRepository
}

func HandlerArt(ArtRepository repositories.ArtRepository) *handlerArt {
	return &handlerArt{ArtRepository}
}

func (h *handlerArt) FindArts(c echo.Context) error {
	arts, err := h.ArtRepository.FindArts()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Get All Data Success", Data: arts})
}

func (h *handlerArt) GetArt(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	art, err := h.ArtRepository.GetArt(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Get Art Data Completed", Data: art})
}

func (h *handlerArt) CreateArt(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)
	filepath := c.Get("dataFile").(string)

	var cntx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")
	cloud, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	respons, err := cloud.Upload.Upload(cntx, filepath, uploader.UploadParams{Folder: "waysgallery"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	art := models.Art{
		Image:         respons.SecureURL,
		ImagePublicID: respons.PublicID,
		ProfileID:     int(userId),
	}
	data, err := h.ArtRepository.CreateArt(art)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Art Created", Data: convertResponseArt(data)})
}

func (h *handlerArt) DeleteArt(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	art, err := h.ArtRepository.GetArt(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.ArtRepository.DeleteArt(art)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Art data deleted successfully", Data: convertResponseArt(data)})
}

func convertResponseArt(u models.Art) artsdto.ArtResponse {
	return artsdto.ArtResponse{
		ID:            u.ID,
		Image:         u.Image,
		ImagePublicID: u.ImagePublicID,
	}
}