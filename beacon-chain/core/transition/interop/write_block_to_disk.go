package interop

import (
	"fmt"
	"os"
	"path"

	"github.com/prysmaticlabs/prysm/v4/config/features"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

// WriteBlockToDisk as a block ssz. Writes to temp directory. Debug!
func WriteBlockToDisk(block interfaces.ReadOnlySignedBeaconBlock, failed bool) {
	if !features.Get().WriteSSZStateTransitions {
		return
	}

	filename := fmt.Sprintf("beacon_block_%d.ssz", block.Block().Slot())
	if failed {
		filename = "failed_" + filename
	}
	fp := path.Join(os.TempDir(), filename)
	log.Warnf("Writing block to disk at %s", fp)
	enc, err := block.MarshalSSZ()
	if err != nil {
		log.WithError(err).Error("Failed to ssz encode block")
		return
	}
	if err := file.WriteFile(fp, enc); err != nil {
		log.WithError(err).Error("Failed to write to disk")
	}
}

func WriteAttestationToDisk(att *eth.Attestation) {
	d := att.Data
	filename := fmt.Sprintf("att_%d_%d.ssz", d.Slot, d.CommitteeIndex)
	fp := path.Join(os.TempDir(), filename)
	log.Warnf("Writing attestation to disk at %s", fp)
	enc, err := att.MarshalSSZ()
	if err != nil {
		log.WithError(err).Error("Failed to ssz encode block")
		return
	}
	if err := file.WriteFile(fp, enc); err != nil {
		log.WithError(err).Error("Failed to write to disk")
	}
}
