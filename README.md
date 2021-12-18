# INSTALACIÓN Y DESPLIEGUE

Instalar go es un poco movida, mejor buscar en internet. Luego tienes que actualizar el path para que te reconozca el comando.

```
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/$USER/go (cambiar por la ruta de win2)
export GOBIN=/usr/local/go/bin
export GOROOT=/usr/local/go
```

Cuando esté instalado tienes que crear esta estructura de carpeta 'go/src/github.com/drg809/apibase' para que reconozca bien los módulos y demas movidas.

Descargas ahí el repo y haces 'go get' en el terminal para descargar los módulos (npm i en go) y 'go mod tigy' para que queden listos.

Lo bonito de go es que puedes compilar un binario y ejecutar la api desde el ejecutando 'go build' en consola (estando en la raíz del proyecto) y luego ejecutar el sh obtenido.

También puedes ejecutar 'go run main.go' para levantar el servicio más rápido.

Para ponerlo en prod es un poco más movida porque hay que dar de alta un servicio y ponerlo apuntando al binario obtenido del build.