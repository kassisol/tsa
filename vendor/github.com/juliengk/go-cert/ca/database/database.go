package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/juliengk/go-cert/ca/database/backend"
)

type Initialize func(string) (backend.Backender, error)

var initializers = make(map[string]Initialize)

var supportedBackends = func() string {
	backends := make([]string, 0, len(initializers))

	for b := range initializers {
		backends = append(backends, string(b))
	}

	sort.Strings(backends)

	return strings.Join(backends, ",")
}()

func NewBackend(backend, config string) (backend.Backender, error) {
	if init, exists := initializers[backend]; exists {
		return init(config)
	}

	return nil, fmt.Errorf("The Database Backend: %s is not supported. Supported backends are %s", backend, supportedBackends)
}

func RegisterBackend(backend string, init Initialize) {
	initializers[backend] = init
}
