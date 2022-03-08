// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hcaptcha

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/flamego/flamego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.roundTrip(r)
}

func TestCaptcha(t *testing.T) {
	tests := []struct {
		name         string
		wantSecret   string
		wantToken    string
		wantRemoteIP string
	}{
		{
			name:         "normal",
			wantSecret:   "test-secret",
			wantToken:    "valid-token",
			wantRemoteIP: "",
		},
		{
			name:         "remoteip",
			wantSecret:   "test-secret",
			wantToken:    "valid-token",
			wantRemoteIP: "127.0.0.1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{
				Transport: &mockRoundTripper{
					roundTrip: func(r *http.Request) (*http.Response, error) {
						assert.Equal(t, test.wantSecret, r.PostFormValue("secret"))
						assert.Equal(t, test.wantToken, r.PostFormValue("response"))
						assert.Equal(t, test.wantRemoteIP, r.PostFormValue("remoteip"))
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
							Request:    r,
						}, nil
					},
				},
			}

			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(Captcha(Options{
				Client: client,
				Secret: test.wantSecret,
			}))
			f.Post("/", func(r *http.Request, h HCaptcha) {
				token := r.PostFormValue("h-captcha-response")

				var err error
				var resp *Response
				if test.wantRemoteIP != "" {
					resp, err = h.Verify(token, test.wantRemoteIP)
				} else {
					resp, err = h.Verify(token)
				}
				require.NoError(t, err)
				assert.True(t, resp.Success)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/", nil)
			require.NoError(t, err)

			req.PostForm = url.Values{
				"h-captcha-response": {test.wantToken},
			}
			f.ServeHTTP(resp, req)
		})
	}
}
