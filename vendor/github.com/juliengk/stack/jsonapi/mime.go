package jsonapi

import (
	"fmt"
)

const (
	MIMEApplicationVendor = "application/vnd"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

func BuildVendorMIME(name string) string {
	return fmt.Sprintf("%s.%s+json", MIMEApplicationVendor, name)
}

func BuildVendorMIMECharsetUTF8(name string) string {
	return fmt.Sprintf("%s.%s+json; %s", MIMEApplicationVendor, name, charsetUTF8)
}
