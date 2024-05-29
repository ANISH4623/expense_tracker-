package serializers

type UserSerializer struct {
	ID        uint                `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	CreatedAt string              `json:"CreatedAt"`
	UpdatedAt string              `json:"UpdatedAt"`
	Expenses  []ExpenseSerializer `json:"expenses"`
	Income    []IncomeSerializer  `json:"income"`
}
type ExpenseSerializer struct {
	ID       uint    `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
}

type IncomeSerializer struct {
	ID       uint    `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
}
