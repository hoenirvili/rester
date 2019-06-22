package permission_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hoenirvili/rester/permission"
)

func TestValid(t *testing.T) {
	require := require.New(t)
	p := permission.Anonymous
	valid := p.Valid()
	require.True(valid)
}

func TestValidWithInvalid(t *testing.T) {
	require := require.New(t)
	p := permission.Permissions(10)
	valid := p.Valid()
	require.False(valid)
}

func TestAppend(t *testing.T) {
	require := require.New(t)
	p := permission.Permissions(15)
	permission.Append(p)
	valid := p.Valid()
	require.True(valid)
}
