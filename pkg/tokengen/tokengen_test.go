package tokengen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenSize int
		want      int
	}{
		{
			name:      "Basic test",
			tokenSize: 10,
			want:      10,
		},
		{
			name:      "Small token size",
			tokenSize: 1,
			want:      1,
		},
		{
			name:      "Big token size",
			tokenSize: 1337,
			want:      1337,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := len(GenerateToken(tt.tokenSize))
			require.Equal(t, tt.want, l)
		})
	}
}
