package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	clientCertFile = "../certs/out/client.crt"
	clientKeyFile  = "../certs/out/client.key"
	caCertFile     = "../certs/out/myca.crt"
)

func main() {
	// Request /hello over port 8080 via the GET method
	//r, err := http.Get("http://localhost:8081/hello")

	// failed: x509: certificate relies on legacy Common Name field, use SANs instead
	//r, err := http.Get("https://localhost:8443/hello")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//ignoreTLS()
	verifyTLS()
}

func ignoreTLS() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", body)
}

func verifyTLS() {
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
	}
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}
	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatalf("client request err:%v", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalf("unexpected error reading response body: %s", err)
	}

	fmt.Printf("\nResponse from server: \n\tHTTP status: %s\n\tBody: %s\n", r.Status, body)
}
