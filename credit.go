package bank

type Credit struct {
	Id          int
	CustomerId  int
	Value       int `db:value`
	Percentage  int `db:percentage`
	LoanPeriod  int
	CurrentDebt int
}

// ОБЯЗАТЕЛЬНЫЕ ПЛАТЕЖИ
