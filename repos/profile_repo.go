package repos

import (
	"context"
	"errors"
	"fmt"
	"my-go-backend/internal"
	"my-go-backend/models"
	"sync"
	"time"
)

type ProfileRepository interface {
	GetProfile(lastSync int64) (*models.Profile, error)
	UpdateProfile(profile models.Profile) (int64, error)
}

type profileRepositoryImpl struct {
	firestore *internal.FirestoreClient
	cache     sync.Map
}

func (p *profileRepositoryImpl) GetProfile(lastSync int64) (*models.Profile, error) {
	if cachedDTO := p.loadFromCache(); cachedDTO != nil {
		if cachedDTO.LastUpdate >= lastSync {
			return toModel(cachedDTO), nil
		}
		return nil, nil
	}

	freshDTO, err := p.loadFromRepository()
	if err != nil {
		return nil, err
	}

	p.saveToCache(freshDTO)

	if freshDTO.LastUpdate >= lastSync {
		return toModel(freshDTO), nil
	}
	return nil, nil
}

func (p *profileRepositoryImpl) loadFromRepository() (*ProfileDTO, error) {
	ctx := context.Background()
	doc, err := p.firestore.Get().
		Collection("sync").
		Doc("profile").
		Get(ctx)

	if err != nil {
		return nil, err
	}

	var dto ProfileDTO
	if err := doc.DataTo(&dto); err != nil {
		return nil, errors.New("can not parse profile object")
	}

	return &dto, nil
}

func (p *profileRepositoryImpl) loadFromCache() *ProfileDTO {
	if value, ok := p.cache.Load("profile"); ok {
		fmt.Println("prfoile -> From cache")
		return value.(*ProfileDTO)
	} else {
		fmt.Println("prfoile -> empty cache")
		return nil
	}
}

func (p *profileRepositoryImpl) saveToCache(dto *ProfileDTO) {
	p.cache.Store("profile", dto)
}

func (p *profileRepositoryImpl) UpdateProfile(profile models.Profile) (int64, error) {
	now := time.Now().Unix()
	ctx := context.Background()
	_, err := p.firestore.Get().
		Collection("sync").
		Doc("profile").
		Set(ctx, toDTO(profile, now))

	if err != nil {
		return 0, err
	}

	p.cache.Clear()

	return now, nil
}

func NewProfileRepository(firestore *internal.FirestoreClient) ProfileRepository {
	return &profileRepositoryImpl{
		firestore: firestore,
	}
}

type ProfileDTO struct {
	LastUpdate  int64      `firestore:"lastUpdate"`
	Description string     `firestore:"description"`
	Skills      []SkillDTO `firestore:"skills"`
}

type SkillDTO struct {
	Name  string `firestore:"name"`
	Level int    `firestore:"level"`
}

func toDTO(p models.Profile, lastUpdate int64) ProfileDTO {
	skills := make([]SkillDTO, len(p.Skills))
	for i, s := range p.Skills {
		skills[i] = SkillDTO{
			Name:  s.Name,
			Level: s.Level,
		}
	}

	return ProfileDTO{
		Description: p.Description,
		Skills:      skills,
		LastUpdate:  lastUpdate,
	}
}

func toModel(p *ProfileDTO) *models.Profile {
	skills := make([]models.Skill, len(p.Skills))
	for i, s := range p.Skills {
		skills[i] = models.Skill{
			Name:  s.Name,
			Level: s.Level,
		}
	}

	fmt.Println("DESC:" + p.Description)

	return &models.Profile{
		Description: p.Description,
		Skills:      skills,
	}
}
