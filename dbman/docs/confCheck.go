package docs

// swagger:route GET /conf/check CONFIGURATION CHECK
//
// CHECK
//
// Checks that a connection can be established to the database and remote git repository based on the specified configuration.
//
// Configuration is provided as environment variables to the running service.
//
// responses:
//        200: confCheckResponse

// A list of checks and their result displayed in text form.
// swagger:response confCheckResponse
type confCheckResponse struct {
	// in:body
	Body string
}
