package tarot_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/shallowclouds/tarot"
)

func TestCardZhString(t *testing.T) {
	testCases := []struct {
		name     string
		card     tarot.Card
		expected string
	}{
		{
			name: "upright position",
			card: tarot.Card{
				Name:     "The Fool",
				ZhName:   "愚者",
				Position: tarot.PositionUpright,
			},
			expected: "愚者（正位）",
		},
		{
			name: "reversed position",
			card: tarot.Card{
				Name:     "The Fool",
				ZhName:   "愚者",
				Position: tarot.PositionReversed,
			},
			expected: "愚者（逆位）",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.card.ZhString())
		})
	}
}

func TestGetDefaultAssets(t *testing.T) {
	assets := tarot.GetDefaultAssets()
	require.NotEmpty(t, assets.Cards)
	require.NotNil(t, assets.BackgroundImg)
	require.NotNil(t, assets.Font)
	require.NotNil(t, assets.AskerImg)
	require.NotNil(t, assets.ReaderImg)
}
