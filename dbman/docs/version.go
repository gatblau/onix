package docs

// swagger:route GET /version DATABASE VERSION
//
// VERSION
//
// List all releases performed on the target database.
//
// responses:
//        200: versionResponse

// The history of all database releases
// swagger:response versionResponse
type versionResponse struct {
	// in:body
	Body string
}
