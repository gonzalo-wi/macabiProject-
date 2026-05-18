package expensesusecases

import (
	"fmt"

	expensesdomain "macabi-back/internal/expenses/domain"
)

func formatExpenseAmount(exp *expensesdomain.Expense) string {
	currency := exp.Currency
	if currency == "" {
		currency = expensesdomain.DefaultCurrency
	}
	return fmt.Sprintf("%s %s", exp.Amount.StringFixed(2), currency)
}
