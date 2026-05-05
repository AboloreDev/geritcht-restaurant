package services

import (
	"errors"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type AllergenService struct {
	db *gorm.DB
}

type DietaryTagsService struct {
	db *gorm.DB
}

func NewAllergenService(db *gorm.DB) *AllergenService {
	return &AllergenService{db: db}
}

func NewDietaryTagsService(db *gorm.DB) *DietaryTagsService {
	return &DietaryTagsService{db: db}
}

func (s *AllergenService) CreateAllergenServices(req *dto.CreateAllergenRequest) error {
	var allergen models.Allergen

	err := s.db.Where("name = ?", req.Name).First(&allergen).Error
	if err == nil {
		return domain.ErrNameConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	allergen = models.Allergen{
		Name: req.Name,
	}

	err = s.db.Create(&allergen).Error
	if err != nil {
		return err
	}

	return nil
}
func (s *DietaryTagsService) CreateDietaryTags(req *dto.CreateDietaryTagRequest) error {
	var tags models.DietaryTag

	err := s.db.Where("name = ?", req.Name).First(&tags).Error
	if err == nil {
		return domain.ErrNameConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tags = models.DietaryTag{
		Name: req.Name,
	}

	err = s.db.Create(&tags).Error
	if err != nil {
		return err
	}

	return nil
}
