package service

import (
	"log"
	"net/http"
	"strings"
	"time"

	"zincsearch.com/mailindex/api/helpers"
)

type IndexListRequest struct {
	Page_num  int
	Page_size int
	Sort_by   string
	Desc      bool
	Name      string
}

func (s *ZincSearch) GetIndexList(request IndexListRequest) (result string, httpError helpers.ErrorResponse) {
	const resource string = "/api/index"

	h := http.Client{Timeout: 20 * time.Second}

	//obtiene string del URL
	urlQuery := helpers.GetUrlQueryFromStruct(request)
	url := helpers.GetUrl(s.https, s.host, s.port, resource, urlQuery)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, helpers.GetErrorResponse(-1, err.Error())
	} else {
		req.Header.Add("Content-Type", "application/json")
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

func (s *ZincSearch) ExistsIndex(indexName string) (result bool, httpError helpers.ErrorResponse) {
	var sb strings.Builder
	sb.WriteString("/api/index/")
	sb.WriteString(indexName)

	h := http.Client{Timeout: 20 * time.Second}

	//obtiene string del URL
	url := helpers.GetUrl(s.https, s.host, s.port, sb.String(), "")

	req, err := http.NewRequest(http.MethodHead, url, nil)

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
	result = response.StatusCode == 200

	//retorna respuesta
	return result, httpError

}

func (s *ZincSearch) GetIndex(indexName string) (result string, httpError helpers.ErrorResponse) {
	var sb strings.Builder
	sb.WriteString("/api/index/")
	sb.WriteString(indexName)

	h := http.Client{Timeout: 20 * time.Second}

	//obtiene string del URL
	url := helpers.GetUrl(s.https, s.host, s.port, sb.String(), "")

	req, err := http.NewRequest(http.MethodGet, url, nil)

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

// guarda indice
func (s *ZincSearch) SaveIndex(indexName string, jsonBody string) (result string, httpError helpers.ErrorResponse) {
	const resource string = "/api/index"

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

// Elimina indice
func (s *ZincSearch) DeleteIndex(indexName string) (result string, httpError helpers.ErrorResponse) {
	var sb strings.Builder
	sb.WriteString("/api/index/")
	sb.WriteString(indexName)

	h := http.Client{Timeout: 20 * time.Second}

	//obtiene string del URL
	url := helpers.GetUrl(s.https, s.host, s.port, sb.String(), "")

	req, err := http.NewRequest(http.MethodDelete, url, nil)

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
