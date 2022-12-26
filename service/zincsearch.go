package service

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

const USUARIO string = "ZINC_FIRST_ADMIN_USER"
const PSWD string = "ZINC_FIRST_ADMIN_PASSWORD"
const HOST string = "ZINC_SERVER_HOST" //variable de entorno custom
const ZincSearchPort string = "ZINC_SERVER_PORT"
const ZincSearchHttps string = "ZINC_SERVER_HTTPS"
const defaultHost string = "localhost"
const defaultPort string = "4080"

// nombre por defecto de indice utilizado por este aplicativo
const INDEX_NAME string = "mailindex"

var debugEnabled bool = false

type ZincSearch struct {
	usuario  string
	password string
	host     string
	port     string
	https    bool
}

// inicializa configuracion para ejecucion de peticiones hace ZincSearch
func (s *ZincSearch) Inicia() {
	//obtener credenciales de variables de ambiente que utiliza zinc search
	s.usuario = os.Getenv(USUARIO)
	s.password = os.Getenv(PSWD)

	if s.usuario == "" {
		log.Fatal("No esta definida la variable de ambiente ", USUARIO, " para el usuario de ZincSearch")
	}
	if s.password == "" {
		log.Fatal("No esta definida la variable de ambiente ", PSWD, " para el password de acceso ZincSearch")
	}

	//obtiene "host" de variable de entorno, o utiliza valor por defecto: localhost
	s.host = os.Getenv(HOST)
	if s.host == "" {
		s.host = defaultHost
	}

	//obtiene "puerto" de variable de entorno, o utiliza valor por defecto: localhost
	var port string = os.Getenv(ZincSearchPort)
	if port == "" {
		s.port = defaultPort
	} else {
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("El valor definido para el puerto debe ser numerico. Valor recibido: ", port)
		}
		s.port = port
	}

	var https string = os.Getenv(ZincSearchHttps)
	s.https = (https == "1" || https == "s" || https == "S")

	s.initDebug()
}

func (s *ZincSearch) initDebug() {
	debugTxt := os.Getenv("ZINC_LOCAL_DEBUG_ENABLED")
	debugEnabled = (debugTxt != "" && (strings.ToLower(debugTxt) == "true" || debugTxt == "1"))
}

func (s *ZincSearch) IsDebug() bool {
	return debugEnabled
}

func (s *ZincSearch) debugReq(request *http.Request) {
	if !debugEnabled {
		return
	}
	data, err := httputil.DumpRequestOut(request, true)
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func (s *ZincSearch) debugRes(response *http.Response) {
	if !debugEnabled {
		return
	}
	data, err := httputil.DumpResponse(response, true)
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
