package permission

type Permissions int

const (
	Anonymous Permissions = 1 << iota
	Basic
	Admin
	Super
)
