import socket
import logging
import signal

from common.protocol import recv_bet, send_answer
from common.utils import store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._client_sock = None

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        
        signal.signal(signal.SIGTERM, self.__signal_handler)

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as e:
                if self._running:
                    logging.error(f"action: server_loop | result: fail | error: {e}")
                return
            except Exception as e:
                logging.error(f"action: server_loop | result: fail | error: {e}")
                return
            
        logging.info('action: shutdown_server | result: success')

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet = recv_bet(client_sock)
            store_bets([bet])
            logging.info(
                f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}'
            )
            send_answer(client_sock, True)
        except OSError as e:
            logging.error(f"action: apuesta_almacenada | result: fail | error: {e}")
            send_answer(client_sock, False)
        finally:
            client_sock.close()
            self._client_sock = None

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __signal_handler(self, signum, frame):
        """
        Handle SIGTERM signal

        When SIGTERM signal is received, the server socket is shutdown and closed.
        """
        logging.info('action: shutdown_server | result: in_progress')
        self._running = False
        if self._client_sock:
            self._client_sock.close()
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()