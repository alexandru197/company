package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Company represents a company entity.
type Company struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"size:15;unique;not null"`
	Description string    `json:"description" gorm:"size:3000"`
	Employees   int       `json:"employees" gorm:"not null"`
	Registered  bool      `json:"registered" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AllowedTypes defines valid company types.
var AllowedTypes = map[string]bool{
	"Corporations":        true,
	"NonProfit":           true,
	"Cooperative":         true,
	"Sole Proprietorship": true,
}

// validateCompany ensures the company meets the required rules.
func ValidateCompany(comp *Company) error {
	if len(comp.Name) == 0 || len(comp.Name) > 15 {
		return errors.New("name is required and must be 15 characters or less")
	}
	if len(comp.Description) > 3000 {
		return errors.New("description cannot exceed 3000 characters")
	}
	if !AllowedTypes[comp.Type] {
		return errors.New("invalid company type")
	}
	return nil
}
