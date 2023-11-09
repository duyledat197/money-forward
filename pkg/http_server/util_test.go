package http_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzMatchPath(f *testing.F) {
	path := "/users/{id}"

	testCases := []struct {
		path string
		want bool
	}{
		{
			path: "/users/123",
			want: true,
		},
		{
			path: "/users/123/",
			want: true,
		},
		{
			path: "/users",
			want: false,
		},
		{
			path: "/users/123/accounts",
			want: false,
		},
	}

	for _, testCase := range testCases {
		f.Add(testCase.path, testCase.want)
	}

	f.Fuzz(func(t *testing.T, a string, b bool) {
		assert.Equal(t, isMatchPath(path, a), b)
	})
}
