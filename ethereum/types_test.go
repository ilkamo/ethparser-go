package ethereum

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/tests/testdata"
	"github.com/ilkamo/ethparser-go/types"
)

var expectedBlock = types.Block{
	Number:     19697111,
	Hash:       "0xc8c7f99d64c6678ac5910f569167356808550b4e8fe22e8787963a62aff66d88",
	ParentHash: "0x91c90676cab257a59cd956d7cb0bceb9b1a71d79755c23c7277a0697ccfaf8c4",
	Timestamp:  time.Unix(1439799153, 0),
	Transactions: []types.Transaction{
		{
			BlockHash:   "0xed6fe3d8722be4b4614bc4fd2cc452d1d03ccdf453bc664b756a626d32ee91af",
			BlockNumber: 19697111,
			Hash:        "0xfa5109806d00fdfe9d0b73f9e9c2c59efd61a197900dcb03faff88c5fe263207",
			From:        "0x2e220f48eab381507f627a3e96f5387885619e83",
			To:          "0xb584d4be1a5470ca1a8778e9b86c81e165204599",
			Value:       "1300000000000",
		},
		{
			BlockHash:   "0xed6fe3d8722be4b4614bc4fd2cc452d1d03ccdf453bc664b756a626d32ee91af",
			BlockNumber: 19697111,
			Hash:        "0xe7d8be4e841d3ccda0f790ec0c57e483b1795c2a2f4f3b0a6b37dfa1f1ee8fd2",
			From:        "0x264bd8291fae1d75db2c5f573b07faa6715997b5",
			To:          "0xa6e127536a7b9aca15c928f6332fc9d2cd2e93c8",
			Value:       "636084590000000000",
		},
	},
}

func Test_block_ToBlock(t *testing.T) {
	t.Run("should convert block to types.Block", func(t *testing.T) {
		b := block{}

		// Unmarshal a real block
		err := json.Unmarshal(testdata.BlockJSON, &b)
		require.NoError(t, err)

		// Convert the block to a types.Block
		gotBlock, err := b.ToBlock()
		require.NoError(t, err)
		require.Equal(t, expectedBlock, gotBlock)
	})

	t.Run("should error because of bad block number", func(t *testing.T) {
		b := block{Number: "0x"}
		_, err := b.ToBlock()
		require.ErrorContains(t, err, "could not decode block number")
	})

	t.Run("should error because of bad block timestamp", func(t *testing.T) {
		b := block{
			Number:    "0x1",
			Timestamp: "0x",
		}
		_, err := b.ToBlock()
		require.ErrorContains(t, err, "could not decode block timestamp")
	})

	t.Run("should error because of bad transaction block number", func(t *testing.T) {
		b := block{
			Number:    "0x1",
			Timestamp: "0x55d19771",
			Transactions: []transaction{
				{
					BlockNumber: "0x",
				},
			},
		}
		_, err := b.ToBlock()
		require.ErrorContains(t, err, "could not decode tx block number")
	})

	t.Run("should error because of bad transaction value", func(t *testing.T) {
		b := block{
			Number:    "0x1",
			Timestamp: "0x55d19771",
			Transactions: []transaction{
				{
					BlockNumber: "0x1",
					Value:       "0x",
				},
			},
		}
		_, err := b.ToBlock()
		require.ErrorContains(t, err, "could not decode tx value")
	})
}
