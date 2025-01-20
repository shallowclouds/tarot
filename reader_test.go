package tarot_test

import (
	"context"
	"testing"

	"github.com/shallowclouds/tarot"
	"github.com/stretchr/testify/require"
)

func TestReaderChoose(t *testing.T) {
	r, err := tarot.NewReader(&tarot.DumbGPTReader{}, "", "", tarot.GetDefaultAssets())
	require.NoError(t, err)

	cards, err := r.Choose()
	require.NoError(t, err)
	require.Len(t, cards, 3)

	// Verify positions are either upright or reversed
	for _, card := range cards {
		require.Contains(t, []tarot.Position{tarot.PositionUpright, tarot.PositionReversed}, card.Position)
	}
}

func TestReaderPrompt(t *testing.T) {
	r, err := tarot.NewReader(&tarot.DumbGPTReader{}, "", "", tarot.GetDefaultAssets())
	require.NoError(t, err)

	cards := [3]tarot.Card{
		{Name: "The Fool", ZhName: "愚者", Position: tarot.PositionUpright},
		{Name: "The Magician", ZhName: "魔术师", Position: tarot.PositionReversed},
		{Name: "The High Priestess", ZhName: "女祭司", Position: tarot.PositionUpright},
	}

	prompt := r.Prompt(cards, "test question", "{{thing}} with {{card1}}, {{card2}}, {{card3}}")
	require.Contains(t, prompt, "test question")
	require.Contains(t, prompt, "愚者（正位）")
	require.Contains(t, prompt, "魔术师（逆位）")
	require.Contains(t, prompt, "女祭司（正位）")
}

func TestReaderDivineWithOption(t *testing.T) {
	r, err := tarot.NewReader(&tarot.DumbGPTReader{}, "", "", tarot.GetDefaultAssets())
	require.NoError(t, err)

	result, err := r.DivineWithOption(context.Background(), tarot.DivineOption{
		Question: "test question",
		Asker:    "test asker",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Cards, 3)
	require.NotNil(t, result.Img)
	require.NotEmpty(t, result.Result)
}
