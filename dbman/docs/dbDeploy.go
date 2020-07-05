package docs

// swagger:route POST /db/deploy DATABASE DEPLOY
//
// DEPLOY
//
// Deploys the database schemas and objects for the current database release.
//
// responses:
//        200: dbDeployResponse

// The output of the deploy command execution.
// swagger:response dbDeployResponse
type dbDeployResponse struct {
	// in:body
	Body string
}
