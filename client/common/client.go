package common

import (
	"net"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
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
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	bet, err := NewBetFromEnv()
	if err != nil {
		log.Errorf("action: create_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		if err := c.sendBet(bet); err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

// sendBet sends the bet over the current connection and waits for the server answer.
// The connection is always closed before returning, even on error.
func (c *Client) sendBet(bet Bet) error {
	defer c.Close()

	if err := SendBet(c.conn, bet); err != nil {
		return err
	}

	answer, err := ReceiveAnswer(c.conn)
	if err != nil {
		return err
	}

	if answer == answerSuccess {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet.Document,
			bet.Number,
		)
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v",
			bet.Document,
			bet.Number,
		)
	}
	return nil
}

// Close Closes the client connection.
func (c *Client) Close() {
	c.conn.Close()
}