// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package mango is a library for the MangoPay service v2.
//
// http://www.mangopay.com
package mango

import (
	/*
		"bytes"
		"crypto"
		"encoding/pem"
		"io"
		"os"
		"strconv"
		"time"
	*/
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request execution environment (production or sandbox).
type ExecEnvironment int

const (
	Production ExecEnvironment = iota
	Sandbox
)

// Base URLs to execution environements
var rootURLs = map[ExecEnvironment]string{
	Production: "https://api.mangopay.com/v2/",
	Sandbox:    "https://api.sandbox.mangopay.com/v2/",
}

// The Mangopay service.
type MangoPay struct {
	clientId string // MangoPay partner ID
	password string
	env      ExecEnvironment // Live or testing env
	rootURL  *url.URL        // Base API URL for the current execution environment
}

// MangoPay object identity data.
type Identity struct {
	ID                       int
	CreationDate, UpdateDate int
}

// NewMangoPay creates a suitable environment for sending accessing
// the web service.
func NewMangoPay(clientId, password string, env ExecEnvironment) (*MangoPay, error) {
	u, err := url.Parse(rootURLs[env])
	if err != nil {
		return nil, err
	}
	return &MangoPay{clientId, password, env, u}, nil
}

/*
func request() {
	resp, err := s.handleResponse(action, data)
	if err != nil {
		return nil, err
	}
	u := new(User)
	if err := unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil

}
*/

// makeHashedKey hashes a string with SHA-1 algorithm and returns a
// 20 bytes hash.
//
// String to sign for GET methods :		{method}|{urlPath}|
// For POST and PUT methods :    		{method}|{urlPath}|{request body}|
/*
func (s *MangoPay) makeHashedKey(method, path string, data JsonObject) ([]byte, error) {
	if method != "GET" {
		if data == nil {
			return nil, errors.New("missing body data for POST or PUT request")
		}
	}
	key := fmt.Sprintf("%s|%s|", method, path)
	if method != "GET" {
		marshaled, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		key += string(marshaled) + "|"
	}

	// Compute hash
	hash := sha1.New()
	io.WriteString(hash, key)

	return hash.Sum(nil), nil
}

// makeSignature computes a request signature, supplied in the request's header
// with the X-Leetchi-Signature key.
//
// The signature is calculated by hashing a string with SHA-1 algorithm,
// signing it with the client’s RSA key and converting to Base64.
func (s *MangoPay) makeSignature(method, path string, data JsonObject) (string, error) {
	hashed, err := s.makeHashedKey(method, path, data)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.pkey, crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}

	// Convert to base64
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(signature)
	encoder.Close()
	return buf.String(), nil
}
*/

// request prepares and sends a well formatted HTTP request to the
// mangopay service.
// TODO: only basic access auth supported at the moment. Add support
// for OAuth2.0.
func (s *MangoPay) request(ma MangoAction, data JsonObject) (*http.Response, error) {
	mr, ok := mangoRequests[ma]
	if !ok {
		return nil, errors.New("Action not implemented.")
	}

	// Create the submit url
	path := ""
	if mr.PathValues != nil {
		// Substitute path variables, if any
		for name, _ := range mr.PathValues {
			if _, ok := data[name]; !ok {
				return nil, errors.New(fmt.Sprintf("missing keyword %s", name))
			}
			path = strings.Replace(mr.Path, "{{"+name+"}}", fmt.Sprintf("%v", data[name]), -1)
		}
	} else {
		path = mr.Path
	}

	uri, err := url.Parse(fmt.Sprintf("%s%s%s", s.rootURL, s.clientId, path))
	if err != nil {
		return nil, err
	}
	fmt.Println(uri)

	// Create a request with body
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(mr.Method, uri.String(), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Set header for basic auth
	credential := s.clientId + ":" + s.password
	credential = base64.StdEncoding.EncodeToString([]byte(credential))
	req.Header.Set("Authorization", "Basic "+credential)

	// Post the data
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	// Handle reponse status code
	if resp.StatusCode != http.StatusOK {
		j := JsonObject{}
		err = unMarshalJSONResponse(resp, &j)
		if err != nil {
			return nil, err
		}
		if msg, ok := j["Message"]; ok {
			err = errors.New(msg.(string))
		} else {
			err = errors.New(fmt.Sprintf("HTTP status %d; body: '%s'", resp.StatusCode, j))
		}
	}
	return resp, err
}

// Unmarshal a JSON HTTP response into an instance.
func unMarshalJSONResponse(resp *http.Response, v interface{}) error {
	if resp == nil {
		return errors.New("can't unmarshal nil response")
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
