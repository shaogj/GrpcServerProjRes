package types

import (
	"fmt"
	//"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"strings"
)
type TrustPubKey string

func(t TrustPubKey) Hash() []byte {
	return tmhash.Sum([]byte(t))
}

func(t TrustPubKey) PubKey() string {
	return strings.Split(string(t), "/")[0]
}

type TrustPubKeys []TrustPubKey

func (ts *TrustPubKeys) Hash() []byte {
	tPks := make([][]byte, ts.Len())
	for i := 0; i < ts.Len(); i++ {
		tPks[i] = tmhash.Sum([]byte((*ts)[i]))
	}
	//return merkle.SimpleHashFromByteSlices(tPks)
	//sgj 1124
	return []byte("merkle.SimpleHashFromByteSlices")

}

func (ts *TrustPubKeys) Len() int {
	if ts == nil {
		return 0
	}
	return len(*ts)
}

func (ts *TrustPubKeys) StringIndented(indent string) string {
	if ts == nil {
		return "nil-TrustData"
	}
	tsStrings := make([]string, cmn.MinInt(ts.Len(), 11))
	for i, tp := range *ts {
		if i == 10 {
			tsStrings[i] = fmt.Sprintf("... (%v total)", ts.Len())
			break
		}
		tsStrings[i] = fmt.Sprintf("%X (%d bytes)", tp.Hash(), len(tp))
	}
	return fmt.Sprintf(`TrustData{
%s  %v
%s}#%v`,
		indent, strings.Join(tsStrings, "\n"+indent+"  "),
		indent, ts.Hash())
}

type TrustNodes struct {
	Nodes []string `json:"nodes"`
}

type ranking struct {
	Address string `json:"attestee"`
}

type Ranking struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data []ranking `json:"result"`
}
