package util

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HTTPClientIface interface {
	Get(*http.Request) ([]byte, int, error)
	Post(*http.Request) ([]byte, int, error)
	Do(*http.Request) ([]byte, int, error)
}

type HTTPClient struct {
	Timeout int
	client  *http.Client
}

func NewHTTPClient(timeout int, client *http.Client) *HTTPClient {
	var c = http.DefaultClient
	if client != nil {
		c = client
	}
	hc := &HTTPClient{
		Timeout: timeout,
		client:  c,
	}
	hc.client.Timeout = time.Duration(timeout) * time.Second

	return hc
}

// NewHTTPReqeust makes a http request
func NewHTTPReqeust(method, url string, queries, headers map[string]string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if queries != nil {
		query := req.URL.Query()
		for key, val := range queries {
			query.Set(key, val)
		}
		req.URL.RawQuery = query.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}

	if body != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	return req, nil
}

func (hc *HTTPClient) Do(req *http.Request) ([]byte, int, error) {
	res, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("resp code is %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, res.StatusCode, err
	}

	return b, res.StatusCode, nil
}

// Get returns http response body in []byte, timeout in second
func (hc *HTTPClient) Get(req *http.Request) ([]byte, int, error) {
	req.Method = "GET"
	return hc.Do(req)
}

// Post do http post
func (hc *HTTPClient) Post(req *http.Request) ([]byte, int, error) {
	req.Method = "POST"
	return hc.Do(req)
}

// GetReadCloser for downloading file
func (hc *HTTPClient) GetReadCloser(req *http.Request) (io.ReadCloser, string, int, error) {
	var err error
	var res *http.Response

	res, err = hc.client.Do(req)
	if err != nil {
		return nil, "", 0, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, "", 0, fmt.Errorf("resp code is %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	contentLength, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return nil, "", 0, err
	}

	return res.Body, contentType, contentLength, nil
}

func MustLoadCertificates(privateKeyFile, certificateFile, caFile string) (tls.Certificate, *x509.CertPool) {
	mycert, err := tls.LoadX509KeyPair(certificateFile, privateKeyFile)
	if err != nil {
		panic(err)
	}
	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return mycert, certPool
}

func MustGetTlsConfiguration(privateKeyFile, certificateFile, caFile string) *tls.Config {
	config := &tls.Config{}
	mycert, certPool := MustLoadCertificates(privateKeyFile, certificateFile, caFile)
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = mycert

	config.RootCAs = certPool
	config.ClientCAs = certPool
	config.InsecureSkipVerify = true

	//config.ClientAuth = tls.RequireAndVerifyClientCert

	//Optional stuff

	//Use only modern ciphers
	config.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}

	//Use only TLS v1.2
	config.MinVersion = tls.VersionTLS12

	//Don't allow session resumption
	// config.SessionTicketsDisabled = true
	return config
}
