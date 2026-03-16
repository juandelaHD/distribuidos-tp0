# TP0: Docker + Comunicaciones + Concurrencia - Resolución

Este documento presenta la resolución de los ejercicios propuestos en el TP0 de la materia Sistemas Distribuidos de la Facultad de Ingeniería de la Universidad de Buenos Aires. Cada ejercicio se aborda de manera individual, con un enlace a las instrucciones detalladas y al commit específico que contiene la solución implementada.

> Puedes encontrar la consigna completa haciendo click [acá](README.md).

## Ejercicio 1

> **Instrucciones:** [README.md](README.md#ejercicio-n1)

Se incorporó un **generador automático de configuraciones de Docker Compose**. Este generador permite crear dinámicamente un archivo `docker-compose` con:

* Un **server**
* **N instancias de clientes** (`client1`, `client2`, ..., `clientN`)

### Uso del generador de Docker Compose

1. Dar permisos de ejecución al script:

```bash
chmod +x generar-compose.sh
```

2. Generar el archivo `docker-compose` indicando:

- El **nombre del archivo de salida**
- La **cantidad de clientes a generar**

```bash
./generar-compose.sh <nombre_archivo> <cantidad_clientes>
```

3. Levantar los contenedores con Docker Compose:

```bash
docker compose -f <nombre_archivo> up
```

## Ejercicio 2

> **Instrucciones:** [README.md](README.md#ejercicio-n2)

Se montaron los archivos de configuración del servidor y del cliente (server/config.ini y client/config.yaml) como volúmenes en los contenedores mediante [bind mounts](https://docs.docker.com/engine/storage/bind-mounts/).

De esta forma, los contenedores acceden directamente a los archivos de configuración ubicados en el sistema de archivos del host en tiempo de ejecución. Esto permite modificar la configuración sin reconstruir las imágenes, ya que los cambios en los archivos montados se reflejan inmediatamente dentro de los contenedores.

## Ejercicio 3

> **Instrucciones:** [README.md](README.md#ejercicio-n3)

Se implementó el script de validación del "echo server" para verificar que el servidor responde correctamente a los mensajes enviados por los clientes. El script utiliza `nc` (netcat) dentro de un contenedor [Alpine](https://hub.docker.com/_/alpine) (una imagen ligera de Linux) para enviar un mensaje al servidor y verificar que la respuesta sea la esperada.

Cómo ejecutarlo:

1. Dar permisos al script:

```
chmod +x validar-echo-server.sh
```

2. Ejecutar la validación desde la terminal:

```
./validar-echo-server.sh
```

Resultado esperado:

```
action: test_echo_server | result: success
```

> Nota: Asegúrese de que los containers estén corriendo antes de ejecutar la validación.

## Ejercicio 4

> **Instrucciones:** [README.md](README.md#ejercicio-n4)

Para implementar el *graceful shutdown* primero se identificaron los recursos que deben cerrarse de forma ordenada: el socket del servidor y las conexiones activas del cliente.

En el servidor, se agregó un handler para SIGTERM que invoca shutdown() y close() sobre el socket del servidor y registra el evento de apagado.
En el cliente, se incorporó la captura de SIGTERM, que ejecuta Close() sobre la conexión, registra el shutdown y finaliza el proceso.

De esta manera, tanto el servidor como los clientes liberan correctamente los sockets antes de terminar, evitando dejar recursos abiertos o conexiones colgantes.

## Ejercicio 5

> **Instrucciones:** [README.md](README.md#ejercicio-n5)

## Ejercicio 6

> **Instrucciones:** [README.md](README.md#ejercicio-n6)

## Ejercicio 7

> **Instrucciones:** [README.md](README.md#ejercicio-n7)

## Ejercicio 8

> **Instrucciones:** [README.md](README.md#ejercicio-n8)
