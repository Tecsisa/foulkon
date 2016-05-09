package http

import (
	"net/http"

	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
)

type PolicyHandler struct {
	core *authorizr.Core
}

// Requests
type CreatePolicyRequest struct {
	Name       string          `json:"Name, omitempty"`
	Path       string          `json:"Path, omitempty"`
	Statements []api.Statement `json:"Statements, omitempty"`
}

// Responses
type CreatePolicyResponse struct {
	Policy api.Policy
}

type ListPoliciesResponse struct {
	Policies []api.Policy
}

func (p *PolicyHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org from path
	org := ps.ByName(ORG_ID)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := p.core.PolicyApi.GetPolicies(org, pathPrefix)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondInternalServerError(w)
		return
	}

	// Create response
	response := &ListPoliciesResponse{
		Policies: result,
	}

	// Return data
	RespondOk(w, response)
}

func (p *PolicyHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Retrieve Organization
	org := ps.ByName(ORG_ID)

	// Decode request
	request := CreatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Validate policy
	policy, err := validatePolicy(createPolicy(request.Name, request.Path, org, &request.Statements))

	// Check errors
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Store this policy
	storedPolicy, err := p.core.PolicyApi.AddPolicy(policy)

	// Error handling
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondInternalServerError(w)
		return
	}

	response := &CreatePolicyResponse{
		Policy: *storedPolicy,
	}

	// Write group to response
	RespondOk(w, response)
}

func (p *PolicyHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleGetPolicyAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// This method validates policies created
func validatePolicy(policy api.Policy) (api.Policy, error) {
	// TODO rsoleto: Crear validador
	return policy, nil
}

// It returns a policy with its parameters setted according to method parameters
func createPolicy(name string, path string, org string, statements *[]api.Statement) api.Policy {
	// TODO rsoleto: Hay que validar la entrada acorde a una expresion regular
	// y quitar los elementos repetidos o no validos
	completePath := path + "/" + name
	urn := fmt.Sprintf("urn:iws:iam:%v:policy/%v", org, completePath)
	policy := api.Policy{
		ID:         uuid.NewV4().String(),
		Name:       name,
		Path:       completePath,
		Org:        org,
		Urn:        urn,
		Statements: statements,
	}

	return policy
}
