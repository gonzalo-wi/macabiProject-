package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
)

type OptionRepository interface {
	CreateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error
	FindOptionGroupByID(ctx context.Context, id string) (*eventdomain.EventOptionGroup, error)
	UpdateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error
	DeleteOptionGroup(ctx context.Context, id string) error
	CreateOption(ctx context.Context, o *eventdomain.EventOption) error
	UpdateOption(ctx context.Context, o *eventdomain.EventOption) error
	FindOptionByID(ctx context.Context, id string) (*eventdomain.EventOption, error)
	DeleteOption(ctx context.Context, id string) error
}
