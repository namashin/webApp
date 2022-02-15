package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetCookie makes new cookie (if it is not made before)
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

// Save receive filename and value. And save its value in filename
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

// JsonStatus return status in JSON style
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

func checkPullDownMenu(w http.ResponseWriter, req *http.Request) (string, bool) {
	slice := []string{"apple", "pear", "banana"}
	val := req.Form.Get("fruit")

	for _, v := range slice {
		if v == val {
			return val, true
		}
	}
	return "", false
}

func GetBaseDirectory() string {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal("ベースとなるディレクトリを取得できてない可能性あり => ", err)
	}
	return baseDir
}

func makeToken() string {
	nowTime := time.Now().UnixNano()
	h := md5.New()
	_, err := io.WriteString(h, strconv.FormatInt(nowTime, 10))
	if err != nil {
		log.Fatal("io.WriteString error", err)
	}

	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("tokenの値 =>" + token)
	return token
}
