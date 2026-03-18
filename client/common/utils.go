package common

type Bet struct {
	Agency    uint8
	FirstName string
	LastName  string
	Document  uint32
	Birthdate string
	Number    uint16
}

func NewBet(agency uint8, firstName, lastName string, document uint32, birthdate string, number uint16) Bet {
	return Bet{
		Agency:    agency,
		FirstName: firstName,
		LastName:  lastName,
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}
}