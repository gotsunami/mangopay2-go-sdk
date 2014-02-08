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
	clientId   string // MangoPay partner ID
	password   string
	env        ExecEnvironment // Live or testing env
	rootURL    *url.URL        // Base API URL for the current execution environment
	verbosity  Level
	authMethod AuthMode
}

// MangoPay object identity data.
type Identity struct {
	ID                       int
	CreationDate, UpdateDate int
}

// NewMangoPay creates a suitable environment for accessing
// the web service. Default verbosity level is set to Info, default authentication
// mode to BasicAuth. They can be changed through the use of Option().
func NewMangoPay(clientId, password string, env ExecEnvironment) (*MangoPay, error) {
	u, err := url.Parse(rootURLs[env])
	if err != nil {
		return nil, err
	}
	return &MangoPay{clientId, password, env, u, Info, BasicAuth}, nil
}

// Option set various options like verbosity etc.
func (m *MangoPay) Option(opts ...option) {
	for _, opt := range opts {
		opt(m)
	}
}

// request prepares and sends a well formatted HTTP request to the
// mangopay service.
// TODO: only basic access auth supported at the moment. Add support
// for OAuth2.0.
func (s *MangoPay) request(ma mangoAction, data JsonObject) (*http.Response, error) {
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

	if s.verbosity == Debug {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>> DEBUG REQUEST")
		fmt.Printf("%s %s\n\n", req.Method, req.URL.String())
		for k, v := range req.Header {
			for _, j := range v {
				fmt.Printf("%s: %v\n", k, j)
			}
		}
		rb := string(body)
		if rb != "null" {
			fmt.Printf("\n%s\n", rb)
		}
		fmt.Println("\nSending request ...")
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<< DEBUG REQUEST")
	}

	// Send request
	resp, err := http.DefaultClient.Do(req)

	// Handle reponse status code
	if resp.StatusCode != http.StatusOK {
		j := JsonObject{}
		err = s.unMarshalJSONResponse(resp, &j)
		if err != nil {
			return nil, err
		}
		if msg, ok := j["Message"]; ok {
			err = errors.New(fmt.Sprintf("%s (%d)", msg.(string), resp.StatusCode))
		} else {
			err = errors.New(fmt.Sprintf("HTTP status %d; body: '%s'", resp.StatusCode, j))
		}
	}
	return resp, err
}

// Unmarshal a JSON HTTP response into an instance.
func (m *MangoPay) unMarshalJSONResponse(resp *http.Response, v interface{}) error {
	if resp == nil {
		return errors.New("can't unmarshal nil response")
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if m.verbosity == Debug {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>> DEBUG RESPONSE")
		fmt.Println("Raw body:", string(b))
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<< DEBUG RESPONSE")
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
