package docs

// swagger:route GET /conf CONFIGURATION SHOW
//
// SHOW
//
// List all configuration variables.
//
// responses:
//        200: confResponse

// The list of configuration variables used by DbMan.
// swagger:response confResponse
type confResponse struct {
	// in:body
	Body string
}
