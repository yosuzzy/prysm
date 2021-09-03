package fuzz

import (
	"github.com/gogo/protobuf/proto"
	pubsub_pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
)

func FuzzMessageID(b []byte) int {
	MessageIDFuzz(b)
	return 0
}

func MessageIDFuzz(b []byte) {
	msg := &pubsub_pb.Message{}
	if err := proto.Unmarshal(b, msg); err != nil {
		return
	}
	genesisValidatorsRoot := bytesutil.PadTo([]byte{'A'}, 32)
	p2p.MsgID(genesisValidatorsRoot, msg)
}