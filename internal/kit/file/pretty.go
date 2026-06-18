package file

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"path/filepath"
	"strings"
)

func formatByFileExt(content []byte, filePath string) []byte {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		if b, ok := prettyJSONBytes(content); ok {
			return b
		}
	case ".xml":
		if b, ok := prettyXMLBytes(content); ok {
			return b
		}
	}
	// plain / unknown ext
	return content
}

func prettyJSONBytes(b []byte) ([]byte, bool) {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return []byte{}, true
	}
	if b[0] != '{' && b[0] != '[' {
		return nil, false
	}
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		return nil, false
	}
	out.WriteByte('\n')
	return out.Bytes(), true
}

func prettyXMLBytes(b []byte) ([]byte, bool) {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return []byte{}, true
	}
	if b[0] != '<' {
		return nil, false
	}

	var out bytes.Buffer
	dec := xml.NewDecoder(bytes.NewReader(b))
	enc := xml.NewEncoder(&out)
	enc.Indent("", "  ")

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if err := enc.EncodeToken(tok); err != nil {
			return nil, false
		}
	}
	if err := enc.Flush(); err != nil {
		return nil, false
	}

	out.WriteByte('\n')
	return out.Bytes(), true
}
