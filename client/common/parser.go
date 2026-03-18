package common

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	AGENCY_BITS   = 8
	DOCUMENT_BITS = 32
	NUMBER_BITS   = 16
	DECIMAL_BASE  = 10
)

func NewBetFromEnv() (Bet, error) {
	agency, err := parseEnvUint("AGENCY", AGENCY_BITS)
	if err != nil {
		return Bet{}, err
	}
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
	birthdate, err := readRequiredEnv("BIRTHDATE")
	if err != nil {
		return Bet{}, err
	}
	if _, err := time.Parse(birthdateLayout, birthdate); err != nil {
		return Bet{}, fmt.Errorf("invalid BIRTHDATE %q: expected YYYY-MM-DD with a valid calendar date", birthdate)
	}
	number, err := parseEnvUint("NUMBER", NUMBER_BITS)
	if err != nil {
		return Bet{}, err
	}

	return NewBet(uint8(agency), firstName, lastName, uint32(document), birthdate, uint16(number)), nil
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
