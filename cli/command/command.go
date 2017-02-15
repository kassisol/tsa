package command

import (
	"path"
)

var AppPath = "/var/lib/tsa"

var DBFilePath = path.Join(AppPath, "data.db")

var CaDir = path.Join(AppPath, "ca")
var CaPrivateDir = path.Join(CaDir, "private")
var CaCertsDir = path.Join(CaDir, "certs")

var CaCrtFile = path.Join(CaCertsDir, "ca.crt")
var CaCrlFile = path.Join(CaDir, "CRL.crl")

var ApiCertsDir = path.Join(AppPath, "certs")

var ApiKeyFile = path.Join(ApiCertsDir, "api.key")
var ApiCrtFile = path.Join(ApiCertsDir, "api.crt")
