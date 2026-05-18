package stockdomain

import (
	"strings"
	"time"
)

type ResourceType string

const (
	ResourceTypeReturnable ResourceType = "returnable"
	ResourceTypeConsumable ResourceType = "consumable"
)

type RequestStatus string

const (
	RequestStatusPending   RequestStatus = "PENDIENTE"
	RequestStatusApproved  RequestStatus = "RESERVADO"
	RequestStatusRejected  RequestStatus = "RECHAZADO"
	RequestStatusDelivered RequestStatus = "ENTREGADO"
	RequestStatusReturned  RequestStatus = "DEVUELTO"
	RequestStatusCancelled RequestStatus = "CANCELADO"
)

type Resource struct {
	ID             string
	Name           string
	Type           ResourceType
	TotalStock     int
	AvailableStock int
	CreatedAt      time.Time
}

func NewResource(name string, rtype ResourceType, totalStock int) (*Resource, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyResourceName
	}
	if totalStock < 0 {
		return nil, ErrInvalidStock
	}
	if rtype != ResourceTypeReturnable && rtype != ResourceTypeConsumable {
		return nil, ErrInvalidResourceType
	}
	return &Resource{
		Name:           name,
		Type:           rtype,
		TotalStock:     totalStock,
		AvailableStock: totalStock,
	}, nil
}

type ResourceRequest struct {
	ID             string
	ProjectID      string
	ResourceID     string
	RequestedByID  string
	Quantity       int
	WithdrawalDate time.Time
	ReturnDate     *time.Time
	Status         RequestStatus
	Notes          string
	CreatedAt      time.Time
}

func NewResourceRequest(
	projectID, resourceID, requestedByID string,
	quantity int,
	withdrawalDate time.Time,
	returnDate *time.Time,
	notes string,
) (*ResourceRequest, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if projectID == "" || resourceID == "" || requestedByID == "" {
		return nil, ErrMissingRequiredField
	}
	return &ResourceRequest{
		ProjectID:      projectID,
		ResourceID:     resourceID,
		RequestedByID:  requestedByID,
		Quantity:       quantity,
		WithdrawalDate: withdrawalDate,
		ReturnDate:     returnDate,
		Notes:          notes,
		Status:         RequestStatusPending,
	}, nil
}

type StockNotification struct {
	ID        string
	UserID    string
	RequestID string
	Message   string
	ReadAt    *time.Time
	CreatedAt time.Time
}

type RequestDetail struct {
	Request       ResourceRequest
	ResourceName  string
	ResourceType  ResourceType
	ProjectName   string
	RequesterName string
}
