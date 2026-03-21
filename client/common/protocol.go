package common

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	AGENCY_SIZE         = 1
	STRING_LENGTH_SIZE  = 1
	STRING_MAX_SIZE     = 255
	DOCUMENT_SIZE       = 4
	BIRTHDATE_SIZE      = 4
	NUMBER_SIZE         = 2
	ANSWER_SIZE         = 1
	NUMBER_OF_BETS_SIZE = 2
)

const birthdateLayout = "2006-01-02" // Go's reference time for parsing dates

// MAX_CHUNK_SIZE is the maximum number of bytes sent in a single write call.
// It limits individual socket writes to avoid large OS-level TCP segmentation.
const MAX_CHUNK_SIZE = 8 * 1024

const (
	answerSuccess = 0
	// answerFailure = 1
)

// betSerializedSize returns the number of bytes appendBet will write for a given bet.
func betSerializedSize(bet Bet) int {
	return STRING_LENGTH_SIZE + len(bet.FirstName) +
		STRING_LENGTH_SIZE + len(bet.LastName) +
		DOCUMENT_SIZE + BIRTHDATE_SIZE + NUMBER_SIZE
}

// SendBatch sends a batch of Bets as one logical message over conn.
// Wire format: N_BETS(2, big-endian) | AGENCY(1) | BET1 | BET2 | ... | BETn
// BET encoding: len_fn(1) | fn(N) | len_ln(1) | ln(M) | document(4, big-endian) | birthdate(4) | number(2, big-endian)
// The header is sent first, then bets are flushed in chunks of at most MAX_CHUNK_SIZE bytes.
// The server reads exactly N_BETS bets and sends one ACK for the whole batch.
func SendBatch(conn net.Conn, agency uint8, bets []Bet) error {
	// Send header: N_BETS(2) | AGENCY(1)
	header := make([]byte, NUMBER_OF_BETS_SIZE+AGENCY_SIZE)
	binary.BigEndian.PutUint16(header[0:NUMBER_OF_BETS_SIZE], uint16(len(bets)))
	header[NUMBER_OF_BETS_SIZE] = agency
	if err := sendAll(conn, header); err != nil {
		return err
	}

	// Serialize bets and flush to conn whenever the buffer reaches MAX_CHUNK_SIZE.
	buf := make([]byte, 0, MAX_CHUNK_SIZE)
	for i, bet := range bets {
		if len(buf)+betSerializedSize(bet) > MAX_CHUNK_SIZE && len(buf) > 0 {
			if err := sendAll(conn, buf); err != nil {
				return err
			}
			buf = buf[:0]
		}
		var err error
		buf, err = appendBet(buf, i, bet)
		if err != nil {
			return err
		}
	}
	if len(buf) > 0 {
		return sendAll(conn, buf)
	}
	return nil
}

// appendBet serializes a single Bet into buf and returns the extended slice.
func appendBet(buf []byte, idx int, bet Bet) ([]byte, error) {
	var err error
	buf, err = appendString(buf, bet.FirstName)
	if err != nil {
		return nil, fmt.Errorf("bet[%d] first_name: %w", idx, err)
	}
	buf, err = appendString(buf, bet.LastName)
	if err != nil {
		return nil, fmt.Errorf("bet[%d] last_name: %w", idx, err)
	}

	// Document (4 bytes, big-endian)
	docBytes := make([]byte, DOCUMENT_SIZE)
	binary.BigEndian.PutUint32(docBytes, bet.Document)
	buf = append(buf, docBytes...)

	// Birthdate (4 bytes): year(2, big-endian) | month(1) | day(1)
	t, err := time.Parse(birthdateLayout, bet.Birthdate)
	if err != nil {
		return nil, fmt.Errorf("bet[%d] birthdate %q is not a valid date (expected YYYY-MM-DD): %w", idx, bet.Birthdate, err)
	}
	bdBytes := make([]byte, BIRTHDATE_SIZE)
	binary.BigEndian.PutUint16(bdBytes[0:2], uint16(t.Year()))
	bdBytes[2] = byte(t.Month())
	bdBytes[3] = byte(t.Day())
	buf = append(buf, bdBytes...)

	// Number (2 bytes, big-endian)
	numBytes := make([]byte, NUMBER_SIZE)
	binary.BigEndian.PutUint16(numBytes, bet.Number)
	buf = append(buf, numBytes...)

	return buf, nil
}

// ReceiveAnswer reads the server's confirmation byte.
func ReceiveAnswer(conn net.Conn) (uint8, error) {
	buf, err := recvAll(conn, ANSWER_SIZE)
	if err != nil {
		return 0, err
	}
	if len(buf) != ANSWER_SIZE {
		return 0, fmt.Errorf("expected %d bytes for answer, got %d", ANSWER_SIZE, len(buf))
	}
	return buf[0], nil
}

// appendString encodes strings as length-prefixed bytes (1-byte length + content).
func appendString(buf []byte, s string) ([]byte, error) {
	if len(s) > STRING_MAX_SIZE {
		return nil, fmt.Errorf("string exceeds max length of %d bytes", STRING_MAX_SIZE)
	}
	buf = append(buf, byte(len(s)))
	buf = append(buf, []byte(s)...)
	return buf, nil
}
