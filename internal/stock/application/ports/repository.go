package stockports

type StockRepository interface {
	ResourceRepository
	RequestRepository
	NotificationRepository
}
