package userdomain

import "fmt"

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"
)

var validRoles = map[Role]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleUser:       true,
}

func NewRole(value string) (Role, error) {
	r := Role(value)
	if !validRoles[r] {
		return "", fmt.Errorf("%w: %s", ErrInvalidRole, value)
	}
	return r, nil
}

func (r Role) String() string {
	return string(r)
}

func (r Role) IsAtLeast(required Role) bool {
	return roleHierarchy[r] >= roleHierarchy[required]
}

var roleHierarchy = map[Role]int{
	RoleUser:       1,
	RoleAdmin:      2,
	RoleSuperAdmin: 3,
}
