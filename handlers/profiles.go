package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	profilesdto "waysgallery/dto/profiles"
	dto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerProfile struct {
	ProfileRepository repositories.ProfileRepository
}

func HandlerProfile(ProfileRepository repositories.ProfileRepository) *handlerProfile {
	return &handlerProfile{ProfileRepository}
}

func (h *handlerProfile) FindProfiles(c echo.Context) error {
	profiles, err := h.ProfileRepository.FindProfiles()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Getting all data success", Data: profiles})
}

func (h *handlerProfile) GetProfile(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	var profile models.Profile
	profile, err := h.ProfileRepository.GetProfile(int(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Get Data Profile Success", Data: profile})
}

func (h *handlerProfile) UpdateProfile(c echo.Context) error {
	filepath := c.Get("datafile").(string)
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	request := profilesdto.ProfileRequest{
		Name:     c.FormValue("name"),
		Greeting: c.FormValue("greeting"),
		Image:    filepath,
	}

	var cntx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")
	cloud, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	res, err := cloud.Upload.Upload(cntx, filepath, uploader.UploadParams{Folder: "waysgallery"})
	if err != nil {
		fmt.Println(err.Error())
	}
	profile, err := h.ProfileRepository.GetProfile(int(userId))

	if err != nil {

		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})

	}

	if request.Name != "" {
		profile.Name = request.Name
	}

	if request.Greeting != "" {
		profile.Greeting = request.Greeting
	}

	if request.Image != "" {
		profile.Image = res.SecureURL
	}
	profile.ImagePublicID = res.PublicID
	profile.UpdateAt = time.Now()

	data, err := h.ProfileRepository.UpdateProfile(profile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Update Success", Data: convertResponseProfile(data)})

}

func convertResponseProfile(u models.Profile) profilesdto.ProfileResponse {
	return profilesdto.ProfileResponse{
		ID:            u.ID,
		Name:          u.Name,
		Greeting:      u.Greeting,
		Image:         u.Image,
		ImagePublicID: u.ImagePublicID,
	}
}
