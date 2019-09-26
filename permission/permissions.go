// Package permission holds types and functions used for defining
// the permission of resources
package permission

// Permissions holds a byte pattern of permissions
type Permissions int

const (
	// NoPermission means that the handler resource does not have any
	// permission set
	NoPermission Permissions = 1 << iota
	// Anonymous means that anyone can access this handler resource
	Anonymous
	// Basic means that anyone with basic permissions can access the
	// requested handler resource
	Basic
	// Admin means that anyone with admin permissions can access the
	// requested handler resource
	Admin
	// Super means that anyone with super permissions can access the
	// requested handler resource
	Super
)

// supportedPermission by default
var supportedPermission = []Permissions{
	Anonymous,
	Basic,
	Admin,
	Super,
}

// Valid checks if the given permission is valid and supported
func (p Permissions) Valid() bool {
	for _, permission := range supportedPermission {
		if permission == p {
			return true
		}
	}
	return false
}

// Append appends a new permission in the underlying supported permissions
// The newly permission added now can be used in the route guard API
func Append(p Permissions) {
	supportedPermission = append(supportedPermission, p)
}
