package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"testing"
)

// launch the server
func TestServer_Serve(t *testing.T) {
	// create a new server
	s := New("test")
	// set auth credentials
	err := os.Setenv("OX_HTTP_UNAME", "test.user")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("OX_HTTP_PWD", "test-pwd")
	if err != nil {
		t.Fatal(err)
	}
	s.Http = func(router *mux.Router) {
		router.HandleFunc("/", doSomething).Methods("GET")
		router.Use(s.AuthenticationMiddleware)
	}
	// serve
	s.Serve()
}

func doSomething(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
