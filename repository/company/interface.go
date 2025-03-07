package repository

import (
	"context"

	"github.com/alexandru197/company/model"
)

type CompanyRepository interface {
	FindById(ctx context.Context, id string) (model.Company, error)
	CreateCompany(ctx context.Context, company model.Company) (model.Company, error)
	DeleteCompany(ctx context.Context, company model.Company) error
	PatchCompany(ctx context.Context, company model.Company, updates map[string]interface{}) (model.Company, error)
}
