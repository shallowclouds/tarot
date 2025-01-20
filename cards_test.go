package tarot_test

import (
	"testing"

	"github.com/shallowclouds/tarot"
	"github.com/stretchr/testify/require"
)

// TestCardZhString 测试卡牌中文字符串格式化功能
func TestCardZhString(t *testing.T) {
	testCases := []struct {
		name     string
		card     tarot.Card
		expected string
	}{
		{
			name: "正位测试",
			card: tarot.Card{
				Name:     "The Fool",
				ZhName:   "愚者",
				Position: tarot.PositionUpright,
			},
			expected: "愚者（正位）",
		},
		{
			name: "逆位测试",
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

// TestGetDefaultAssets 测试默认资源（卡牌、图片、字体）的加载功能
func TestGetDefaultAssets(t *testing.T) {
	assets := tarot.GetDefaultAssets()
	require.NotEmpty(t, assets.Cards)
	require.NotNil(t, assets.BackgroundImg)
	require.NotNil(t, assets.Font)
	require.NotNil(t, assets.AskerImg)
	require.NotNil(t, assets.ReaderImg)
}
