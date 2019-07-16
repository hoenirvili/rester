package permission

type Permissions int

const (
	NoPermission Permissions = 1 << iota
	Anonymous
	Basic
	Admin
	Super
)

var supportedPermission = []Permissions{
	Anonymous, Basic, Admin, Super,
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
