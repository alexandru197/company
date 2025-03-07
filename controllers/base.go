package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexandru197/company/model"
	services "github.com/alexandru197/company/services/company"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BaseController struct {
	companyService services.CompanyService
}

func NewBaseController(companyService services.CompanyService) *BaseController {
	return &BaseController{
		companyService: companyService,
	}
}

// GetCompany - Returns company based on ID
// @Summary This API can be used to retrieve a Company instance from the DB.
// @Description Retrieves a Company instance from the DB.
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router / [get]
func (c *BaseController) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["id"])
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var comp model.Company
	comp, err = c.companyService.GetCompanyById(r.Context(), model.ID(companyID.String()))
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(comp)
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, resp, http.StatusOK)
}

// AddCompany - Adds a new Company
// @Summary This API can be used to add a new Company instance to the DB.
// @Description Adds a new Company instance to the DB.
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router / [create]
func (c *BaseController) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var comp model.Company
	if err := json.NewDecoder(r.Body).Decode(&comp); err != nil {
		ApplicationErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if comp.ID == uuid.Nil {
		ApplicationErrorResponse(w, errors.New("company ID is required"), http.StatusBadRequest)
		return
	}
	if err := model.ValidateCompany(&comp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comp, err := c.companyService.CreateCompany(r.Context(), comp)
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(comp)
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, resp, http.StatusCreated)
}

// PatchCompany - Updates company based on ID
// @Summary This API can be used to update a Company instance from the DB.
// @Description Updates a Company instance from the DB.
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router / [patch]
func (c *BaseController) PatchCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["id"])
	if err != nil {
		ApplicationErrorResponse(w, errors.New("invalid company id"), http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		ApplicationErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	var comp model.Company
	comp, err = c.companyService.GetCompanyById(r.Context(), model.ID(companyID.String()))
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	comp, err = c.companyService.PatchCompany(r.Context(), comp, updates)
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(comp)
	if err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, resp, http.StatusOK)
}

// GetCompany - Deletes company based on ID
// @Summary This API can be used to delete a Company instance from the DB.
// @Description Deletes a Company instance from the DB.
// @Accept  json
// @Produce  json
// @Success 200 {string} response "api response"
// @Router / [delete]
func (c *BaseController) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["id"])
	if err != nil {
		ApplicationErrorResponse(w, errors.New("invalid company ID"), http.StatusBadRequest)
		return
	}

	if err = c.companyService.DeleteCompany(r.Context(), companyID.String()); err != nil {
		ApplicationErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, nil, http.StatusNoContent)
}
