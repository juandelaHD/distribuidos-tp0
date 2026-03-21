import socket
import logging
import signal
import threading
import multiprocessing

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
        self.processes = []
        self._manager = multiprocessing.Manager()
        self._store_lock = self._manager.Lock()
        self._winners_dict = self._manager.dict()
        self._barrier = self._manager.Barrier(total_agencies)
        self._lottery_lock = self._manager.Lock()
        self._lottery_done = self._manager.Event()

    def run(self):
        signal.signal(signal.SIGTERM, self.__shutdown)

        while self._running and len(self.processes) < self.total_agencies:
            try:
                client_sock = self.__accept_new_connection()
                p = multiprocessing.Process(
                    target=self.__handle_client_connection,
                    args=(client_sock,)
                )
                p.start()
                client_sock.close()  # Parent closes its copy after fork
                self.processes.append(p)
            except Exception as e:
                if self._running:
                    logging.error(f"action: server_loop | result: fail | error: {e}")
                break

        for p in self.processes:
            p.join()
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
        except threading.BrokenBarrierError:
            logging.info(f"action: client_connection | result: fail | reason: server_shutdown | agency: {agency}")
        except Exception as e:
            logging.error(f"action: apuesta_recibida | result: fail | error: {e}")
            self.__close_client_with_error(client_sock)
        finally:
            self.__close_socket(client_sock)

    def __is_done(self, bets, agency, client_sock):
        if len(bets) != 0:
            return False
        logging.info(f"action: done_received | result: success | agency: {agency}")
        self._barrier.wait()
        self.__run_lottery_once()
        winners = self._winners_dict.get(agency, [])
        send_results(client_sock, winners)
        return True

    def __process_batch(self, client_sock, bets):
        with self._store_lock:
            store_bets(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
        send_answer(client_sock, True)

    def __run_lottery_once(self):
        with self._lottery_lock:
            if not self._lottery_done.is_set():
                self.__run_lottery()
                self._lottery_done.set()
        self._lottery_done.wait()

    def __run_lottery(self):
        logging.info("action: sorteo | result: success")
        bets = load_bets()
        winners_by_agency = {}
        for bet in bets:
            if has_won(bet):
                agency = int(bet.agency)
                winners_by_agency.setdefault(agency, []).append(int(bet.document))
        for agency, docs in winners_by_agency.items():
            self._winners_dict[agency] = docs

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
        logging.info('action: shutdown_server | result: in_progress')
        self._running = False

        try:
            self._barrier.abort()
        except Exception:
            pass

        self.__close_socket(self._server_socket)

        for p in self.processes:
            p.join()

        self._manager.shutdown()
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
