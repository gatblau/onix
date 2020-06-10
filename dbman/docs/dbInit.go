package docs

// swagger:route POST /db/init DATABASE INIT
//
// INITIALISE
//
// Executes the database initialisation scripts defined in the remote git repository.
//
// Due to the fact some scripts might need to execute operations such as creating the database, they are not executed
// within a transaction block.
//
// Configuration is provided as environment variables to the running service.
//
// responses:
//        200: dbInitResponse

// The output of the initialisation command execution as a plain text log.
// swagger:response dbInitResponse
type dbInitResponse struct {
	// in:body
	Body string
}
