import socket
import logging
import signal

from common.protocol import recv_batch, send_answer, send_results
from common.utils import store_bets, load_bets, has_won

class Server:
    def __init__(self, port, listen_backlog, total_agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self.total_agencies = total_agencies
        self.finished_clients = {}  # agency (int) -> socket

    def run(self):
        signal.signal(signal.SIGTERM, self.__shutdown)

        while self._running:
            try:
                if len(self.finished_clients) == self.total_agencies:
                    self.__run_lottery()
                    self.__shutdown()
                else:
                    client_sock = self.__accept_new_connection()
                    if client_sock:
                        self.__handle_client_connection(client_sock)
            except Exception as e:
                if self._running:
                    logging.error(f"action: server_loop | result: fail | error: {e}")
                return

        logging.info('action: server_shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        agency = None
        try:
            while True:
                bets, recv_agency = recv_batch(client_sock)
                if agency is None:
                    agency = recv_agency
                if self.__is_done(bets, agency, client_sock):
                    return
                self.__process_batch(client_sock, bets)
        except OSError as e:
            logging.info(f"action: client_connection | result: fail | error: {e}")
            self.__close_socket(client_sock)
        except Exception as e:
            logging.error(f"action: apuesta_recibida | result: fail | error: {e}")
            self.__close_client_with_error(client_sock)

    def __is_done(self, bets, agency, client_sock):
        if len(bets) != 0:
            return False
        self.finished_clients[agency] = client_sock
        logging.info(f"action: done_received | result: success | agency: {agency}")
        return True

    def __process_batch(self, client_sock, bets):
        store_bets(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
        send_answer(client_sock, True)


    def __run_lottery(self):
        logging.info("action: sorteo | result: success")
        bets = load_bets()
        winners = [(bet.agency, int(bet.document)) for bet in bets if has_won(bet)]
        send_results(self.finished_clients, winners)

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

    def __shutdown(self):
        self._running = False
        for sock in self.finished_clients.values():
            self.__close_socket(sock)
        self.__close_socket(self._server_socket)
        logging.info('action: shutdown_server | result: success')

    def __close_socket(self, sock):
        try:
            sock.close()
        except OSError:
            pass # Socket already closed, ignore

    def __close_client_with_error(self, client_sock):
        try:
            send_answer(client_sock, False)
        except OSError:
            pass # Socket already closed to send error message, ignore
        self.__close_socket(client_sock)
