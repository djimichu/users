package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorDto struct {
	Message string
	Time    time.Time
}

func (e ErrorDto) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
func WriteHttpError(w http.ResponseWriter, err error, code int) {
	errDto := ErrorDto{
		Message: err.Error(),
		Time:    time.Now(),
	}
	http.Error(w, errDto.ToString(), code)
}
