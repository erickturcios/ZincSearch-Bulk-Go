package helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code  int
	Error string
}

// Obtiene respuesta HTTP como string
func GetResponseString(response *http.Response) (str string, err error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return str, err
	}
	str = string(body)
	return str, nil
}

// obtiene respuesta HTTP como []byte
func GetResponseBytes(response *http.Response) (str []byte, err error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return str, err
	}
	return body, nil
}

func GetError(response *http.Response) (httpError ErrorResponse, err error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return httpError, err
	}
	err = json.Unmarshal(body, &httpError)
	if err != nil {
		log.Fatal(err)
	}
	httpError.Code = response.StatusCode

	return httpError, err

}

func GetErrorResponse(code int, description string) (err ErrorResponse) {
	err.Code = code
	err.Error = description
	return err
}
