package common

import (
	"encoding/csv"
	"fmt"
	"os"
)

const FILEPATH = "/data/agency-"

type Bet struct {
	FirstName string
	LastName  string
	Document  uint32
	Birthdate string
	Number    uint16
}

func NewBet(firstName, lastName string, document uint32, birthdate string, number uint16) Bet {
	return Bet{
		FirstName: firstName,
		LastName:  lastName,
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}
}

type BetReader struct {
	Agency    string
	Finished  bool
	Reader    *csv.Reader
	File      *os.File
	BatchSize int
}

func NewBetReader(agency string, batchSize int) (*BetReader, error) {
	path := FILEPATH + agency + ".csv"
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return &BetReader{
		Agency:    agency,
		Finished:  false,
		Reader:    csv.NewReader(f),
		File:      f,
		BatchSize: batchSize,
	}, nil
}