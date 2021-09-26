# Onix HTTP client

A go client for the Onix Web API.

## Usage 

An example of how to use is below:

```go
// import the library
package main 

import "github.com/gatblau/oxc"

func main() {
    // prepares the client configuration
    cfg := &oxc.ClientConf{
         BaseURI:            "http://localhost:8080",
         InsecureSkipVerify: true,
         AuthMode:           oxc.Basic,
         Username:           "admin",
         Password:           "0n1x",
         // uncomment below & reset configuration vars
         // to test using an OAuth bearer token
         // AuthMode:           	OIDC,
         // TokenURI:     		"https://dev-447786.okta.com/oauth2/default/v1/token",
         // ClientId:			"0oalyh...356",
         // AppSecret:			"Tsed........OP0oEf9H7",
    }
    // create an instance of the web api client
    client, err := oxc.NewClient(cfg)

    if err != nil {
       panic(err)
    }

    // create a new model
    model := &oxc.Model {
        Key:         "test_model",
        Name:        "Test Model",
        Description: "Test Model",
    }

    // put the model
    result, err := client.PutModel(model)
    
    if err != nil {
       panic(err)
    }

    if result.Error {
        panic(result.Message)
    }
}

// create an instance of the client
func createClient() *oxc.Client {
    client, err := oxc.NewClient(&oxc.ClientConf{
        BaseURI:            "http://localhost:8080",
        InsecureSkipVerify: true,
        AuthMode:           oxc.Basic,
        Username:           "admin",
        Password:           "0n1x",
        // uncomment below & reset configuration vars
        // to test using an OAuth bearer token
        // AuthMode:           	OIDC,
        // TokenURI:     		"https://dev-447786.okta.com/oauth2/default/v1/token",
        // ClientId:			"0oalyh...356",
        // AppSecret:			"Tsed........OP0oEf9H7",
	})
	if err != nil { panic(err) }
	return client
}
```

More examples can be found [here](client_test.go).