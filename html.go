package email

import (
	"bytes"
	"fmt"
	"html"
)

// HTML creates a complete email compatible XHTML 1.0 document
// around the passed bodyTag.
func HTML(title string, bodyTag []byte) []byte {
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, htmlHead, html.EscapeString(title))
	b.Write(bodyTag)
	b.WriteString("\n</html>")
	return b.Bytes()
}

const htmlHead = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
<meta name="viewport" content="width=device-width"/>
<title>%s</html>
</head>
`
