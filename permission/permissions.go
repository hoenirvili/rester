package permission

type Permissions int

const (
	Anonymous Permissions = 1 << iota
	Basic
	Admin
	Super
)

var supportedPermission = [...]Permissions{
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
