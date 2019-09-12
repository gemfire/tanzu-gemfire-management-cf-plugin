package common

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Requester struct{}

func (requester *Requester) Exchange(url string, method string, bodyReader io.Reader, username string, password string) (urlResponse string, err error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return getURLOutput(resp)
}

func getURLOutput(resp *http.Response) (urlResponse string, err error) {
	respInASCII, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInASCII)

	return
}
