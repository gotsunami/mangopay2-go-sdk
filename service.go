// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package mango is a library for the MangoPay service v2.
//
// MangoPay is a platform that allows to accept payments and manage e-money
// using wallets. See http://www.mangopay.com.
//
// First, create an account with a unique Id to authenticate to the service:
//  conf, err := mango.RegisterClient("myclientid", "My Company",
//      "contact@company.com", mango.Sandbox)
//  if err != nil {
//      panic(err)
//  }
// Or use existing credentials:
//  conf, err := mango.NewConfig("myclientid", "My Company",
//      "contact@company.com", "passwd", "sandbox")
//
// Then, choose an authentication mode (OAuth2.0 or Basic) to use with the service:
//  service, err := mango.NewMangoPay(conf, mango.OAuth)
package mango

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
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

// The default HTTP client to use with the MangoPay api.
var DefaultClient = &http.Client{
	Transport: &http.Transport{
		// Use TLS 1.1 as maximum version acceptable. TLS 1.2 (which is the
		// default used if not specified) seems no more supported by MangoPay
		// servers (used to be working though). Using TLS 1.2 results in
		// "connection reset by peer" errors.
		TLSClientConfig: &tls.Config{MaxVersion: tls.VersionTLS11},
	},
}

// The Mangopay service.
type MangoPay struct {
	clientId   string // MangoPay partner ID
	password   string
	env        ExecEnvironment // Live or testing env
	rootURL    *url.URL        // Base API URL for the current execution environment
	verbosity  Level
	authMethod AuthMode
	// To track the current token during its lifetime
	oauth *oAuth2
}

// ProcessIdent identifies the current operation.
type ProcessIdent struct {
	Id           string
	Tag          string
	CreationDate int64
}

// ProcessReply holds commong fields part of MangoPay API replies.
type ProcessReply struct {
	ProcessIdent
	Status        string
	ResultCode    string
	ResultMessage string
	ExecutionDate int64
}

// NewMangoPay creates a suitable environment for accessing
// the web service. Default verbosity level is set to Info, which can be
// changed through the use of Option().
func NewMangoPay(auth *Config, mode AuthMode) (*MangoPay, error) {
	if auth == nil {
		return nil, errors.New("nil config")
	}
	if auth.Env == "sandbox" {
		auth.env = Sandbox
	} else if auth.Env == "production" {
		auth.env = Production
	} else {
		return nil, errors.New("unknown exec environment: " + auth.Env)
	}
	u, err := url.Parse(rootURLs[auth.env])
	if err != nil {
		return nil, err
	}
	return &MangoPay{auth.ClientId, auth.Passphrase, auth.env, u, Info, mode, nil}, nil
}

// Option set various options like verbosity etc.
func (m *MangoPay) Option(opts ...option) {
	for _, opt := range opts {
		opt(m)
	}
}

// request prepares and sends a well formatted HTTP request to the
// mangopay service.
func (s *MangoPay) request(ma mangoAction, data JsonObject) (*http.Response, error) {
	mr, ok := mangoRequests[ma]
	if !ok {
		return nil, errors.New("Action not implemented.")
	}

	// Create the submit url
	path := mr.Path
	if mr.PathValues != nil {
		// Substitute path variables, if any
		for name, _ := range mr.PathValues {
			if _, ok := data[name]; !ok {
				return nil, errors.New(fmt.Sprintf("missing keyword %s", name))
			}
			path = strings.Replace(path, "{{"+name+"}}", fmt.Sprintf("%v", data[name]), -1)
		}
	} else {
		path = mr.Path
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := s.rawRequest(mr.Method, "application/json",
		fmt.Sprintf("%s%s%s", s.rootURL, s.clientId, path), body, true)
	return resp, err
}

// rawRequest sends an HTTP request with method method to an arbitrary URI.
func (s *MangoPay) rawRequest(method, contentType string, uri string, body []byte, useAuth bool) (*http.Response, error) {
	if contentType == "" {
		return nil, errors.New("empty request's content type")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Set header for basic auth
	if useAuth {
		if s.authMethod == BasicAuth {
			req.Header.Set("Authorization", basicAuthorization(s.clientId, s.password))
		} else {
			o, err := oAuthAuthorization(s)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", o)
		}
	}
	req.Header.Set("Content-Type", contentType)

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
	resp, err := DefaultClient.Do(req)

	// Handle reponse status code
	if err == nil && resp.StatusCode != http.StatusOK {
		j := JsonObject{}
		err = s.unMarshalJSONResponse(resp, &j)
		if err != nil {
			return nil, err
		}
		errmsg := ""
		if msg, ok := j["Message"]; ok {
			errmsg = fmt.Sprintf("Status %d: %s", resp.StatusCode, msg.(string))
		} else {
			errmsg = fmt.Sprintf("Status %d; body: '%s'", resp.StatusCode, j)
		}
		if private, ok := j["errors"]; ok {
			errmsg += fmt.Sprintf("(details :%v)", private)
		}
		err = errors.New(errmsg)
	}
	return resp, err
}

// Unmarshal a JSON HTTP response into an instance.
func (m *MangoPay) unMarshalJSONResponse(resp *http.Response, v interface{}) error {
	if resp == nil {
		return errors.New("can't unmarshal nil response")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if m.verbosity == Debug {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>> DEBUG RESPONSE")
		fmt.Printf("Status code: %d\n\n", resp.StatusCode)
		for k, v := range resp.Header {
			for _, j := range v {
				fmt.Printf("%s: %v\n", k, j)
			}
		}
		fmt.Printf("\n%s\n", string(b))
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<< DEBUG RESPONSE")
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

// Generic request for any object.
func (m *MangoPay) anyRequest(o interface{}, action mangoAction, data JsonObject) (interface{}, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(o)
		t = reflect.Indirect(v).Type()
	}
	ins := reflect.New(t).Interface()
	if err := m.unMarshalJSONResponse(resp, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func unixTimeToString(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).String()
	}
	return "Never"
}

// Use reflection to print data structures.
func struct2string(c interface{}) string {
	var b bytes.Buffer
	e := reflect.ValueOf(c).Elem()
	for i := 0; i < e.NumField(); i++ {
		sfield := e.Type().Field(i)
		// Skip unexported fields
		if sfield.PkgPath != "" {
			continue
		}
		name := sfield.Name
		val := e.Field(i).Interface()
		// Handle embedded types
		if sfield.Anonymous {
			b.Write([]byte(struct2string(e.Field(i).Addr().Interface())))
		} else {
			if name == "CreationDate" || name == "ExecutionDate" ||
				name == "Birthday" {
				val = unixTimeToString(val.(int64))
			}
			b.Write([]byte(fmt.Sprintf("%-24s: %v\n", name, val)))
		}
	}
	return b.String()
}

func consumerId(c Consumer) string {
	id := ""
	switch c.(type) {
	case *LegalUser:
		id = c.(*LegalUser).Id
	case *NaturalUser:
		id = c.(*NaturalUser).Id
	}
	return id
}
