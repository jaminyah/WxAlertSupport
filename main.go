/* Main driver function. 
 * Run: go run main.go fetch.go submit.go
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"encoding/json"

	_"github.com/mattn/go-sqlite3"
	"github.com/mojocn/base64Captcha"
)

var database, db_err = sql.Open("sqlite3", "./wxalert.db")


type UserComment struct {
	Id int				`json:"id"`
	Name string			`json:"name"`
	Message string		`json:"comment"`
	Timestamp string	`json:"date"`
}

//configJsonBody json request body.
type configJsonBody struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
}

// base64Captcha create http handler
func generateCaptchaHandler(w http.ResponseWriter, r *http.Request) {

	//parse request parameters
	decoder := json.NewDecoder(r.Body)
	var postParameters configJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	//create base64 encoding captcha
	var config interface{}
	switch postParameters.CaptchaType {
	case "audio":
		config = postParameters.ConfigAudio
	case "character":
		config = postParameters.ConfigCharacter
	default:
		config = postParameters.ConfigDigit
	}
	captchaId, captcaInterfaceInstance := base64Captcha.GenerateCaptcha(postParameters.Id, config)
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)

	//or you can just write the captcha content to the httpResponseWriter.
	//before you put the captchaId into the response COOKIE.
	//captcaInterfaceInstance.WriteTo(w)

	//set json response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]interface{}{"code": 1, "data": base64blob, "captchaId": captchaId, "msg": "success"}
	json.NewEncoder(w).Encode(body)
}

// base64Captcha verify http handler
func captchaVerifyHandle(w http.ResponseWriter, r *http.Request) {

	//parse request parameters
	decoder := json.NewDecoder(r.Body)
	var postParameters configJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	//verify the captcha
	verifyResult := base64Captcha.VerifyCaptcha(postParameters.Id, postParameters.VerifyValue)

	//set json response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]interface{}{"code": "error", "data": "validation failed", "msg": "captcha failed"}
	if verifyResult {
		body = map[string]interface{}{"code": "success", "data": "validation verified", "msg": "captcha verified"}
	}
	json.NewEncoder(w).Encode(body)
}

func main() {


	if db_err != nil {
		log.Fatal(db_err)
	}
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY, username TEXT, comment TEXT, date TEXT)")
	statement.Exec()	

	http.Handle("/", http.FileServer(http.Dir("static")))

    http.HandleFunc("/submit", submitAjax)

	http.HandleFunc("/comments", fetchComments)

	//api for creating captcha
	http.HandleFunc("/api/getCaptcha", generateCaptchaHandler)

	//api for verify captcha
	http.HandleFunc("/api/verifyCaptcha", captchaVerifyHandle)

	fmt.Printf("Starting HTTP server...\n")
    if server_err := http.ListenAndServe(":8088", nil); server_err != nil {
        log.Fatal(server_err)
    }
}