package adf

type App struct {
	Dir AppDir
}

type AppDir struct {
	Root  string
	Certs string
}

type TLSOptions struct {
	CaFile  string
	KeyFile string
	CrtFile string
}
