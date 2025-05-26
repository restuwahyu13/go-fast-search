package opt

type (
	Application struct {
		ENV          string
		PORT         string
		INBOUND_SIZE int
	}

	Redis struct {
		URL string
	}

	Postgres struct {
		URL string
	}

	Jwt struct {
		SECRET  string
		EXPIRED int
	}

	RabbitMQ struct {
		URL    string
		VSN    string
		SECRET string
	}

	MeiliSearch struct {
		URL string
		KEY string
	}

	Environtment struct {
		APP         Application
		REDIS       Redis
		POSTGRES    Postgres
		JWT         Jwt
		RABBITMQ    RabbitMQ
		MEILISEARCH MeiliSearch
	}
)
