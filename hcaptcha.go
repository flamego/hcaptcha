// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hcaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/flamego/flamego"
)

// Options contains options for both hcaptcha.Captcha middleware.
type Options struct {
	// Secret is the secret key to check user captcha codes. This field is required.
	Secret string
}

// Captcha returns a middleware handler that injects hcaptcha.HCaptcha into the
// request context, which is used for verifying hCaptcha requests.
func Captcha(opts Options) flamego.Handler {
	if opts.Secret == "" {
		panic("hcaptcha: empty secret")
	}

	return flamego.ContextInvoker(func(c flamego.Context) {
		hcaptcha := &hCaptcha{
			secret: opts.Secret,
		}
		c.MapTo(hcaptcha, (*HCaptcha)(nil))
	})
}

// HCaptcha is a hCaptcha verify service.
type HCaptcha interface {
	// Verify verifies the given token. An optional remote IP of the user may be
	// passed as extra security criteria.
	Verify(token string, remoteIP ...string) (*Response, error)
}

// Response is the response struct which hCaptcha sends back.
type Response struct {
	// Success indicates whether the passcode valid, and does it meet security
	// criteria you specified.
	Success bool `json:"success"`
	// ChallengeTS is the timestamp of the challenge (ISO format
	// yyyy-MM-dd'T'HH:mm:ssZZ).
	ChallengeTS time.Time `json:"challenge_ts"`
	// Hostname is the hostname of the site where the challenge was solved.
	Hostname string `json:"hostname"`
	// Credit indicates whether the response will be credited.
	Credit bool `json:"credit"`
	// ErrorCodes contains the error codes when verify failed.
	ErrorCodes []string `json:"error-codes"`
}

var _ HCaptcha = (*hCaptcha)(nil)

type hCaptcha struct {
	secret string
}

func (h *hCaptcha) Verify(token string, remoteIP ...string) (*Response, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	data := url.Values{
		"secret":   {h.secret},
		"response": {token},
	}
	if len(remoteIP) > 0 {
		data.Add("remoteip", remoteIP[0])
	}

	resp, err := http.PostForm("https://hcaptcha.com/siteverify", data)
	if err != nil {
		return nil, errors.Wrap(err, "request hCaptcha server")
	}
	defer func() { _ = resp.Body.Close() }()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decode response body")
	}
	return &response, nil
}
