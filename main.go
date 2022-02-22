package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"github.com/kardianos/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger service.Logger

var SERVER_PORT = os.Getenv("SERVER_PORT")
var DEFAULT_NTLM_URL = os.Getenv("DEFAULT_NTLM_URL")
var DEFAULT_NTLM_USERNAME = os.Getenv("DEFAULT_NTLM_USERNAME")
var DEFAULT_NTLM_PASSWORD = os.Getenv("DEFAULT_NTLM_PASSWORD")
var BASICAUTH_USERNAME = os.Getenv("BASICAUTH_USERNAME")
var BASICAUTH_PASSWORD = os.Getenv("BASICAUTH_PASSWORD")

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Info("starting service...")

	http.HandleFunc("/", handler)
	logger.Error(http.ListenAndServe(":" + SERVER_PORT, nil))
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func auth(r *http.Request) error {
	if len(BASICAUTH_USERNAME) > 0 && len(BASICAUTH_PASSWORD) > 0 {
		clientUsername, clientPassword, ok := r.BasicAuth()
		if !ok {
			logger.Info("basic auth required")
			return errors.New("basic auth required")
		}
		if clientUsername != BASICAUTH_USERNAME || clientPassword != BASICAUTH_PASSWORD {
			logger.Info("invalid username or password")
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
	logger.Info("sending request to", url)
	res, err := client.Do(req)
	if err != nil {
		logger.Info("erroneous call", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cnt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Info("error reading response body", err)
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
	parseArgs()

	svcConfig := &service.Config{
		Name:        "NTLM-proxy",
		DisplayName: "NTLM-proxy",
		Description: "This is a NTLM-proxy.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func parseArgs() {
	var port string
	flag.StringVar(&port, "port", "8080", "Specify port. Default is 8080.")

	var ntlm_url string
	flag.StringVar(&ntlm_url, "ntlm_url", "", "Specify NTLM_URL.")

	var ntlm_username string
	flag.StringVar(&ntlm_username, "ntlm_username", "", "Specify NTLM_USERNAME.")

	var ntlm_password string
	flag.StringVar(&ntlm_password, "ntlm_password", "", "Specify NTLM_PASSWORD.")

	var basicauth_username string
	flag.StringVar(&basicauth_username, "basicauth_username", "", "Specify BASICAUTH_USERNAME.")

	var basicauth_password string
	flag.StringVar(&basicauth_password, "basicauth_password", "", "Specify BASICAUTH_PASSWORD.")

	flag.Usage = func() {
		fmt.Printf("Usage: \n")
		fmt.Printf("./NTLM-proxy -port 8080 -ntlm_url=http://NTLM-server:7047 -ntlm_username=username -ntlm_password=password -basicauth_username=clientuser -basicauth_password=clientpass \n")
	}
	flag.Parse()

	SERVER_PORT = getVar(port, SERVER_PORT)
	DEFAULT_NTLM_URL = getVar(ntlm_url, DEFAULT_NTLM_URL)
	DEFAULT_NTLM_USERNAME = getVar(ntlm_username, DEFAULT_NTLM_USERNAME)
	DEFAULT_NTLM_PASSWORD = getVar(ntlm_password, DEFAULT_NTLM_PASSWORD)
	BASICAUTH_USERNAME = getVar(basicauth_username, BASICAUTH_USERNAME)
	BASICAUTH_PASSWORD = getVar(basicauth_password, BASICAUTH_PASSWORD)
}
