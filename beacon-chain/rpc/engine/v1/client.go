package v1

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	enginev1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
)

const (
	newPayloadRequest        = "engine_newPayloadV1"
	forkchoiceUpdatedRequest = "engine_forkchoiceUpdatedV1"
	getPayloadRequest        = "engine_getPayloadV1"
)

// Custom error for handling -32000 server errors from the engine API's
// JSON-RPC specification. More details can be found in the execution-apis
// repository.
type ServerError struct {
	data string
}

// Error interface method.
func (s *ServerError) Error() string {
	return fmt.Sprintf("server error: %v", s.data)
}

// ForkChoiceUpdates response for the engine_forkchoiceUpdatedV1 method.
type ForkchoiceUpdatedResponse struct {
	PayloadStatus *enginev1.PayloadStatus
	PayloadID     [8]byte
}

// Client --
type Client struct {
	client *rpc.Client
}

// New --
func New(endpoint string) (*Client, error) {
	rpcClient, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: rpcClient,
	}, nil
}

// NewPayload is a wrapper to call the engine_newPayloadV1 JSON-RPC endpoint.
func (c *Client) NewPayload(payload *enginev1.ExecutionPayload) (*enginev1.PayloadStatus, error) {
	return nil, errors.New("unimplemented")
}

// ForkchoiceUpdated is a wrapper to call the engine_forkchoiceUpdatedV1 JSON-RPC endpoint.
func (c *Client) ForkchoiceUpdated(
	state *enginev1.ForkchoiceState, attributes *enginev1.PayloadAttributes,
) (*ForkchoiceUpdatedResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetPayload is a wrapper to call the engine_getPayloadV1 JSON-RPC endpoint.
func (c *Client) GetPayload(payloadID [8]byte) (*enginev1.ExecutionPayload, error) {
	return nil, errors.New("unimplemented")
}
