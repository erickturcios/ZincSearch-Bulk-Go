package service

import (
	"log"
	"net/http"
	"strings"
	"time"

	"zincsearch.com/mailindex/api/helpers"
)

// guarda indice
func (s *ZincSearch) CreateDocumentBulk(jsonBody string) (result string, httpError helpers.ErrorResponse) {
	const resource string = "/api/_bulkv2"

	h := http.Client{Timeout: 20 * time.Second}

	//obtiene string del URL
	url := helpers.GetUrl(s.https, s.host, s.port, resource, "")

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))

	if err != nil {
		return result, helpers.GetErrorResponse(-1, err.Error())
	}

	s.debugReq(req)

	//agrega credenciales para autenticacion
	helpers.AddBasicAuth(req, s.usuario, s.password)
	//ejecuta peticion
	response, err := h.Do(req)

	if err != nil {
		return result, helpers.GetErrorResponse(-1, err.Error())
	}

	s.debugRes(response)

	if response.Body != nil {
		defer response.Body.Close()
	}

	//obtiene resultado como string
	if response.StatusCode == 200 {
		result, err = helpers.GetResponseString(response)
		if err != nil {
			return result, helpers.GetErrorResponse(-1, err.Error())
		}
	} else {
		httpError, err := helpers.GetError(response)
		if err != nil {
			log.Fatal(err)
		}
		return result, httpError

	}

	//retorna respuesta
	return result, httpError
}
