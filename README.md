Codigo para realizar carga de archivos a ZincSearch, utilizando Go. 
El proceso utiliza go routines y channels para realizar la carga de miles de registros en paralelo, y utiliza carga bulk de ZincSearch para dimsinuir cantidad de llamados de su API