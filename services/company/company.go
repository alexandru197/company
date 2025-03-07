package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/alexandru197/company/model"
	repository "github.com/alexandru197/company/repository/company"
	"github.com/segmentio/kafka-go"
)

const (
	COMPANY_CREATED_SUCCESS = "COMPANY_CREATED_SUCCESS"
	COMPANY_CREATED_FAILED  = "COMPANY_CREATED_FAILED"

	COMPANY_DELETED_SUCCESS = "COMPANY_DELETED_SUCCESS"
	COMPANY_DELETED_FAILED  = "COMPANY_DELETED_FAILED"

	COMPANY_PATCHED_SUCCESS = "COMPANY_PATCHED_SUCCESS"
	COMPANY_PATCHED_FAILED  = "COMPANY_PATCHED_FAILED"
)

type companyService struct {
	CompanyRepository repository.CompanyRepository
	KafkaWriter       *kafka.Writer
}

func NewCompanyService(cr repository.CompanyRepository, kr *kafka.Writer) CompanyService {
	return &companyService{
		CompanyRepository: cr,
		KafkaWriter:       kr,
	}
}

func (cs *companyService) GetCompanyById(ctx context.Context, id model.ID) (model.Company, error) {
	company, err := cs.CompanyRepository.FindById(ctx, string(id))
	if err != nil {
		log.Printf("error retrieving Company with ID %v: %v", id, err)
		return model.Company{}, err
	}
	return company, nil
}

func (cs *companyService) CreateCompany(ctx context.Context, company model.Company) (model.Company, error) {
	company, err := cs.CompanyRepository.CreateCompany(ctx, company)
	if err != nil {
		log.Printf("error creating Company: %v", err)
		publishCompanyEvent(COMPANY_CREATED_FAILED, cs.KafkaWriter, company)
		return model.Company{}, err
	}

	publishCompanyEvent(COMPANY_CREATED_SUCCESS, cs.KafkaWriter, company)
	return company, nil
}

func (cs *companyService) DeleteCompany(ctx context.Context, id string) error {
	company, err := cs.CompanyRepository.FindById(ctx, id)
	if err != nil {
		log.Printf("error retrieving company with ID %v: %v", id, err)
		return err
	}

	err = cs.CompanyRepository.DeleteCompany(ctx, company)
	if err != nil {
		log.Printf("error deleting Company with id %v: %v", id, err)
		publishCompanyEvent(COMPANY_DELETED_FAILED, cs.KafkaWriter, company)
		return err
	}

	publishCompanyEvent(COMPANY_DELETED_SUCCESS, cs.KafkaWriter, company)
	return nil
}

func (cs *companyService) PatchCompany(ctx context.Context, company model.Company, updates map[string]interface{}) (model.Company, error) {
	company, err := cs.CompanyRepository.PatchCompany(ctx, company, updates)
	if err != nil {
		log.Printf("error patching company with id %v", company.ID)
		publishCompanyEvent(COMPANY_PATCHED_FAILED, cs.KafkaWriter, company)
		return model.Company{}, err
	}

	publishCompanyEvent(COMPANY_PATCHED_SUCCESS, cs.KafkaWriter, company)
	return company, nil
}

func publishCompanyEvent(eventType string, kafkaWriter *kafka.Writer, company model.Company) error {
	event := map[string]interface{}{
		"event":   eventType,
		"company": company,
		"time":    time.Now().Format(time.RFC3339),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return err
	}
	msg := kafka.Message{
		Key:   []byte(company.ID.String()),
		Value: payload,
	}
	if err = kafkaWriter.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("Error publishing event: %v", err)
		return err
	}

	return nil
}
