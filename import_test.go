package main

import (
	"testing"

	"github.com/go-zoox/testify"
)

func TestBuildGoImport(t *testing.T) {
	expected := `<!DOCTYPE html>
<html>
  <head>
    <meta name="go-import" content="zmicro.com/zoox git http://10.0.0.1:8080/go-zoox/zoox">
    <meta name="go-source" content="zmicro.com/zoox http://10.0.0.1:8080/go-zoox/zoox http://10.0.0.1:8080/go-zoox/zoox/tree/master{/dir} http://10.0.0.1:8080/go-zoox/zoox/blob/master{/dir}/{file}#L{line}">
  </head>
  <body></body>
</html>`
	testify.Equal(t, expected, BuildGoImport("zmicro.com/zoox", "http://10.0.0.1:8080/go-zoox/zoox"))
}
