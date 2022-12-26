Codigo de ejemplo para realizar carga de archivos a ZincSearch, utilizando Go. 
El proceso utiliza go routines y channels para realizar la carga de miles de registros en paralelo, y utiliza carga bulk de ZincSearch para dimsinuir cantidad de llamados de su API

## Configuración
Previo a ejecutar el proceso, se necesita crear un archivo .env con las siguientes propiedades:
- ZINC_FIRST_ADMIN_USER=
- ZINC_FIRST_ADMIN_PASSWORD=
- ZINC_LOCAL_CREATE_MAIN_INDEX=false
- ZINC_LOCAL_DEBUG_ENABLED=false
- ZINC_LOCAL_PROFILING_ENABLED=false

Donde:
- ZINC_FIRST_ADMIN_USER: es el usuario para acceder al API de ZincSearch
- ZINC_FIRST_ADMIN_PASSWORD: contraseña de usuario de API ZincSearch
- ZINC_LOCAL_CREATE_MAIN_INDEX: boolean (true/false) que indica si el indice se creara previo a la carga de lso datos (opcional)
- ZINC_LOCAL_DEBUG_ENABLED: boolean (true/false) habilita mensajes de consoola en modo debug
- ZINC_LOCAL_PROFILING_ENABLED: boolean (true/false) habilita profiling de carga


## Ejecución
Ejemplo de llamado:

    go run indexer.go [DIRECTORIO]

Donde:
- DIRECTORIO: es la carpeta donde se encuentra lso archivos que seran cargados a la instancia destino de ZincSearch


