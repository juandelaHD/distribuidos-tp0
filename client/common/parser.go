package common

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	AGENCY_BITS   = 8
	DOCUMENT_BITS = 32
	NUMBER_BITS   = 16
	DECIMAL_BASE  = 10
	BET_FIELDS_COUNT  = 5
	IDX_FIRST_NAME    = 0
	IDX_LAST_NAME     = 1
	IDX_DOCUMENT      = 2
	IDX_BIRTHDATE     = 3
	IDX_BET_NUMBER    = 4
)

// NextBatch reads up to batchSize bets from the CSV. Returns an empty slice when finished.
// Malformed rows are skipped with a warning log.
func (r *BetReader) NextBatch() ([]Bet, error) {
	if r.Finished {
		return nil, nil
	}
	var batch []Bet
	for len(batch) < r.BatchSize {
		record, err := r.Reader.Read()
		if err == io.EOF {
			r.Finished = true
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV: %w", err)
		}
		if len(record) < BET_FIELDS_COUNT {
			log.Warningf("action: parse_bet | result: skip | row: %v | error: expected %d fields, got %d", record, BET_FIELDS_COUNT, len(record))
			continue
		}
		bet, err := parseBetRecord(record)
		if err != nil {
			log.Warningf("action: parse_bet | result: skip | row: %v | error: %v", record, err)
			continue
		}
		batch = append(batch, bet)
	}
	return batch, nil
}

// Close releases the underlying file handle.
func (r *BetReader) Close() error {
	return r.File.Close()
}

// NewBetFromEnv builds a Bet from environment variables.
func NewBetFromEnv() (Bet, error) {
	firstName, err := readRequiredEnv("FIRST_NAME")
	if err != nil {
		return Bet{}, err
	}
	lastName, err := readRequiredEnv("LAST_NAME")
	if err != nil {
		return Bet{}, err
	}
	document, err := parseEnvUint("DOCUMENT", DOCUMENT_BITS)
	if err != nil {
		return Bet{}, err
	}
	birthDate, err := readRequiredEnv("BIRTHDATE")
	if err != nil {
		return Bet{}, err
	}
	if _, err := time.Parse(birthdateLayout, birthDate); err != nil {
		return Bet{}, fmt.Errorf("invalid BIRTHDATE %q: expected YYYY-MM-DD with a valid calendar date", birthDate)
	}
	number, err := parseEnvUint("NUMBER", NUMBER_BITS)
	if err != nil {
		return Bet{}, err
	}
	return NewBet(firstName, lastName, uint32(document), birthDate, uint16(number)), nil
}

// ParseAgency converts a client ID string to a uint8 agency ID.
func ParseAgency(id string) (uint8, error) {
	v, err := strconv.ParseUint(id, DECIMAL_BASE, AGENCY_BITS)
	if err != nil {
		return 0, fmt.Errorf("client ID %q is not a valid agency (expected 0-255): %w", id, err)
	}
	return uint8(v), nil
}

func parseBetRecord(record []string) (Bet, error) {
	firstName := record[IDX_FIRST_NAME]
	lastName := record[IDX_LAST_NAME]

	document, err := strconv.ParseUint(record[IDX_DOCUMENT], DECIMAL_BASE, DOCUMENT_BITS)
	if err != nil {
		return Bet{}, fmt.Errorf("invalid document %q: %w", record[IDX_DOCUMENT], err)
	}

	birthDate := record[IDX_BIRTHDATE]
	if _, err := time.Parse(birthdateLayout, birthDate); err != nil {
		return Bet{}, fmt.Errorf("invalid birthdate %q: expected YYYY-MM-DD", birthDate)
	}

	number, err := strconv.ParseUint(record[IDX_BET_NUMBER], DECIMAL_BASE, NUMBER_BITS)
	if err != nil {
		return Bet{}, fmt.Errorf("invalid number %q: %w", record[IDX_BET_NUMBER], err)
	}

	return NewBet(firstName, lastName, uint32(document), birthDate, uint16(number)), nil
}

func readRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing required env var %s", key)
	}
	return value, nil
}

func parseEnvUint(key string, bitSize int) (uint64, error) {
	raw, err := readRequiredEnv(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(raw, DECIMAL_BASE, bitSize)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", key, err)
	}
	return value, nil
}
