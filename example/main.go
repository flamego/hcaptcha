// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/flamego/flamego"

	"github.com/flamego/hcaptcha"
)

func main() {
	siteKey := flag.String("site-key", "", "The hCaptcha site key")
	secret := flag.String("secret", "", "The hCaptcha account secret")
	flag.Parse()

	f := flamego.Classic()
	f.Use(hcaptcha.Captcha(
		hcaptcha.Options{
			Secret: *secret,
		},
	))
	f.Get("/", func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		_, _ = w.Write([]byte(fmt.Sprintf(`
<html>
<head>
	<script src="https://hcaptcha.com/1/api.js"></script>
</head>
<body>
	<form method="POST">
		<div class="h-captcha" data-sitekey="%s"></div>
		<input type="submit" name="button" value="Submit">
	</form>
</body>
</html>
`, *siteKey)))
	})

	f.Post("/", func(w http.ResponseWriter, r *http.Request, h hcaptcha.HCaptcha) {
		token := r.PostFormValue("h-captcha-response")
		resp, err := h.Verify(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		} else if !resp.Success {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Verification failed, error codes %v", resp.ErrorCodes)))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Verified!"))
	})

	f.Run()
}
