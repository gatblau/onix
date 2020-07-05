package docs

// swagger:route GET /db/server DATABASE SERVER
//
// SERVER
//
// Gets database server information.
//
// responses:
//        200: dbServerResponse

// The output of the server command execution.
// swagger:response dbServerResponse
type dbServerResponse struct {
	// in:body
	Body string
}

// // The application version tag.
// //
// // The command output as a plain text log.
// // swagger:parameters database deploy
// type AppVersion struct {
// 	// The application version to deploy
// 	//
// 	// in: path
// 	// required: true
// 	AppVersion string `json:"appVersion"`
// }
//
// // The name of the task to run.
// //
// // The command output as a plain text log.
// // swagger:parameters DATABASE RUN
// type TaskName struct {
// 	// The application version to deploy
// 	//
// 	// in: path
// 	// required: true
// 	TaskName string `json:"taskName"`
// }
