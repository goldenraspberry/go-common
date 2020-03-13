package server

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	m := th(1)
	RunEngine(&m)
}

type th int

func (*th) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("receiver request")
}
