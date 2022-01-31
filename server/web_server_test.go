package main

import (
	"fmt"
	"testing"
)

func TestGetCookie(t *testing.T) {
	// Init

	// Execution

	// Validation

}

func TestWebServer_Port(t *testing.T) {
	// Init
	ws := &WebServer{port: 8000}
	p1 := uint16(8000)

	// Execution
	p2 := ws.Port()

	// Validation
	if !(p1 == p2) {
		t.Error("error")
	}
}

func TestJsonStatus(t *testing.T) {
	// Init

	// Execution
	result := string(JsonStatus("success"))

	// Validation
	fmt.Println(result)

}

func TestWebServer_ShowCookieValue(t *testing.T) {
	//Init

	//

}
