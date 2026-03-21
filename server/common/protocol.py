import struct

from common.parser import parse_bet

# Protocol constants
AGENCY_SIZE = 1
STRING_LENGTH_SIZE = 1
# STRING_MAX_SIZE = 255
DOCUMENT_SIZE = 4
BIRTHDATE_SIZE = 4
NUMBER_SIZE = 2
NUMBER_OF_BETS_SIZE = 2
ANSWER_SIZE = 1

# Answer codes
ANSWER_SUCCESS = 0
ANSWER_FAIL = 1


"""
Receive a batch of Bets from the socket.
Packet encoding: number_of_bets(2, big-endian) | agency(1) | BET1 | BET2 | ...
BET encoding: len_fn(1) | first_name(N) | len_ln(1) | last_name(M) | document(4, big-endian) | birthdate(4) | number(2, big-endian)

Returns a list of Bets.
Raises OSError if the connection is closed or a read fails.
Raises ValueError if any bet contains invalid field values.
"""
def recv_batch(socket):
    n_bets = struct.unpack('!H', _recv_all(socket, NUMBER_OF_BETS_SIZE))[0]
    agency = struct.unpack('!B', _recv_all(socket, AGENCY_SIZE))[0]

    bets = []
    for _ in range(n_bets):
        first_name = _recv_string(socket)
        last_name = _recv_string(socket)
        document = struct.unpack('!I', _recv_all(socket, DOCUMENT_SIZE))[0]
        bd_bytes = _recv_all(socket, BIRTHDATE_SIZE)
        number = struct.unpack('!H', _recv_all(socket, NUMBER_SIZE))[0]
        bets.append(parse_bet(agency, first_name, last_name, document, bd_bytes, number))

    return bets


"""Send a 1-byte result code to the client."""
def send_answer(socket, success):
    code = ANSWER_SUCCESS if success else ANSWER_FAIL
    _send_all(socket, struct.pack('!B', code))


"""Receive a length-prefixed string (1-byte length + content)."""
def _recv_string(socket):
    length = struct.unpack('!B', _recv_all(socket, STRING_LENGTH_SIZE))[0]
    return _recv_all(socket, length).decode('utf-8')


"""Send all bytes to socket, retrying on partial writes."""
def _send_all(socket, data):
    total = 0
    while total < len(data):
        sent = socket.send(data[total:])
        if sent == 0:
            raise OSError("Connection closed during send")
        total += sent


"""Read exactly n bytes from socket, retrying on partial reads."""
def _recv_all(socket, n):
    buf = b''
    while len(buf) < n:
        chunk = socket.recv(n - len(buf))
        if not chunk:
            raise OSError("Connection closed during recv")
        buf += chunk
    return buf
