package main

import (
	"errors"
	"github.com/Azure/go-ntlmssp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var SERVER_PORT = getVar(os.Getenv("SERVER_PORT"), "8080")
var DEFAULT_NTLM_URL = os.Getenv("DEFAULT_NTLM_URL")
var DEFAULT_NTLM_USERNAME = os.Getenv("DEFAULT_NTLM_USERNAME")
var DEFAULT_NTLM_PASSWORD = os.Getenv("DEFAULT_NTLM_PASSWORD")
var BASICAUTH_USERNAME = os.Getenv("BASICAUTH_USERNAME")
var BASICAUTH_PASSWORD = os.Getenv("BASICAUTH_PASSWORD")

func auth(r *http.Request) error {
	if len(BASICAUTH_USERNAME) > 0 && len(BASICAUTH_PASSWORD) > 0 {
		clientUsername, clientPassword, ok := r.BasicAuth()
		if !ok {
			log.Println("basic auth required")
			return errors.New("basic auth required")
		}
		if clientUsername != BASICAUTH_USERNAME || clientPassword != BASICAUTH_PASSWORD {
			log.Println("invalid username or password")
			return errors.New("invalid username or password")
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// auth the client request
	if err := auth(r); err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	_NTLM_URL := getVar(r.Header.Get("NTLM_URL"), DEFAULT_NTLM_URL)
	_NTLM_USERNAME := getVar(r.Header.Get("NTLM_USERNAME"), DEFAULT_NTLM_USERNAME)
	_NTLM_PASSWORD := getVar(r.Header.Get("NTLM_PASSWORD"), DEFAULT_NTLM_PASSWORD)

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper:&http.Transport{},
		},
	}

	url := _NTLM_URL + r.URL.RequestURI()
	req, _ := http.NewRequest(r.Method, _NTLM_URL + "/" + r.URL.RequestURI(), r.Body)
	req.SetBasicAuth(_NTLM_USERNAME, _NTLM_PASSWORD)
	log.Println("sending request to", url)
	res, err := client.Do(req)
	if err != nil {
		log.Println("erroneous call", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cnt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// copy NTLM response headers to client response
	for k, v := range res.Header {
		for i := range v {
			w.Header().Add(k, v[i])
		}
	}

	// copy same status code
	w.WriteHeader(res.StatusCode)

	// copy the NTLM response body to client response
	w.Write(cnt)
}

func getVar(newValue string, defaultValue string) string {
	if len(newValue) == 0 {
		return defaultValue
	}
	return newValue
}

func main() {
	log.Println("starting service...")

	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(":" + SERVER_PORT, nil))
}
