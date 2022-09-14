package kv

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/v3/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v3/testing/assert"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
	"github.com/prysmaticlabs/prysm/v3/testing/util"
	"go.etcd.io/bbolt"
)

func Test_migrateBlindedBeaconBlocks(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, db *bbolt.DB)
		eval  func(t *testing.T, db *bbolt.DB)
	}{
		{
			name: "migrates and deletes entries",
			setup: func(t *testing.T, db *bbolt.DB) {
				b := util.NewBeaconBlockBellatrix()
				b.Block.Slot = 10
				b.Block.Body.ExecutionPayload.BlockNumber = 5
				signedBlock, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				encodedBlock, err := marshalBlock(context.Background(), signedBlock)
				require.NoError(t, err)
				blockRoot, err := signedBlock.Block().HashTreeRoot()
				require.NoError(t, err)
				err = db.Update(func(tx *bbolt.Tx) error {
					return tx.Bucket(blocksBucket).Put(blockRoot[:], encodedBlock)
				})
				assert.NoError(t, err)
			},
			eval: func(t *testing.T, db *bbolt.DB) {
				err := db.View(func(tx *bbolt.Tx) error {
					b := util.NewBeaconBlockBellatrix()
					b.Block.Slot = 10
					b.Block.Body.ExecutionPayload.BlockNumber = 5
					signedBlock, err := blocks.NewSignedBeaconBlock(b)
					require.NoError(t, err)
					blockRoot, err := signedBlock.Block().HashTreeRoot()
					require.NoError(t, err)
					v := tx.Bucket(blocksBucket).Get(blockRoot[:])
					decodedBlock, err := unmarshalBlock(context.Background(), v)
					require.NoError(t, err)
					require.Equal(t, true, decodedBlock.IsBlinded())
					decodedBlockRoot, err := decodedBlock.Block().HashTreeRoot()
					require.NoError(t, err)
					require.DeepEqual(t, blockRoot, decodedBlockRoot)
					return nil
				})
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t).db
			tt.setup(t, db)
			assert.NoError(
				t,
				migrateBlindedBeaconBlocksEnabled(context.Background(), db),
				"migrateBlindedBeaconBlocksEnabled(tx) error",
			)
			tt.eval(t, db)
		})
	}
}
