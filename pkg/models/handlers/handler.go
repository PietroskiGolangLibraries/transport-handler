package handlers_model

type (
	Server interface {
		Handle()
		Start() error
	}
)
