package docs

// swagger:route POST /db/init db-resource db-init
// Initialises the database using the configuration in the current configuration set.
// responses:
//   200:

// This text will appear as description of your response body.
// swagger:response string
type dbInitResponseWrapper struct {
	// in:body
	Body string
}
