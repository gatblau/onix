package docs

// swagger:route POST /db/run/{taskName}/{appVersion} DATABASE RUN
//
// RUN
//
// Executes the database task defined in a remote git repository.
//
// Configuration is provided as environment variables to the running service.
//
// responses:
//        200: dbRunResponse

// The output of the run command execution as a plain text log.
// swagger:response dbRunResponse
type dbRunResponse struct {
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

// The name of the task to run.
//
// The command output as a plain text log.
// swagger:parameters DATABASE RUN
type TaskName struct {
	// The application version to deploy
	//
	// in: path
	// required: true
	TaskName string `json:"taskName"`
}
