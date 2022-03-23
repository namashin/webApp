package server

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
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
		t.Fail()
	}
}

//func TestHandler(t *testing.T) {
//	//Init
//	req := httptest.NewRequest(http.MethodGet, "/test", nil)
//	w := httptest.NewRecorder()
//
//	//Execution
//	Handler(w, req)
//
//	resp := w.Result()
//	fmt.Println(resp)
//
//	body, _ := ioutil.ReadAll(resp.Body)
//
//	// Validation
//	if string(body) != "test" {
//		t.Errorf("error")
//	}
//}

func TestGetCookie(t *testing.T) {
	// Init
	req, err := http.NewRequest(http.MethodGet, "/cookie", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	// execution
	got := GetCookie(w, req)

	// validation
	print(got.Value)
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

const IMAGE_PATH = "../pics/R.jpg"

func TestWebServer_UploadFileHandler(t *testing.T) {
	//Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	//This writers is going to transform
	//what we pass to it to multipart form data
	//and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	ws := NewWebServer(8000, "http://127.0.0.1")

	go func() {
		defer func() {
			err := writer.Close()
			if err != nil {
				t.Error("close pipe writer ERROR")
			}
		}()
		//we create the form data field 'fileupload'
		//which returns another writer to write the actual file

		part, err := writer.CreateFormFile("inputFile", "R.png")
		if err != nil {
			t.Error(err)
		}

		//https://yourbasic.org/golang/create-image/
		img := image.NewRGBA(image.Rect(0, 0, 12, 6))

		//Encode() takes an io.Writer.
		//We pass the multipart field
		//'fileupload' that we defined
		//earlier which, in turn, writes
		//to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	//We read from the pipe which receives data
	//from the multipart writer, which, in turn,
	//receives data from png.Encode().
	//We have 3 chained writers !
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.UploadFileHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 302 {
		t.Errorf("Expected %v, received %d", 302, response.Code)
	}

	f, err := os.Stat("../pics/R.jpg")
	if f.Size() < 0 {
		t.Error("not exist")
	}
	if err != nil {
		t.Error("os.Stat ERROR")
	}
}

func TestJsonStatus(t *testing.T) {
	// Init
	statusMessage := "success"

	// execution
	ret := JsonStatus(statusMessage)

	// validation
	if string(ret) != `{"message":"success"}` {
		t.Errorf("expected %s  but got %s", `{"message":"success"}`, ret)
	}
}
