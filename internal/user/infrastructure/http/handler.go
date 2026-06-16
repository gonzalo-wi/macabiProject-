package userhttp

type Handlers struct {
	Auth       *AuthHandler
	User       *UserHandler
	Invitation *InvitationHandler
}

func NewHandlers(auth *AuthHandler, user *UserHandler, invitation *InvitationHandler) *Handlers {
	return &Handlers{Auth: auth, User: user, Invitation: invitation}
}
