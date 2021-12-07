package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverCertFile = "../certs/out/localhost.crt"
	serverKeyfile  = "../certs/out/localhost.key"
	caCertFile     = "../certs/out/myca.crt"

	//certopt    Optional, specifies the option for authenticating a client via certificate:
	//0 - certificate not required,
	//1 - request a certificate but it's not required,
	//2 - require any client certificate
	//3 - if provided, verify the client certificate is authorized
	//4 - require certificate and verify it's authorized`
	certOpt = 0
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	log.Print("start http server")

	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	server := &http.Server{
		Addr:         ":8443",
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		//TLSConfig:    &tls.Config{ServerName: "localhost"},
		TLSConfig: getTLSConfig("localhost", caCertFile, tls.RequireAndVerifyClientCert),
	}

	// Listen to port 8080 and wait
	//log.Fatal(http.ListenAndServe(":8081", nil))

	if err := server.ListenAndServeTLS(serverCertFile, serverKeyfile); err != nil {
		log.Fatalf("server listen failed.%v", err)
	}

}

func getTLSConfig(host, caCertFile string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Fatal("Error opening cert file", caCertFile, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		// ClientAuth: tls.NoClientCert,				// Client certificate will not be requested and it is not required
		// ClientAuth: tls.RequestClientCert,			// Client certificate will be requested, but it is not required
		// ClientAuth: tls.RequireAnyClientCert,		// Client certificate is required, but any client certificate is acceptable
		// ClientAuth: tls.VerifyClientCertIfGiven,		// Client certificate will be requested and if present must be in the server's Certificate Pool
		// ClientAuth: tls.RequireAndVerifyClientCert,	// Client certificate will be required and must be present in the server's Certificate Pool
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}
