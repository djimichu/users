package response

import (
	"encoding/json"
	"log"
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
func WriteJSON(w http.ResponseWriter, data any, code int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("err encode json:", err)
	}
}
