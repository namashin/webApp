package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type WebServer struct {
	port uint16
}

func NewWebServer(port uint16) *WebServer {
	return &WebServer{port: port}
}

func (ws *WebServer) Port() uint16 {
	return ws.port
}

// GetCookie just return cookie value
func GetCookie(w http.ResponseWriter, req *http.Request) *http.Cookie {
	c, err := req.Cookie("session")
	if err != nil {
		// ランダムにIDを生成してくれる。
		sID := uuid.New()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
	}
	return c
}

func (ws *WebServer) ShowCookieValue(w http.ResponseWriter, req *http.Request) {
	c := GetCookie(w, req)

	switch req.Method {
	case http.MethodGet:
		err := Save("cookie.txt", c.Value)
		if err != nil {
			log.Println("err in save method")
		}
		_, err = io.WriteString(w, string(JsonStatus("success")))
		if err != nil {
			log.Println("error in io.WriteString", err)
		}

		return
	default:
		_, err := io.WriteString(w, string(JsonStatus("fail")))
		if err != nil {
			log.Println("error in io.WriteString", err)
		}
	}
}

func Save(fileName string, value string) error {
	return ioutil.WriteFile(fileName, []byte(value), 0664)
}

func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}

func (ws *WebServer) index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join("templates", "index.html"))
		if err != nil {
			log.Println("err in parse files", err)
			return
		}
		err = t.Execute(w, "")

		if err != nil {
			log.Println("ERROR in execute error", err)
			return
		}
	case http.MethodPost:
		mf, fh, err := req.FormFile("inputFile")
		if err != nil {
			log.Println("error in req.FormFile", err)
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Println("error in os.GetWd()")
			return
		}

		fPath := filepath.Join(wd, "server", "pics", fh.Filename)
		pFile, err := os.Create(fPath)
		if err != nil {
			log.Println("error in os.Create", err)
			return
		}
		defer pFile.Close()

		_, err = mf.Seek(0, 0)
		if err != nil {
			log.Println("error in mf.Seek", err)
			return
		}
		_, err = io.Copy(pFile, mf)
		if err != nil {
			log.Println("error in io.Copy *FILE", err)
			return
		}
		log.Println(string(JsonStatus("success")))

		// ファイル送信成功したら、元の/indexにリダイレクト
		http.Redirect(w, req, "/index", 301)
	default:
		log.Println("ERROR: invalid request method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WebServer) DisplayPort(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		_, err := io.WriteString(w, req.Method+"  and  "+strconv.Itoa(int(ws.Port())))
		if err != nil {
			log.Println("error in io.WriteString", err)
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WebServer) DisplaySavedPic(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./templates/display_pic.html")
	if err != nil {
		log.Println("error in parse display_pic.html file", err)
		return
	}

	// TODO
	// picsに保存した写真をhtmlで表示
	fNamesInPics, err := ioutil.ReadDir("./server/pics/")
	if err != nil {
		log.Println("error in ReadDir", err)
		return
	}
	var Names []string

	for _, fNameInPics := range fNamesInPics {
		Name := fNameInPics.Name()
		Names = append(Names, Name)

	}

	err = t.Execute(w, Names)
	if err != nil {
		log.Println("template execute error", err)
	}
}

func DirWalk(dir string) []byte {
	elems, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println("error in ioutil.ReadFile", err)
	}
	return elems
}

func (ws *WebServer) Run() {
	http.HandleFunc("/", ws.index)
	http.HandleFunc("/cookie", ws.ShowCookieValue)
	http.HandleFunc("/port", ws.DisplayPort)
	http.HandleFunc("/pics", ws.DisplaySavedPic)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
