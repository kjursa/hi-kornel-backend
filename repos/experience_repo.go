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

type ExperienceRepository interface {
	GetExperience(lastSync int64) (*models.Experience, error)
	UpdateExperience(exp models.Experience) (int64, error)
}

type experienceRepositoryImpl struct {
	firestore *internal.FirestoreClient
	cache     sync.Map
}

var CacheKey = "experience"

func (e *experienceRepositoryImpl) GetExperience(lastSync int64) (*models.Experience, error) {
	if cached := e.loadFromCache(); cached != nil {
		if cached.LastUpdate >= lastSync {
			return toExperience(cached), nil
		}
		return nil, nil
	}

	fresh, err := e.loadFromRepository()
	if err != nil {
		return nil, err
	}

	e.saveToCache(fresh)

	if fresh.LastUpdate >= lastSync {
		return toExperience(fresh), nil
	}
	return nil, nil
}

func (p *experienceRepositoryImpl) loadFromRepository() (*ExperienceFirestoreDoc, error) {
	ctx := context.Background()
	doc, err := p.firestore.Get().
		Collection("sync").
		Doc("experience").
		Get(ctx)

	if err != nil {
		return nil, err
	}

	var experienceDoc ExperienceFirestoreDoc
	if err := doc.DataTo(&experienceDoc); err != nil {
		return nil, errors.New("can not parse experience object")
	}

	return &experienceDoc, nil
}

func (p *experienceRepositoryImpl) loadFromCache() *ExperienceFirestoreDoc {
	if value, ok := p.cache.Load(CacheKey); ok {
		fmt.Println("experience -> From cache")
		return value.(*ExperienceFirestoreDoc)
	} else {
		fmt.Println("prfoile -> empty cache")
		return nil
	}
}

func (p *experienceRepositoryImpl) saveToCache(dto *ExperienceFirestoreDoc) {
	p.cache.Store(CacheKey, dto)
}

func (p *experienceRepositoryImpl) UpdateExperience(exp models.Experience) (int64, error) {
	now := time.Now().Unix()
	ctx := context.Background()
	_, err := p.firestore.Get().
		Collection("sync").
		Doc("experience").
		Set(ctx, toExperienceDoc(exp, now))

	if err != nil {
		return 0, err
	}

	p.cache.Clear()

	return now, nil
}

func toExperienceDoc(model models.Experience, lastUpdate int64) *ExperienceFirestoreDoc {
	companies := make([]CompanyFirestore, len(model.Companies))
	for i, v := range model.Companies {
		companies[i] = CompanyFirestore{
			Name:        v.Name,
			Description: v.Description,
			Period:      v.Period,
			Projects:    toFirestoreProjects(v.Projects),
		}
	}
	experienceFirestore := ExperienceFirestore{
		Name:      model.Name,
		Title:     model.Title,
		Companies: companies,
	}
	return &ExperienceFirestoreDoc{
		LastUpdate: lastUpdate,
		Content:    experienceFirestore,
	}
}

func toFirestoreProjects(model []models.Project) []ProjectFirestore {
	projects := make([]ProjectFirestore, len(model))
	for i, v := range model {
		projects[i] = ProjectFirestore{
			Title:       v.Title,
			Description: v.Description,
			Points:      v.Points,
			Technology:  v.Technology,
		}
	}
	return projects
}

func toProjects(model []ProjectFirestore) []models.Project {
	projects := make([]models.Project, len(model))
	for i, v := range model {
		projects[i] = models.Project{
			Title:       v.Title,
			Description: v.Description,
			Points:      v.Points,
			Technology:  v.Technology,
		}
	}
	return projects
}

func toExperience(doc *ExperienceFirestoreDoc) *models.Experience {
	companies := make([]models.Company, len(doc.Content.Companies))
	for i, v := range doc.Content.Companies {
		companies[i] = models.Company{
			Name:        v.Name,
			Description: v.Description,
			Period:      v.Period,
			Projects:    toProjects(v.Projects),
		}
	}
	return &models.Experience{
		Name:      doc.Content.Name,
		Title:     doc.Content.Title,
		Companies: companies,
	}
}

func NewExperienceRepository(firestore *internal.FirestoreClient) ExperienceRepository {
	return &experienceRepositoryImpl{
		firestore: firestore,
	}
}

type ProjectFirestore struct {
	Title       string   `firestore:"title"`
	Description string   `firestore:"description"`
	Points      []string `firestore:"points"`
	Technology  []string `firestore:"technology"`
}

type CompanyFirestore struct {
	Name        string             `firestore:"name"`
	Period      string             `firestore:"period"`
	Description string             `firestore:"description"`
	Projects    []ProjectFirestore `firestore:"projects"`
}

type ExperienceFirestore struct {
	Name      string             `firestore:"name"`
	Title     string             `firestore:"title"`
	Companies []CompanyFirestore `firestore:"companies"`
}

type ExperienceFirestoreDoc struct {
	LastUpdate int64               `firestore:"lastUpdate"`
	Content    ExperienceFirestore `firestore:"content"`
}
