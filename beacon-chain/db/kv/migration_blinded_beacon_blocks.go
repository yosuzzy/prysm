package kv

import (
	"bytes"
	"context"

	"github.com/prysmaticlabs/prysm/v3/config/features"
	"github.com/prysmaticlabs/prysm/v3/monitoring/progress"
	bolt "go.etcd.io/bbolt"
)

var (
	migrationBlindedBeaconBlocksKey = []byte("blinded-beacon-blocks-enabled")
	blindedBlocksMigrationMsg       = "Migrating all blocks in the database to a blinded format, this may take some time " +
		"if you have a big database. Alternatively, we recommend using checkpoint sync according to our documentation " +
		"for a faster initialization"
)

func migrateBlindedBeaconBlocksEnabled(ctx context.Context, db *bolt.DB) error {
	if !features.Get().EnableOnlyBlindedBeaconBlocks {
		return nil // Only write to the migrations bucket if the feature is enabled.
	}
	if updateErr := db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket(migrationsBucket)
		if b := mb.Get(migrationBlindedBeaconBlocksKey); bytes.Equal(b, migrationCompleted) {
			return nil // Migration already completed.
		}

		// Load all blocks and put them to disk as blinded beacon block .
		bkt := tx.Bucket(blocksBucket)
		numItems := bkt.Stats().KeyN
		bar := progress.InitializeProgressBar(numItems, blindedBlocksMigrationMsg)
		if err := bkt.ForEach(func(k, v []byte) error {
			defer func() {
				if err := bar.Add(1); err != nil {
					log.Error(err)
				}
			}()
			if v == nil {
				return nil
			}
			decoded, err := unmarshalBlock(ctx, v)
			if err != nil {
				return err
			}
			blindedFormat, err := decoded.ToBlinded()
			if err != nil {
				return err
			}
			encoded, err := marshalBlock(ctx, blindedFormat)
			if err != nil {
				return err
			}
			return bkt.Put(k, encoded)
		}); err != nil {
			return err
		}

		return mb.Put(migrationBlindedBeaconBlocksKey, migrationCompleted)
	}); updateErr != nil {
		return updateErr
	}
	return nil
}
