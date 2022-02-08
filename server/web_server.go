package main

import (
	"bufio"
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
	port    uint16
	gateway string
}

func NewWebServer(port uint16, gateway string) *WebServer {
	return &WebServer{port: port, gateway: gateway}
}

func (ws *WebServer) Port() uint16 {
	return ws.port
}

func (ws *WebServer) Gateway() string {
	return ws.gateway
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

func (ws *WebServer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		port    uint16
		gateway string
	}{
		port:    ws.port,
		gateway: ws.gateway,
	})
}

func (ws *WebServer) ExpireCookie(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		log.Println(string(JsonStatus("fail")))
		http.Redirect(w, req, "/index", http.StatusSeeOther)
	}
	c.MaxAge = -1
	http.SetCookie(w, c)
	log.Println("cookieを削除しました。")
	http.Redirect(w, req, "/del/back", http.StatusSeeOther)
}

func (ws *WebServer) checkVisitedCount(c *http.Cookie) int {

	var visitedCount int

	b, err := os.Open("cookie.txt")
	if err != nil {
		log.Fatal("error in open file =>", err)
	}

	scan := bufio.NewScanner(b)

	for i := 0; scan.Scan(); i++ {
		line := scan.Text()

		if c.Value == line {
			visitedCount++
		} else {
			continue
		}
	}
	return visitedCount
}

func (ws *WebServer) ShowCookieValue(w http.ResponseWriter, req *http.Request) {
	c := GetCookie(w, req)

	switch req.Method {
	case http.MethodGet:
		err := Save("cookie.txt", c.Value)
		if err != nil {
			log.Println("err in save method =>", err)
		}

		VisitedCount := ws.checkVisitedCount(c)

		_, err = io.WriteString(w, "あなたは今日"+strconv.Itoa(VisitedCount)+"回このサイトに訪れました")
		if err != nil {
			log.Println("error in io.WriteString", err)
		}
		return
	default:
		log.Println("ERROR: invalid request method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Save(fileName string, value string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error in os.OpenFile =>>>", err)
	}
	defer f.Close()

	_, err = f.WriteString(value + "\n")
	if err != nil {
		log.Fatal("Error in f.WriteString =>>>", err)
	}
	return nil
}

func JsonStatus(message string) []byte {
	m, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	if err != nil {
		log.Println("json.Marshal error ", err)
	}

	return m
}

func (ws *WebServer) BackToIndex(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join("templates", "back.html"))
	if err != nil {
		log.Println("ERROR in parse files", err)
		return
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Println("ERROR in Execute template", err)
		return
	}
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
			_, err = io.WriteString(w, "ファイルを選択してください")
			if err != nil {
				log.Println(err)
				return
			}
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

		// ファイル送信成功したら、元の/indexにリダイレクトさせる
		http.Redirect(w, req, "/index", http.StatusFound)
	default:
		log.Println("ERROR: invalid request method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WebServer) DisplayPort(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		_, err := io.WriteString(w, "あなたのポート番号は"+strconv.Itoa(int(ws.Port()))+"番です")
		if err != nil {
			log.Println("error in io.WriteString", err)
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// DisplaySavedPic
// 写真が表示されない
func (ws *WebServer) DisplaySavedPic(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./templates/display_pic.html")
	if err != nil {
		log.Println("error in parse display_pic.html file =>>>", err)
		return
	}

	fNamesInPicsDir, err := ioutil.ReadDir("./server/pics/")
	if err != nil {
		log.Println("error in ReadDir", err)
		return
	}
	var Names []string

	for _, fNameInPicsDir := range fNamesInPicsDir {
		Name := fNameInPicsDir.Name()
		Names = append(Names, Name)
	}

	err = t.Execute(w, Names)
	if err != nil {
		log.Println("template execute error", err)
	}
}

func (ws *WebServer) Run() {
	http.HandleFunc("/", ws.index)
	http.HandleFunc("/del", ws.ExpireCookie)
	http.HandleFunc("/del/back", ws.BackToIndex)
	http.HandleFunc("/cookie", ws.ShowCookieValue)
	http.HandleFunc("/port", ws.DisplayPort)
	http.HandleFunc("/pics", ws.DisplaySavedPic)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
