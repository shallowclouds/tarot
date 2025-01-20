package tarot_test

import (
	"context"
	"testing"

	"github.com/shallowclouds/tarot"
	"github.com/stretchr/testify/require"
)

func TestDumbGPTReader(t *testing.T) {
	testCases := []struct {
		name       string
		systemMsg  string
		userMsg    string
		expectErr  bool
		expectResp bool
	}{
		{
			name:       "empty messages",
			systemMsg:  "",
			userMsg:    "",
			expectErr:  false,
			expectResp: true,
		},
		{
			name:       "with system message",
			systemMsg:  "System message",
			userMsg:    "User message",
			expectErr:  false,
			expectResp: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &tarot.DumbGPTReader{}
			resp, err := r.Chat(context.Background(), tc.systemMsg, tc.userMsg)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.expectResp {
				require.NotEmpty(t, resp)
			} else {
				require.Empty(t, resp)
			}
		})
	}
}
