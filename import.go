package main

import "fmt"

func BuildGoImport(importPrefix, repoRoot string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta name="go-import" content="%s git %s">
    <meta name="go-source" content="%s %s %s/tree/master{/dir} %s/blob/master{/dir}/{file}#L{line}">
  </head>
  <body></body>
</html>`, importPrefix, repoRoot, importPrefix, repoRoot, repoRoot, repoRoot)
}
