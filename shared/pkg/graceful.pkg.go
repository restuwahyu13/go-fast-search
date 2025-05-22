package pkg

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/ory/graceful"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

func Graceful(req dto.Request[dto.Environtment], Handler func() opt.Graceful) error {
	h := Handler()
	secure := true

	if _, ok := os.LookupEnv("GO_ENV"); ok && req.Config.APP.ENV != cons.DEV {
		secure = false
	}

	server := http.Server{
		Handler:        h.HANDLER,
		Addr:           ":" + h.ENV.APP.PORT,
		MaxHeaderBytes: req.Config.APP.INBOUND_SIZE,
		TLSConfig:      &tls.Config{InsecureSkipVerify: secure},
	}

	Logrus(cons.INFO, "Server listening on port %s", h.ENV.APP.PORT)
	return graceful.Graceful(server.ListenAndServe, server.Shutdown)
}
