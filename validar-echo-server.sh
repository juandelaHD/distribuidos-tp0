#!/bin/bash

NETWORK="tp0_testing_net"
SERVER_CONTAINER="server"
PORT=12345
MSG="test message"
COMMAND="echo $MSG | nc $SERVER_CONTAINER $PORT"

RESPONSE=$(docker run --rm --network "$NETWORK" alpine sh -c "$COMMAND")

if [ "$RESPONSE" = "$MSG" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
