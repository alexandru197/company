package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/alexandru197/company/model"
	"gorm.io/gorm"
)

type companyRepository struct {
	DB *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		DB: db,
	}
}

func (cr *companyRepository) FindById(ctx context.Context, id string) (model.Company, error) {
	var comp model.Company
	if err := cr.DB.First(&comp, "id = ?", id).Error; err != nil {
		log.Printf("error retrieving Company with ID %v: %v", id, err)
		return comp, err
	}
	return comp, nil
}

func (cr *companyRepository) CreateCompany(ctx context.Context, company model.Company) (model.Company, error) {
	var existing model.Company
	if err := cr.DB.Where("LOWER(name) = ?", strings.ToLower(company.Name)).First(&existing).Error; err == nil {
		log.Printf("Company name must be unique")
		return company, err
	}
	if err := cr.DB.Create(&company).Error; err != nil {
		log.Printf("Error creating company")
		return company, err
	}
	return company, nil
}

func (cr *companyRepository) DeleteCompany(ctx context.Context, comp model.Company) error {
	if err := cr.DB.Delete(&comp).Error; err != nil {
		log.Printf("Error deleting company with ID %v", comp.ID)
		return err
	}
	return nil
}

func (cr *companyRepository) PatchCompany(ctx context.Context, comp model.Company, updates map[string]interface{}) (model.Company, error) {
	if name, ok := updates["name"].(string); ok {
		if len(name) == 0 || len(name) > 15 {
			err := errors.New("Name is required and must be 15 characters or less")
			return model.Company{}, err
		}
		var existing model.Company
		if err := cr.DB.Where("LOWER(name) = ? AND id <> ?", strings.ToLower(name), comp.ID).First(&existing).Error; err == nil {
			err := errors.New("Company name must be unique")
			return model.Company{}, err
		}
		updates["name"] = name
	}
	if desc, ok := updates["description"].(string); ok {
		if len(desc) > 3000 {
			err := errors.New("Description cannot exceed 3000 characters")
			return model.Company{}, err
		}
	}
	if typ, ok := updates["type"].(string); ok {
		if !model.AllowedTypes[typ] {
			err := errors.New("Invalid company type")
			return model.Company{}, err
		}
	}
	if err := cr.DB.Model(&comp).Updates(updates).Error; err != nil {
		log.Printf("Error updating company")
		err := errors.New("Error updating company")
		return model.Company{}, err
	}

	cr.DB.First(&comp, "id = ?", comp.ID)

	return comp, nil
}
