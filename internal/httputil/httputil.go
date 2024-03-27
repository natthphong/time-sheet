package httputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type HTTPPostRequestFunc func(reqBody interface{}, requestRef *string) ([]byte, error)

type HTTPPostPaymentRequestFunc func(url string, reqBody []byte, auth string, contentType string) ([]byte, error)
type HttpPostOddPaymentRequestFunc func(reqBody []byte) ([]byte, error)

func InitHttpClient(timeout time.Duration, maxIdleConn, maxIdleConnPerHost, maxConnPerHost int) *http.Client {
	certPool := x509.NewCertPool()
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        maxIdleConn,
			MaxIdleConnsPerHost: maxIdleConnPerHost,
			MaxConnsPerHost:     maxConnPerHost,
		},
	}
	return client
}

func InitHttpClientWithCertAndKey(timeout time.Duration, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int, CertFile, KeyFile []byte) (*http.Client, error) {

	cert, err := tls.X509KeyPair(CertFile, KeyFile)
	if err != nil {
		return nil, err
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	rootCAs.AppendCertsFromPEM(CertFile)
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            rootCAs,
				InsecureSkipVerify: false,
			},
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			MaxConnsPerHost:     maxConnsPerHost,
		},
	}
	return client, nil
}

func InitHttpClientWithCert(timeout time.Duration, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int, CertFile []byte) (*http.Client, error) {
	//cert, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if !rootCAs.AppendCertsFromPEM(CertFile) {
		fmt.Println("Cannot AppendCerts\n", string(CertFile))
	}
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAs,
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			MaxConnsPerHost:     maxConnsPerHost,
		},
	}
	return client, nil
}

func NewHttpPostCall(client *http.Client, url string) HTTPPostRequestFunc {
	return func(reqBody interface{}, requestRef *string) ([]byte, error) {

		message, _ := json.Marshal(&reqBody)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(message))
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request")
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Aurora-Secret", "Internal")
		if requestRef != nil {
			req.Header.Add("RequestRef", *requestRef)
		} else {
			req.Header.Add("RequestRef", uuid.NewString())
		}

		res, err := client.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unable to request %s", url))
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request ioutil")
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s", body)
		}

		return body, nil
	}
}

func NewHttpPostPaymentCall(client *http.Client) HTTPPostPaymentRequestFunc {
	return func(url string, reqBody []byte, auth string, contentType string) ([]byte, error) {

		//message, _ := json.Marshal(&reqBody)

		// make SSL certificate and key
		certFile := "star_allgold_arrgx_com.crt"
		keyFile := "_.allgold.arrgx.com.key"

		// read certificate and key from file
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Println("Can't download certificate และ key:", err)
			return nil, err
		}

		// Create HTTP client
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			fmt.Println("Cert error:", err)
			return nil, err
		}

		rootCAs.AppendCertsFromPEM([]byte(certFile))
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      rootCAs,
			},
		}
		client := &http.Client{Transport: transport}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request")
		}

		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("x-test-mode", "true")
		req.Header.Set("env-id", "OAUTH2")

		res, err := client.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unable to request %s", url))
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request ioutil")
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("StatusCode: %d ,Body: %s", res.StatusCode, body)
		}

		return body, nil
	}
}

func NewHttpPostOddPaymentCall(client *http.Client, url string) HttpPostOddPaymentRequestFunc {
	return func(reqBody []byte) ([]byte, error) {

		//message, _ := json.Marshal(&reqBody)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request")
		}

		req.Header.Add("Content-Type", "text/xml;charset=utf-8")

		req.Header.Set("Content-Type", "text/xml")
		req.Header.Set("SOAPAction", "")
		req.Header.Set("Access-Control-Allow-Headers", "Authorization")
		req.Header.Set("Access-Control-Allow-Methods", "POST")
		req.Header.Set("Access-Control-Allow-Origin", "*")

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request ioutil")
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s", body)
		}

		return body, nil
	}
}
