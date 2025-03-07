package adapters

import (
	"github.com/alexandru197/company/controllers"
	"github.com/gorilla/mux"
)

const (
	PATH_PREFIX          = "/companies"
	PATH_SUFFIX          = "/{id}"
	ROUTE_GET_COMPANY    = PATH_PREFIX + PATH_SUFFIX
	ROUTE_ADD_COMPANY    = PATH_PREFIX + "/createCompany"
	ROUTE_PATCH_COMPANY  = PATH_PREFIX + "/patchCompany" + PATH_SUFFIX
	ROUTE_DELETE_COMPANY = PATH_PREFIX + "/deleteCompany" + PATH_SUFFIX
)

// SetupRoutes configures the HTTP routes.
func (c *CompanyApp) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	baseController := controllers.NewBaseController(c.CompanyService)

	// Public endpoint.
	r.HandleFunc(ROUTE_GET_COMPANY, baseController.GetCompany).Methods("GET")

	// Protected endpoints.
	protected := r.PathPrefix(PATH_PREFIX).Subrouter()
	protected.Use(jwtMiddleware)
	protected.HandleFunc(ROUTE_ADD_COMPANY, baseController.CreateCompany).Methods("POST")
	protected.HandleFunc(ROUTE_PATCH_COMPANY, baseController.PatchCompany).Methods("PATCH")
	protected.HandleFunc(ROUTE_DELETE_COMPANY, baseController.DeleteCompany).Methods("DELETE")

	return r
}
