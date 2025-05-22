package con

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/wagslane/go-rabbitmq"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
)

func RabbitConnection(req dto.Request[dto.Environtment]) (*rabbitmq.Conn, error) {
	interval := time.Duration(time.Second * 5)

	return rabbitmq.NewConn(req.Config.RABBITMQ.URL,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(interval),
		rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
			Vhost:           req.Config.RABBITMQ.VSN,
			FrameSize:       http.DefaultMaxHeaderBytes * 5,
			Heartbeat:       interval,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}))
}
