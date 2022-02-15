package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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

func TestHandler(t *testing.T) {
	//Init
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	//Execution
	Handler(w, req)

	resp := w.Result()
	fmt.Println(resp)

	body, _ := ioutil.ReadAll(resp.Body)

	// Validation
	if string(body) != "test" {
		t.Errorf("error")
	}
}

func TestGetCookie(t *testing.T) {
	// Init
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	w := httptest.NewRecorder()

	// execution
	got := GetCookie(w, req)

	// validation
	if !(len(got.Value) > 0) {
		t.Error("not generated cookie")
	}
}

func TestGetBaseDirectory(t *testing.T) {
	// Init
	myBaseDir := "C:\\Users\\Nagahara\\GolandProjects\\webApp\\server"
	fmt.Println(myBaseDir)

	// Execution
	got := GetBaseDirectory()
	fmt.Println(got)

	// validation
	if myBaseDir != got {
		t.Error("dir not match")
	}
}

func TestWebServer_DisplaySavedPic(t *testing.T) {
	//Init
	ws := &WebServer{port: 8000, gateway: "http://127.0.0.1"}

	ts := httptest.NewServer(http.HandlerFunc(ws.DisplaySavedPic))
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Error("http.Get error", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Error("status error")
	}
}

func TestWebServer_ShowCookieValue(t *testing.T) {
	ws := NewWebServer(8000, "http://127.0.0.1")

	req, _ := http.NewRequest("GET", "/cookie", nil)
	w := httptest.NewRecorder()

	c := GetCookie(w, req)
	ws.ShowCookieValue(w, req)

	if !(len(c.Value) > 0) {
		t.Error("cookie not made ")
	}

	f, err := os.Stat("cookie.txt")
	if os.IsNotExist(err) {
		t.Error("ファイルが生成されてません", err)
	}

	fmt.Println(f.Size())
	fmt.Println(f.Name())

}
