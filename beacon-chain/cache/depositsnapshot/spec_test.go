package depositsnapshot

import (
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v3/testing/assert"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	DepositData     depositData `yaml:"deposit_data"`
	DepositDataRoot [32]byte    `yaml:"deposit_data_root"`
	Eth1Data        eth1Data    `yaml:"eth1_data"`
	BlockHeight     uint        `yaml:"block_height"`
	Snapshot        snapshot    `yaml:"snapshot"`
}

func (tc *testCase) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		DepositData     depositData `yaml:"deposit_data"`
		DepositDataRoot string      `yaml:"deposit_data_root"`
		Eth1Data        eth1Data    `yaml:"eth1_data"`
		BlockHeight     string      `yaml:"block_height"`
		Snapshot        snapshot    `yaml:"snapshot"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	tc.DepositDataRoot, err = hexStringToByteArray(raw.DepositDataRoot)
	if err != nil {
		return err
	}
	tc.DepositData = raw.DepositData
	tc.Eth1Data = raw.Eth1Data
	tc.BlockHeight, err = stringToUint(raw.BlockHeight)
	if err != nil {
		return err
	}
	tc.Snapshot = raw.Snapshot
	return nil
}

type depositData struct {
	Pubkey                []byte `yaml:"pubkey"`
	WithdrawalCredentials []byte `yaml:"withdrawal_credentials"`
	Amount                uint64 `yaml:"amount"`
	Signature             []byte `yaml:"signature"`
}

func (dd *depositData) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		Pubkey                string `yaml:"pubkey"`
		WithdrawalCredentials string `yaml:"withdrawal_credentials"`
		Amount                string `yaml:"amount"`
		Signature             string `yaml:"signature"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	dd.Pubkey, err = hexStringToBytes(raw.Pubkey)
	if err != nil {
		return err
	}
	dd.WithdrawalCredentials, err = hexStringToBytes(raw.WithdrawalCredentials)
	if err != nil {
		return err
	}
	dd.Amount, err = strconv.ParseUint(raw.Amount, 10, 64)
	if err != nil {
		return err
	}
	dd.Signature, err = hexStringToBytes(raw.Signature)
	if err != nil {
		return err
	}
	return nil
}

type eth1Data struct {
	DepositRoot  [32]byte `yaml:"deposit_root"`
	DepositCount uint     `yaml:"deposit_count"`
	BlockHash    [32]byte `yaml:"block_hash"`
}

func (ed *eth1Data) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		DepositRoot  string `yaml:"deposit_root"`
		DepositCount string `yaml:"deposit_count"`
		BlockHash    string `yaml:"block_hash"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	ed.DepositRoot, err = hexStringToByteArray(raw.DepositRoot)
	if err != nil {
		return err
	}
	ed.DepositCount, err = stringToUint(raw.DepositCount)
	if err != nil {
		return err
	}
	ed.BlockHash, err = hexStringToByteArray(raw.BlockHash)
	if err != nil {
		return err
	}
	return nil
}

type snapshot struct {
	Finalized            [][32]byte `yaml:"finalized"`
	DepositRoot          [32]byte   `yaml:"deposit_root"`
	DepositCount         uint       `yaml:"deposit_count"`
	ExecutionBlockHash   [32]byte   `yaml:"execution_block_hash"`
	ExecutionBlockHeight uint       `yaml:"execution_block_height"`
}

func (sd *snapshot) UnmarshalYAML(value *yaml.Node) error {
	raw := struct {
		Finalized            []string `yaml:"finalized"`
		DepositRoot          string   `yaml:"deposit_root"`
		DepositCount         string   `yaml:"deposit_count"`
		ExecutionBlockHash   string   `yaml:"execution_block_hash"`
		ExecutionBlockHeight string   `yaml:"execution_block_height"`
	}{}
	err := value.Decode(&raw)
	if err != nil {
		return err
	}
	sd.Finalized = make([][32]byte, len(raw.Finalized))
	for i, finalized := range raw.Finalized {
		sd.Finalized[i], err = hexStringToByteArray(finalized)
		if err != nil {
			return err
		}
	}
	sd.DepositRoot, err = hexStringToByteArray(raw.DepositRoot)
	if err != nil {
		return err
	}
	sd.DepositCount, err = stringToUint(raw.DepositCount)
	if err != nil {
		return err
	}
	sd.ExecutionBlockHash, err = hexStringToByteArray(raw.ExecutionBlockHash)
	if err != nil {
		return err
	}
	sd.ExecutionBlockHeight, err = stringToUint(raw.ExecutionBlockHeight)
	if err != nil {
		return err
	}
	return nil
}

func readTestCases(filename string) ([]testCase, error) {
	var testCases []testCase
	file, err := os.ReadFile(filename)
	if err != nil {
		return []testCase{}, err
	}
	err = yaml.Unmarshal(file, &testCases)
	if err != nil {
		return []testCase{}, err
	}
	return testCases, nil
}

func TestRead(t *testing.T) {
	tcs, err := readTestCases("test_cases.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range tcs {
		t.Log(tc)
	}
}

func hexStringToByteArray(s string) (b [32]byte, err error) {
	var raw []byte
	raw, err = hexStringToBytes(s)
	if err != nil {
		return
	}
	if len(raw) != 32 {
		err = errors.New("invalid hex string length")
		return
	}
	copy(b[:], raw[:32])
	return
}

func hexStringToBytes(s string) (b []byte, err error) {
	b, err = hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return
}

func stringToUint(s string) (uint, error) {
	value, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}

func TestDepositCases(t *testing.T) {
	dt := NewDepositTree()
	testCases, err := readTestCases("test_cases.yaml")
	assert.NoError(t, err)
	for _, c := range testCases {
		err = dt.AddDeposit(c.DepositDataRoot, 0)
		assert.NoError(t, err)
		assert.Equal(t, c.Snapshot.DepositRoot, c.Eth1Data.DepositRoot)
		assert.Equal(t, dt.GetRoot(), c.Eth1Data.DepositRoot)
	}
}
