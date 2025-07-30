package handlers

import (
	"my-go-backend/models"
	"my-go-backend/repos"
	"my-go-backend/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SyncHandler struct {
	profileRepository    repos.ProfileRepository
	experienceRepository repos.ExperienceRepository
}

func NewSyncHandler(profileRepository repos.ProfileRepository, experienceRepository repos.ExperienceRepository) *SyncHandler {
	return &SyncHandler{
		profileRepository:    profileRepository,
		experienceRepository: experienceRepository,
	}
}

func (s *SyncHandler) Sync(c *fiber.Ctx) error {
	queryParamSince := c.Query("last_sync")
	if queryParamSince == "" {
		return c.Status(400).SendString("Brak parametru last_sync")
	}

	lastSync, err := strconv.ParseInt(queryParamSince, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Niepoprawny format last_sync")
	}
	now := time.Now().Unix()
	profile, err := s.profileRepository.GetProfile(lastSync)
	if err != nil {
		return c.Status(500).SendString("internal error")
	}
	experience, err := s.experienceRepository.GetExperience(lastSync)
	if err != nil {
		return c.Status(500).SendString("internal error")
	}

	response := make(map[string]any)
	response["last_sync"] = now
	if profile != nil {
		response["profile"] = *profile
	}
	if experience != nil {
		response["experience"] = *experience
	}

	return c.JSON(response)
}

type profileBody struct {
	Description string         `json:"description"`
	Skills      []models.Skill `json:"skills"`
}

func (s *SyncHandler) UpdateProfile(c *fiber.Ctx) error {
	var body profileBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if len(body.Description) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid param description")
	}
	if len(body.Skills) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid param skills")
	}
	profile := models.Profile{
		Description: body.Description,
		Skills:      body.Skills,
	}
	_, err := s.profileRepository.UpdateProfile(profile)
	if err != nil {
		return err
	}
	return nil

}

func (s *SyncHandler) UpdateExperience(c *fiber.Ctx) error {
	var body models.Experience
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if len(body.Companies) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty data")
	}

	_, err := s.experienceRepository.UpdateExperience(body)
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncHandler) UpdateProjects(c *fiber.Ctx) error {
	time := time.Now().Unix()
	textualTime := utils.Int64ToString(time)
	return c.JSON("Projects ok -> " + textualTime)
}
