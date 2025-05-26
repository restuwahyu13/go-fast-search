package opt

import (
	"github.com/go-chi/chi/v5"
)

type Graceful struct {
	HANDLER *chi.Mux
	ENV     *Environtment
}
