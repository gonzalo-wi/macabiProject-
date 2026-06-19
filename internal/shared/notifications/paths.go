package notifications

import "fmt"

// Rutas del frontend (React Router). El service worker abre URLs relativas al origen.

func StockRequestDetail(requestID string) string {
	return fmt.Sprintf("/app/stock/requests/%s", requestID)
}

func ExpenseDetail(expenseID string) string {
	return fmt.Sprintf("/app/gastos/%s", expenseID)
}

func EventRespond(eventID string) string {
	return fmt.Sprintf("/app/jornadas/%s/responder", eventID)
}

func NewsDetail(newsID string) string {
	return fmt.Sprintf("/app/noticias/%s", newsID)
}
