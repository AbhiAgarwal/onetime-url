package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"text/template"

	AuMail "github.com/aureum/aumail"
	"github.com/gorilla/mux"
)

var (
	sc  string
	au  AuMail.AuMail
	ids map[string]bool
)

type errorStruct struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error"`
}

type emailRequestStruct struct {
	SK    string `json:"sk"`
	Email string `json:"email"`
}

type emailResponseStruct struct {
	Status bool `json:"status"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateKey() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PrintMessage(res http.ResponseWriter, req *http.Request, message string, status int) {
	result := errorStruct{
		Status:       status,
		ErrorMessage: message,
	}
	js, _ := json.Marshal(result)
	res.WriteHeader(status)
	res.Write(js)
	return
}

func emailHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	var result emailRequestStruct
	err := decoder.Decode(&result)
	if err != nil {
		PrintMessage(res, req, "JSON parsing error.", http.StatusInternalServerError)
	}

	if result.SK == sc {
		generatedKey := generateKey()
		ids[generatedKey] = true

		au.From = "abhi@aureum.io"
		email := make([]string, 0)
		email = append(email, result.Email)
		au.Emails = email

		au.Subject = "A key for you!"
		au.Text = "Key: " + generatedKey
		status, _ := au.Send()
		var emailStatus emailResponseStruct
		emailStatus.Status = status

		js := []byte{}
		js, _ = json.Marshal(emailStatus)
		res.WriteHeader(200)
		res.Write(js)
	} else {
		PrintMessage(res, req, "Incorrect secret key.", http.StatusUnauthorized)
	}
	return
}

func homeHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	PrintMessage(res, req, "Unauthorized.", http.StatusUnauthorized)
}

func idHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if _, ok := ids[id]; ok {
		res.Header().Set("Content-Type", "text/html")
		t, _ := template.ParseFiles("public/index.html")
		t.Execute(res, nil)
		delete(ids, id)
		return
	} else {
		res.Header().Set("Content-Type", "application/json")
		PrintMessage(res, req, "Incorrect Key.", http.StatusForbidden)
	}
	return
}

func main() {
	ids = make(map[string]bool)

	// Setup
	au.SendGridUser = os.Getenv("SendGridUser")
	au.SendGridKey = os.Getenv("SendGridKey")
	sc = os.Getenv("OTUSecretKey")

	// Initialize router
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/{id:[0-9a-zA-Z]+}", idHandler).Methods("GET")
	r.HandleFunc("/email", emailHandler).Methods("POST")
	http.ListenAndServe(":3000", r)
}
