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

Para el Ejercicio 5 reorganicé la solución pensando primero en qué datos viajan en cada apuesta y en qué orden debían transmitirse para que cliente y servidor interpretaran la misma estructura, y a partir de eso definí un protocolo binario explícito con campos de tamaño fijo y variable delimitados, incluyendo reglas de longitud y validaciones de formato para evitar ambigüedades. 

De esta forma, cada apuesta se codifica como la secuencia: agency(1 byte) | len_first_name(1 byte) | first_name(N bytes) | len_last_name(1 byte) | last_name(M bytes) | document(4 bytes) | birthdate(4 bytes) | number(2 bytes), donde los campos de longitud variable se interpretan usando sus prefijos de tamaño (N = valor de len_first_name y M = valor de len_last_name). El campo birthdate se serializa en formato binario como year(2 bytes, big-endian) | month(1 byte) | day(1 byte). En lugar de usar texto ASCII, lo trabajé directo con el formato binario, lo que me permitió optimizar un poco más el tamaño del payload con un formato más compacto y estructurado.

Luego implementé la serialización y deserialización de ambos lados, de modo que cada mensaje se transforme de datos de la lógica de negocio a bytes y nuevamente a datos sin pérdida de información.

Busqué separar el modelo de dominio y la lógica de negocio de la capa de transporte, manteniendo la construcción/validación de datos en archivos de dominio y dejando en la capa de comunicación únicamente el empaquetado, envío, recepción e interpretación del protocolo. Finalmente, para el uso correcto de sockets, contemplé lecturas y escrituras parciales con un esquema de envío/recibo total de bytes esperados, manejo de errores de conexión y validaciones de integridad del payload, garantizando que la confirmación de recepción sea consistente con el estado real del procesamiento.

## Ejercicio 6

> **Instrucciones:** [README.md](README.md#ejercicio-n6)

## Ejercicio 7

> **Instrucciones:** [README.md](README.md#ejercicio-n7)

## Ejercicio 8

> **Instrucciones:** [README.md](README.md#ejercicio-n8)
