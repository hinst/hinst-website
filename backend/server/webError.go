package server

type webError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (me webError) Error() string {
	return me.Message
}

func assertWebError(err error, webErr webError) {
	if err != nil {
		panic(webErr)
	}
}
