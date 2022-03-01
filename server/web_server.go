package server

import (
	"bufio"
	"encoding/json"
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

	pFile, err := os.Open("cookie.txt")
	if err != nil {
		log.Fatal("error in open file =>", err)
	}

	scan := bufio.NewScanner(pFile)

	for i := 0; scan.Scan(); i++ {
		line := scan.Text()

		if c.Value == line {
			visitedCount++
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
	default:
		log.Println("ERROR: invalid request method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WebServer) BackToIndex(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join("templates", "back.html"))
	if err != nil {
		log.Println("ERROR in parse files", err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println("ERROR in Execute template", err)
		return
	}
}

func (ws *WebServer) UploadFileHandler(w http.ResponseWriter, req *http.Request) {
	mf, fh, err := req.FormFile("inputFile")
	fileExplanation := req.FormValue("fileExplanation")

	if err != nil {
		log.Println("error in req.FormFile: => ", err)
		_, err = io.WriteString(w, "ファイルを選択してください")
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	err = Save("fileExplanation.txt", fileExplanation)
	if err != nil {
		log.Println("fileExplanation.txt does not saved", err)
		return
	}

	//val, ok := checkPullDownMenu(w, req)
	//if !ok {
	//	log.Println("invalid select form")
	//	return
	//}
	//fmt.Println("選ばれたのは" + val + "でした")

	baseDir := GetBaseDirectory()

	fPath := filepath.Join(baseDir, "pics", fh.Filename)
	pFile, err := os.Create(fPath)
	if err != nil {
		log.Println("error in os.Create", err)
		return
	}

	defer pFile.Close()

	v, err := mf.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("error in mf.Seek", err, "walked", v)
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
}

func (ws *WebServer) index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		t, err := template.ParseFiles(path.Join("templates", "index.html"))
		if err != nil {
			log.Println("err in parse files", err)
			return
		}
		token := makeToken()

		err = t.Execute(w, token)
		if err != nil {
			log.Println("ERROR: in execute error", err)
			return
		}
	case http.MethodPost:
		ws.UploadFileHandler(w, req)
	default:
		log.Println("ERROR: invalid request method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func makeFile() {

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
// 写真が表示されない なぜ??
func (ws *WebServer) DisplaySavedPic(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./templates/display_pic.html")
	if err != nil {
		log.Println("error in parse display_pic.html file =>>>", err)
		return
	}

	fNamesInPicsDir, err := ioutil.ReadDir("./pics")
	if err != nil {
		log.Println("error in ReadDir", err)
		return
	}
	var NamesStrList []string

	for _, fNameInPicsDir := range fNamesInPicsDir {
		Name := fNameInPicsDir.Name()
		NamesStrList = append(NamesStrList, Name)
	}

	err = t.Execute(w, NamesStrList)
	if err != nil {
		log.Println("template execute error", err)
	}
}

//func Handler(w http.ResponseWriter, req *http.Request) {
//	io.WriteString(w, "test")
//}

func (ws *WebServer) Run() {
	http.HandleFunc("/", ws.index)
	//	http.HandleFunc("/test", Handler)
	http.HandleFunc("/del", ws.ExpireCookie)
	http.HandleFunc("/del/back", ws.BackToIndex)
	http.HandleFunc("/cookie", ws.ShowCookieValue)
	http.HandleFunc("/port", ws.DisplayPort)
	http.HandleFunc("/pics", ws.DisplaySavedPic)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
