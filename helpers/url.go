package helpers

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

const SEP string = "/"
const SEP_QUERY string = "?"

func GetUrl(
	https bool, host string, port string,
	resource string, urlquery string) (url string) {

	//stringbuilder donde se ira formando el URL a utilzar
	var sb strings.Builder

	if https {
		sb.WriteString("https://")
	} else {
		sb.WriteString("http://")
	}
	sb.WriteString(host)
	sb.WriteString(":")
	sb.WriteString(port)

	//Si el recurso no inicia con pleca
	if strings.Index(resource, SEP) > 0 {
		sb.WriteString(SEP)
	}

	sb.WriteString(resource)

	if urlquery != "" {
		if strings.Index(urlquery, SEP_QUERY) > 0 {
			sb.WriteString(SEP_QUERY)
		}
		sb.WriteString(urlquery)
	}

	url = sb.String()
	return url
}

func GetUrlQueryFromStruct(v interface{}) string {
	var sb strings.Builder

	//referencia por reflection a estructura
	vReference := reflect.ValueOf(v)

	//obtiene el tipo de la estructura
	vType := reflect.TypeOf(v)

	if vType.Kind() != reflect.Struct {
		log.Fatal("type ", vType.Kind(), " is not supported")
	}

	//obtiene los campos
	fields := reflect.VisibleFields(vType)

	//recorre los campos para irlos agregando al query
	for _, sf := range fields {
		if sb.Len() == 0 {
			//inicia url query
			sb.WriteString("?")
		} else {
			//inicia url query
			sb.WriteString("&")
		}

		//obtiene referencia a propiedad
		field := vReference.FieldByName(sf.Name)

		//agrega nombre del parametro al query
		sb.WriteString(strings.ToLower(sf.Name))
		sb.WriteString("=")

		//agrega valor del parametro al query
		switch sf.Type.Kind() {
		//Texto
		case reflect.String:
			if field.String() != "" {
				sb.WriteString(field.String())
			}
		//Enteros
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			stringVal := strconv.FormatInt(field.Int(), 10)
			sb.WriteString(stringVal)
		//Enteros sin signo
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			stringVal := strconv.FormatUint(field.Uint(), 10)
			sb.WriteString(stringVal)
		//Bool
		case reflect.Bool:
			stringVal := strconv.FormatBool(field.Bool())
			sb.WriteString(stringVal)

		}
	}
	return sb.String()
}
