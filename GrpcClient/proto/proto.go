package proto

import (
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/types"
	"encoding/json"
	//"github.com/tendermint/tendermint/crypto"
	"20210810BFLProj/grpcSimpleService1017/crypto"

)

//import common "github.com/tendermint/tendermint/libs/common"

// 用于解析证实信息
// Use to parse the attest information from the broadcast.
type AttestInformation struct {
	Action  string          `json:"action"`
	Ranking [][]interface{} `json:"ranking"`
}

type ranking struct {
	Address string `json:"attestee"`
}

type Ranking struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data []ranking `json:"result"`
}
//http://192.168.1.221:46657/tri_net_info
type TrustQueryReq struct {
	TxSignReq string `json:"tx"`

}
type TrustDataRespInfo struct {
	Address []string `json:"attestee"`

}


type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type RPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	CODE    int    `json:"code"`
	//Result	 []byte `json:"result,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	//sgj 1115add
	//Data interface{} `json:"Data,omitempty"`
	//Result  interface{} `json:"result,omitempty"`
	//Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}
type ResultBroadcastTxCommit struct {
	//CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
	//CheckTx   ResponseCheckTx   `json:"check_tx"`
	//DeliverTx ResponseDeliverTx `json:"deliver_tx"`
	Hash      []byte            `json:"hash"`
	Height    int64             `json:"height"`
}

//1118add
// Validators for a height
type ResultValidators1117 struct {
	BlockHeight int64              `json:"block_height"`
	Validators  []*Validator `json:"validators"`
	//Validators  []*types.Validator `json:"validators"`
}

type ResultValidators struct {
	BlockHeight string              `json:"block_height"`
	Validators  []*Validator `json:"validators"`
	//Validators  []*types.Validator `json:"validators"`
}

//Address is hex bytes.
type Address = crypto.Address

type Validator struct {
	Address     Address       `json:"address"`
	PubKey      crypto.PubKey `json:"pub_key"`
	//PubKey      PubKey `json:"pub_key"`
	VotingPower int64         `json:"voting_power"`
	ProposerPriority int64 `json:"proposer_priority"`
	//1119
	//VotingPower string         `json:"voting_power"`
	//ProposerPriority string `json:"proposer_priority"`
}

type PubKey interface {
	Address() Address
	Bytes() []byte
	VerifyBytes(msg []byte, sig []byte) bool
	Equals(PubKey) bool
}

//1124
type ResultBlockInfo struct {
	BlockMeta *types.BlockMeta `json:"block_meta"`
	Block     *types.Block     `json:"block"`
	Size      int              `json:"size"`
}
