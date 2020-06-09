package docs

// swagger:route POST /db/deploy/{appVersion} database deploy
//
// deploy
//
// Executes the database schema / functions deployment scripts defined in the remote git repository.
//
// Due to the fact some scripts might need to execute operations such as creating the database, they are not executed
// within a transaction block.
//
// Configuration is provided as environment variables to the running service.
//
// responses:
//        200: dbDeployResponse

// The output of the deploy command execution as a plain text log.
// swagger:response dbDeployResponse
type dbDeployResponse struct {
	// in:body
	Body string
}

// The application version tag.
//
// The command output as a plain text log.
// swagger:parameters database deploy
type AppVersion struct {
	// The application version to deploy
	//
	// in: path
	// required: true
	AppVersion string `json:"appVersion"`
}
