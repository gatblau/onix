package docs

// swagger:route POST /db/create DATABASE CREATE
//
// CREATE
//
// Creates a new database based on the release manifest.
//
// responses:
//        200: dbCreateResponse

// The output of the create command execution.
// swagger:response dbCreateResponse
type dbCreateResponse struct {
	// in:body
	Body string
}
