package adf

type App struct {
	Dir AppDir
}

type AppDir struct {
	Root    string
	Tenants string
	Certs   string
}

type TLSOptions struct {
	CaFile  string
	KeyFile string
	CrtFile string
}
