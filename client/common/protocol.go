package common

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	AGENCY_SIZE        = 1
	STRING_LENGTH_SIZE = 1
	STRING_MAX_SIZE    = 255
	DOCUMENT_SIZE      = 4
	BIRTHDATE_SIZE     = 4
	NUMBER_SIZE        = 2
	ANSWER_SIZE        = 1
)

const birthdateLayout = "2006-01-02" // Go's reference time for parsing dates

const (
	answerSuccess = 0
	// answerFailure = 1
)

// SendBet serializes a Bet and sends it over the provided connection.
// Bet encoding: agency(1) | len_first_name(1) | first_name(N) | len_last_name(1) | last_name(M) | document(4) | birthdate(4) | number(2).
// birthdate encoding: year(2, big-endian) | month(1) | day(1).
func SendBet(conn net.Conn, bet Bet) error {
	buf := make([]byte, 0, AGENCY_SIZE+STRING_LENGTH_SIZE+len(bet.FirstName)+STRING_LENGTH_SIZE+len(bet.LastName)+DOCUMENT_SIZE+BIRTHDATE_SIZE+NUMBER_SIZE)

	// Agency (1 byte)
	buf = append(buf, bet.Agency)

	// First name and last name (length-prefixed strings)
	var err error
	buf, err = appendString(buf, bet.FirstName)
	if err != nil {
		return fmt.Errorf("first_name: %w", err)
	}
	buf, err = appendString(buf, bet.LastName)
	if err != nil {
		return fmt.Errorf("last_name: %w", err)
	}

	// Document (4 bytes, big-endian)
	docBytes := make([]byte, DOCUMENT_SIZE)
	binary.BigEndian.PutUint32(docBytes, bet.Document)
	buf = append(buf, docBytes...)

	// Birthdate (4 bytes): year(2, big-endian) | month(1) | day(1)
	t, err := time.Parse(birthdateLayout, bet.Birthdate)
	if err != nil {
		return fmt.Errorf("birthdate %q is not a valid date (expected YYYY-MM-DD): %w", bet.Birthdate, err)
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

	return sendAll(conn, buf)
}

// ReceiveAnswer reads the server's confirmation byte, returning true on success.
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
