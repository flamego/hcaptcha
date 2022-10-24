# hcaptcha

[![GitHub Workflow Status](https://img.shields.io/github/checks-status/flamego/hcaptcha/main?logo=github&style=for-the-badge)](https://github.com/flamego/hcaptcha/actions?query=branch%3Amain)
[![Codecov](https://img.shields.io/codecov/c/gh/flamego/hcaptcha?logo=codecov&style=for-the-badge)](https://app.codecov.io/gh/flamego/hcaptcha)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/hcaptcha?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/flamego/hcaptcha)

Package hcaptcha is a middleware that provides hCaptcha rendering integration for [Flamego](https://github.com/flamego/flamego).

## Installation

The minimum requirement of Go is **1.18**.

	go get github.com/flamego/hcaptcha

## Getting started

```html
<!-- templates/home.tmpl -->
<html>
<head>
  <script src="https://hcaptcha.com/1/api.js"></script>
</head>
<body>
  <form method="POST">
    <div class="h-captcha" data-sitekey="{{.SiteKey}}"></div>
    <input type="submit" name="button" value="Submit">
  </form>
</body>
</html>
```

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/hcaptcha"
	"github.com/flamego/template"
)

func main() {
	f := flamego.Classic()
	f.Use(template.Templater())
	f.Use(hcaptcha.Captcha(
		hcaptcha.Options{
			Secret: "<SECRET>",
		},
	))
	f.Get("/", func(t template.Template, data template.Data) {
		data["SiteKey"] = "<SITE KEY>"
		t.HTML(http.StatusOK, "home")
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
```

## Getting help

- Read [documentation and examples](https://flamego.dev/middleware/hcaptcha.html).
- Please [file an issue](https://github.com/flamego/flamego/issues) or [start a discussion](https://github.com/flamego/flamego/discussions) on the [flamego/flamego](https://github.com/flamego/flamego) repository.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
