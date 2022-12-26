package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"

	"zincsearch.com/mailindex/api/override/godotenv"
	"zincsearch.com/mailindex/api/service"
)

var api service.ZincSearch
var profiling bool = false
var createMainIndex bool = false
var queueMsgQuantity int = 0

var archivos int = 0
var carpetas int = 0

// cola utilizada para acumular documentos por enviar
var queue chan string = make(chan string)

// Crea grupos de espera para trabajos en paralelo
var wg sync.WaitGroup

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error en la carga de variables de ambiente: %s", err)
	}

	txtProf := os.Getenv("ZINC_LOCAL_PROFILING_ENABLED")
	profiling = (txtProf != "" && (strings.ToLower(txtProf) == "true" || txtProf == "1"))
	fmt.Println("profiling: ", txtProf)

	createIndex := os.Getenv("ZINC_LOCAL_CREATE_MAIN_INDEX")
	createMainIndex = createIndex != "" && (strings.ToLower(createIndex) == "true" || createIndex == "1")
	fmt.Println("create main index: ", createIndex)

}

func main() {
	//profiling config
	if profiling {
		//crea archivo de salida
		f, err := os.Create(`cpu.pprof`)

		if err != nil {
			log.Fatal(err)
		}

		//inicia perfil
		err = pprof.StartCPUProfile(f)

		if err != nil {
			log.Fatal("Error al iniciar el perfilamiento: ", err)
		}

		defer f.Close()

		defer pprof.StopCPUProfile()
	}

	fmt.Println("Inicia", time.Now().Format(time.RFC1123))

	if len(os.Args) == 1 {
		log.Fatal("Es obligatorio ingresar la ruta del directorio. Ej. C:\\enron_mail_20110402")
	}

	var dirname = os.Args[1]

	//inicializa servicio, cargando su configuracion del archivo .env
	api.Inicia()
	verificaIndice()

	wg.Add(1)
	go importaArchivos(dirname)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(queue)
	}(&wg)

	enviarDocsAZincSearch()

	fmt.Println(" Folders procesados: ", carpetas)
	fmt.Println(" Archivos procesados: ", archivos)
	fmt.Println(" Mensajes procesados: ", queueMsgQuantity)

	fmt.Println("Termina", time.Now().Format(time.RFC1123))
}

// Veririca la existencia de indice, y en caso de no existir lo crea
func verificaIndice() {
	if !createMainIndex {
		return
	}
	//TODO remove, only for testing
	//api.DeleteIndex(service.INDEX_NAME)

	//busca existencia de indice
	resultadoCreacion, _ := api.ExistsIndex(service.INDEX_NAME)
	if !resultadoCreacion {
		crearIndice()
	}

	_, httpError := api.ExistsIndex(service.INDEX_NAME)
	if httpError.Code != 0 {
		log.Fatal("No se creó el indice, no se puede continuar")
	}
}

// funcion llamada para iniciar procesamiento de archivos de correo
func importaArchivos(dir string) (folders int, documents int) {
	defer wg.Done()

	//obtiene listado de archivos en directorio indicado
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	//almacenara subfolders encontrados
	var subfolders []string

	//recorrer listado de archivos y subfolders
	for _, file := range files {

		if file.IsDir() {

			//utiliza separador del sistema operativo
			path := filepath.Join(dir, file.Name())

			//acumula cantidad de carpetas
			carpetas++

			//almacena lista de subfolders en slice para su procesamiento posterior
			subfolders = append(subfolders, path)

		} else {
			archivos++

			filename := filepath.Join(dir, file.Name())
			// procesa en paralelo archivos
			procesaArchivo(filename)

		}
	}

	if len(subfolders) > 0 {
		for _, p := range subfolders {
			wg.Add(1)
			go importaArchivos(p)
		}
	}

	return carpetas, archivos
}

// Procesa envío de archivo individual
func procesaArchivo(filename string) {

	dat, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	//parseo de texto a estructura email (headers/body)
	msg, err := mail.ReadMessage(bytes.NewBuffer(dat))

	//Si no se pudo obtener la estructura del mail, se omite el registro
	if err != nil {
		return
	}

	stmail, err := parsearDatosEmail(msg)

	if err != nil {
		log.Fatal(err)
	}

	//envia JSON a canal
	queue <- stmail
}

// crea indice como primer paso del proceso (cuando no existe)
func crearIndice() {

	if !createMainIndex {
		return
	}

	jsonIndex, err := os.Open("json/index_mailindex.json")
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(jsonIndex)
	if err != nil {
		log.Fatal(err)
	}

	result, errorHttp := api.SaveIndex(service.INDEX_NAME, string(content))

	if result == "" {
		log.Fatal("Error en creación de indice: ", errorHttp)
	}

}

func enviarDocsAZincSearch() {
	const MAX_POR_LOTE int = 5000 //debe ser multiplo de 1000
	const COMMA byte = ','
	const REQUEST_BEGIN string = "{\"index\": \"" + service.INDEX_NAME + "\",\"records\": ["
	var sb strings.Builder
	var i int = 0

	for strJson := range queue {
		if i > 0 {
			sb.WriteByte(COMMA)
		} else {
			//prepara peticion
			sb.WriteString(REQUEST_BEGIN)
		}

		//aumenta contador
		i++

		//agrega json actual a peticion
		sb.WriteString(strJson)

		queueMsgQuantity++
		if queueMsgQuantity%MAX_POR_LOTE == 0 {
			enviarDocs(&sb)

			//reinicia datos para siguiente bloque
			i = 0
		}
	}

	//el ultimo lote puede no haber alcanzado el tamaño maximo
	//por lo que se procesa si hay al menos un registro incluido
	if i > 0 {
		enviarDocs(&sb)
	}
}

func enviarDocs(sb *strings.Builder) {

	const REQUEST_END string = "]}"
	//cierra estructura JSON
	sb.WriteString(REQUEST_END)
	defer sb.Reset()

	result, errorHttp := api.CreateDocumentBulk(sb.String())

	if result == "" {
		log.Fatal("Error en creación de indice: ", errorHttp)
	}
}

func parsearDatosEmail(info *mail.Message) (emailJson string, err error) {
	email := stEmail{}
	email.Bcc = info.Header.Get("Bcc")
	email.Cc = info.Header.Get("Cc")
	email.ContentType = info.Header.Get("Content-Type")

	//email.Date = info.Header.Get("Date")
	//cambiar campo texto a tipo time.Time
	date, err := mail.ParseDate(info.Header.Get("Date"))

	if err != nil {
		log.Fatal(err)
	}

	//formatea fecha con formato por defecto de ZincSearch
	email.Date = date.Format("2006-01-02T15:04:05Z07:00")
	email.From = info.Header.Get("From")
	email.MessageID = info.Header.Get("MessageID")
	email.ReplyTo = info.Header.Get("Reply-To")
	email.Sender = info.Header.Get("Sender")
	email.Subject = info.Header.Get("Subject")
	email.To = info.Header.Get("To")
	txtBytes, _ := ioutil.ReadAll(info.Body)
	email.TextBody = strings.TrimSuffix(string(txtBytes[:]), "\n")

	jsonBytes, err := json.Marshal(email)
	emailJson = string(jsonBytes)

	return emailJson, err
}

// Estructura de email
type stEmail struct {
	Subject string
	Sender  string
	From    string
	ReplyTo string
	To      string
	Cc      string
	Bcc     string
	Date    string

	MessageID   string
	ContentType string

	TextBody string
}
