package common

import (
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

// StartClientLoop reads bets from the agency CSV file and sends them to the server in batches
// over a single persistent TCP connection. After all batches are sent, notifies the server
// with a done signal and waits for the lottery winners list.
func (c *Client) StartClientLoop() {
	agency, err := ParseAgency(c.config.ID)
	if err != nil {
		log.Errorf("action: parse_agency | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	reader, err := NewBetReader(c.config.ID, c.config.BatchMaxAmount)
	if err != nil {
		log.Errorf("action: open_csv | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}
	defer reader.Close()

	if err := c.createClientSocket(); err != nil {
		return
	}
	defer c.Close()

	for {
		batch, err := reader.NextBatch()
		if err != nil {
			log.Errorf("action: read_batch | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}
		if len(batch) == 0 {
			break
		}

		if err := SendBatch(c.conn, agency, batch); err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}

		answer, err := ReceiveAnswer(c.conn)
		if err != nil {
			log.Errorf("action: receive_answer | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}

		if answer == answerSuccess {
			log.Infof("action: apuesta_enviada | result: success")
		} else {
			log.Errorf("action: apuesta_enviada | result: fail")
		}
	}

	if err := SendDone(c.conn, agency); err != nil {
		log.Errorf("action: notify_done | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	winners, err := ReceiveWinners(c.conn)
	if err != nil {
		log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
}

// Close Closes the client connection.
func (c *Client) Close() {
	c.conn.Close()
}