import socket
import logging
import signal

from common.protocol import recv_batch, send_answer
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
        Read batches of bets from a client socket and respond to each one.
        Loops until the client closes the connection.
        Each batch is processed atomically: all bets are stored or none are.
        """
        try:
            while True:
                bets = recv_batch(client_sock)
                self.__process_batch(client_sock, bets)
        except OSError:
            pass  # client disconnected normally
        except ValueError as e:
            logging.error(f'action: apuesta_recibida | result: fail | error: {e}')
            try:
                send_answer(client_sock, False)
            except OSError:
                pass
        finally:
            client_sock.close()
            self._client_sock = None

    def __process_batch(self, client_sock, bets):
        n = len(bets)
        try:
            store_bets(bets)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {n}')
            send_answer(client_sock, True)
        except Exception:
            logging.error(f'action: apuesta_recibida | result: fail | cantidad: {n}')
            send_answer(client_sock, False)

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