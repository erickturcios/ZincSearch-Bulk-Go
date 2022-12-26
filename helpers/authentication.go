package helpers

import "net/http"

//agrega credenciales para autenticacion basica
func AddBasicAuth(req *http.Request, usuario string, password string) {
	if usuario != "" && password != "" {
		req.SetBasicAuth(usuario, password)
	}
}
