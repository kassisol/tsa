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

var EngineCertsDir = "/etc/docker/tls"

var EngineCaFile = path.Join(EngineCertsDir, "ca.pem")
var EngineKeyFile = path.Join(EngineCertsDir, "server-key.pem")
var EngineCrtFile = path.Join(EngineCertsDir, "server-cert.pem")
