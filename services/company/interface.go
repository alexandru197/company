package services

import (
	"context"

	"github.com/alexandru197/company/model"
)

type CompanyService interface {
	GetCompanyById(ctx context.Context, id model.ID) (model.Company, error)
	CreateCompany(ctx context.Context, company model.Company) (model.Company, error)
	DeleteCompany(ctx context.Context, id string) error
	PatchCompany(ctx context.Context, company model.Company, updates map[string]interface{}) (model.Company, error)
}
