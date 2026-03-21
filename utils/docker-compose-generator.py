import sys
import os

COMPOSE_NAME ="""name: tp0
"""

COMPOSE_SERVER = """
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""

CLIENT_TEMPLATE = """
  client{n}:
    container_name: client{n}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={n}
      - AGENCY={n}
      - FIRST_NAME=Santiago Lionel
      - LAST_NAME=Lorca
      - DOCUMENT=3090446{n}
      - BIRTHDATE=1999-03-17
      - NUMBER=7574
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data:/data
"""

COMPOSE_NETWORK = """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

def build_compose_yaml(clients):
    yaml = COMPOSE_NAME + COMPOSE_SERVER
    for i in range(1, clients + 1):
        yaml += CLIENT_TEMPLATE.format(n=i)
    yaml += COMPOSE_NETWORK
    return yaml

def write_file(path, content):
	with open(path, "w") as f:
		f.write(content)

def parse_output_path(path):
    if not path.strip():
        raise ValueError("output_file cannot be empty")

    if not path.endswith((".yaml", ".yml")):
        raise ValueError("output_file must end with .yaml or .yml")

    if os.path.isdir(path):
        raise ValueError("output_file cannot be a directory")

    parent = os.path.dirname(path) or "."

    if not os.access(parent, os.W_OK):
        raise ValueError("Cannot write to target directory")

    return path


def parse_clients(num_clients):
    try:
        clients = int(num_clients)
    except ValueError:
        raise ValueError("<num_clients> must be an integer")

    if clients < 0:
        raise ValueError("<num_clients> must be grater than or equal to 0")

    return clients


def parse_inputs(argv):
    if len(argv) != 3:
        raise ValueError(
            "Run: ./generar-compose.sh <output_file.yaml> <num_clients>"
        )

    output_file = parse_output_path(argv[1])
    clients = parse_clients(argv[2])

    return output_file, clients

def main():
    try:
        output_file, clients = parse_inputs(sys.argv)
    except ValueError as e:
        print(f"Error: {e}")
        return 1

    compose_yaml = build_compose_yaml(clients)
    write_file(output_file, compose_yaml)
    
    return 0

if __name__ == "__main__":
	raise SystemExit(main())

