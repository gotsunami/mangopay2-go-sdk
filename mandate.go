package mango

type Mandate struct {
	ProcessReply
	//An ID of a Bank Account
	BankAccountId string
	// Id of the author.
	UserId string
	// Currency of the registered card.
	ReturnURL     string
	RedirectURL   string
	DocumentURL   string
	Culture       string
	Scheme        string
	ExecutionType string
	MandateType   string
	BankReference string

	service *MangoPay
}

func (m *MangoPay) Mandate(id string) (*Mandate, error) {
	any, err := m.anyRequest(new(Mandate), actionFetchMandate, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return any.(*Mandate), nil
}
