package bankly_test

type AccountToTest struct {
	BankCode string
	Branch   string
	Account  string
	Document string
	Name     string
}

func accountA() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "189162",
		Document: "82895341000137",
		Name:     "NOME DA EMPRESA 1245312",
	}
}

func accountB() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "189081",
		Document: "59619372000143",
		Name:     "Nome da Empresa XVlBzgbaiC",
	}
}

func accountC() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "190420",
		Document: "45515165000134",
		Name:     "Nome da Empresa XVlBzgbaiC",
	}
}

func accountD() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "190411",
		Document: "90953987000151",
		Name:     "Nome da Empresa XVlBzgbaiC",
	}
}

func accountE() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "190403",
		Document: "50720146000180",
		Name:     "Nome da Empresa XVlBzgbaiC",
	}
}

func accountF() *AccountToTest {
	return &AccountToTest{
		BankCode: "332",
		Branch:   "0001",
		Account:  "190470",
		Document: "87568289000128",
		Name:     "Nome da Empresa XVlBzgbaiC",
	}
}
