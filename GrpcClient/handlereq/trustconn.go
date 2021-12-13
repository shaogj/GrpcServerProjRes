package handlereq
//1124Testing
//package main

import (
	auth "20210810BFLProj/grpcSimpleService1017/GrpcClient/KeyStore"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/client_http"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"strconv"

	"net/http"
	//TMcrypte "github.com/tendermint/tendermint/crypto"
	cryptoloc "20210810BFLProj/grpcSimpleService1017/crypto"
	//TMcrypteCur "github.com/tendermint/tendermint/crypto"
	//"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoed255519 "20210810BFLProj/grpcSimpleService1017/crypto/ed25519"
	//124:
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/types"
)
import (
	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()
const (
	PrivKeyAminoName = "tendermint/PrivKeyEd25519"
	PubKeyAminoName  = "tendermint/PubKeyEd25519"
)
func init() {
	//RegisterƒBlockchainMessages(cdc)
	//types.RegisterBlockAmino(cdc)
}

//1108add for generate
func CreateKey() (privs, addrs string) {
	//创建私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//fmt.Println(address)
	return priv, address
}
//1105


//11119add

func RequestTrustInfo(curUrlVerify string) (uint, interface{},error) {
	var reqInfo proto.TrustQueryReq
	reqInfo.TxSignReq = "signstr123===004"

	//UrlVerify := "http://192.168.1.221:46657/tri_bc_tx_commit"
	//1118doing:
	//UrlVerify :="http://49.233.196.33:46657/tri_block_validators"
	UrlVerify := curUrlVerify
		//fix to this,每次变量初始化
	resQuerySign := proto.RPCResponse{}
	//trustInfo := &proto.ResultBroadcastTxCommit{}
	//1118doing
	trustInfo := &proto.ResultValidators1117{}
	//resQuerySign.Result = &trustInfo

	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo is:%v", UrlVerify, reqInfo)

	reqBody, err := json.Marshal(&reqInfo)
	if nil != err {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return 0, nil,nil
	}
	//1107add
	/*
	getPriv,getAddr :=CreateKey()
	fmt.Println("after CreateKey() get getPriv is:%s,getAddr is:%s,len(getAddr) is:%d", getPriv,getAddr,len(getAddr))
	return 0, nil,nil
	*/
	GSettleAccessKey :="a01b1efa9cdc076ed4f09769a62546d033604b6925e174d948475d24b5c31ab7"
	AccessKeyAddr := "0x3B174bf7027CbA74B6d58bC1030132a287F65C67"
	//1207doing,,验证key的签名
	/*
	privkey:
	a01b1efa9cdc076ed4f09769a62546d033604b6925e174d948475d24b5c31ab7
	addr:
	0x3B174bf7027CbA74B6d58bC1030132a287F65C67
	HActionTurstSign     = “BFL-ActionSign"
	*/
	//	if signData, err = auth.KSign(reqBody, GSettleAccessKey); err != nil {
	var HActionTurstSignKey = []byte("BFL-ActionSign")
	var signData string
	if signData, err = auth.KSign(HActionTurstSignKey, GSettleAccessKey); err != nil {
		fmt.Println("In trustQuery(),auth.KSign failed,signData is :%v,len(signData) is:%d,err is:%v", signData, len(signData),err)
		return 0,nil, nil
	}else{
		fmt.Println("In trustQuery(),auth.KSign succ!,get signData is :%v,len(signData) is:%d,err is:%v", signData, len(signData),err)

	}
	//1207,,to add to head:signData

	//err :
	err = auth.KAuth(AccessKeyAddr, signData, HActionTurstSignKey)//reqBody
	fmt.Println("In trustQuery(),auth.KAuth finished!,get err is:%v", err)
	//return 0,nil, nil
	//end 1107add
	reader := bytes.NewReader(reqBody)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, reader) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))
	//1207,,add to headsigndata
	//r.Header.Add("BFL-ActionSign", signData)
	//1209,
	r.Header.Add("signature", signData)

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0, nil,err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return 0, nil,err
	}
	fmt.Println("post sendTransactionPostForm success-----55,get res is :%v", string(content))
	//getResp := &proto.NodeResponse{}
	//if _, err = n.rpcClient.Call("validators", nil, vals); err != nil {

	cdc :=amino.NewCodec()
	cdc.RegisterInterface((*cryptoloc.PubKey)(nil), nil)
	//cdc.RegisterConcrete(&tendermint.PubKeyEd25519{}, "amino_test/MyStruct", nil)
	cdc.RegisterConcrete(cryptoed255519.PubKeyEd25519{},
		PubKeyAminoName, nil)

	//gettrustvals
	err = json.Unmarshal(content, &resQuerySign)
	if err != nil {
		fmt.Println("error content form rpc response to resQuerySign,err is: %v", err)
	}
	/*1207skep
	err = json.Unmarshal(resQuerySign.Result, trustInfo)
	if err !=nil {
		fmt.Println("resQuerySign.Result ummarsshal resQuerySign.Result is:%s,err=%v", string(resQuerySign.Result), err.Error())
	}
	*/
	_,err = client_http.UnmarshalResponseBytes(cdc, content, trustInfo)
	//1119====err = json.Unmarshal(content, &resQuerySign)
	if nil != err {
		fmt.Println("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		fmt.Println("json.Unmarshal err!!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
		return 0, nil,err
	}

	fmt.Println("json.Unmarshal succ!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
	for _,val := range trustInfo.Validators {
		//fmt.Println("get cur trustInfo id:%d,Validator.PubKey is:%s",key,string(val.PubKey.Bytes()))//02x
		pubkeyAddr :=val.PubKey.Address()
		//1207--skip==fmt.Println("get cur trustInfo id:%d,Validator.PubKey'Address--444, len is:%d,Address is:%s",key,len(pubkeyAddr),val.PubKey.Address())
		envcodeStr :=hex.EncodeToString([]byte(pubkeyAddr))
		fmt.Println("Validator.PubKey'Address--555 is: %s",envcodeStr)

		//result[i] = common.BytesToAddress(address)
		//hexstrpubaddr := fmt.Sprintf("%x",val.PubKey.Bytes())
		//log.Println("HexToAddress val.Address======44444-=== len is:%d, val is:%s",len(hexstrpubaddr),hexstrpubaddr)


		//1207watching err:
		//addresshex :=common.HexToAddress(string(val.Address))	//("12345")
		//log.Println("HexToAddress addresshex len is:%d, val is:%s",len(addresshex),addresshex)


		//PubKeystr :=common.HexToAddress(string(val.PubKey.Bytes()))	//("12345")
		//log.Println("HexToAddress PubKeystr val is:%v",PubKeystr)

	}
	return 1,trustInfo,nil

	//getValidator3,err := UnmarshalResponseTest(cdc,[]byte(stringinfo),&getValidator)
	//log.Println("RequestTrustInfo(),return getValidator3 is: %v", getValidator3)

	//fmt.Println("json.Unmarshal succ!,cur get resQuerySign.Result :%v,get resQuerySign is :%v", resQuerySign.Result, resQuerySign)
	//err = json.Unmarshal(resQuerySign.Result, trustInfo)
	//1118doing,,
	var getresult interface{}
	if false {
		err = json.Unmarshal(resQuerySign.Result, trustInfo)
	}else{
		err = cdc.UnmarshalJSON(resQuerySign.Result, trustInfo)
	}
	if nil != err {
		fmt.Println("Could not unmarshal bytes.resp's result is:=%s,err=%v", string(resQuerySign.Result), err.Error())
		return 0, nil,err
	}
	fmt.Println("json.Unmarshal succ!,cur getresult info is :%v",getresult)

	//CheckTx
	/*
	PeterAddr :=common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527")// Peter
	//MartinAddr :=common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595")// Martin
	ExtraData	:= hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa")
	fmt.Println("json.Unmarshal succ!,cur get ExtraData addr is :%v",ExtraData)
	*/
	/*
	GSettleAccessKey :="GSettleAccessKey"
	if signData, err = auth.KSign(reqBody, GSettleAccessKey); err != nil {
		log.Error("In trustQuery(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return 0, nil
	}
	//strRes, statusCode, errorCode, err := SendTransactionPostForm(UrlVerify, 9000, &reqInfo, &resQuerySign)
	for ino, curSettItem := range transInfo.Withdraws {
		// 2 审核通过(运营审核通过)
		log.Info("get SettlApiQuery queue info, cur ino is:%d,SettItem record Status is:%v,curSettItem is :%v", ino, curSettItem.Status, curSettItem)
	}*/
	return 1,resQuerySign.Result,nil

}

//1124add
type ResultBlock struct {
	BlockMeta *types.BlockMeta `json:"block_meta"`
	Block     *types.Block     `json:"block"`
}
func reqGetUrl(vertify string)(respstr string,err error){
	resp, err := http.Get(vertify)
	if err != nil {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return "",err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return "",err
	}
	fmt.Println(string(body))
	return string(body),nil
}
//11.24get trustInfo good!
func RequestNodeTrustScore(defUrlVerify string) (uint, interface{},error) {
	var reqInfo proto.TrustQueryReq
	reqInfo.TxSignReq = "signstr123===004"
	//UrlVerify ="http://49.233.196.33:46657/tri_block_validators"
	UrlVerify := "http://192.168.1.221:46657/tri_block_info"
	resQuerySign := proto.RPCResponse{}
	//trustInfo := &proto.ResultBroadcastTxCommit{}
	trustInfo := &proto.ResultValidators1117{}
	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo is:%v", UrlVerify, "reqInfo")
	trustBlockInfo :=&ResultBlock{}
	//resQuerySign.Result = &trustInfo
	if true {
		getblockInfo,err :=reqGetUrl(UrlVerify)
		fmt.Println("RequestNodeTrustScore.reqGetUrl(),get getblockInfo is:%s,err is:%v", getblockInfo, err)
		_,err = client_http.UnmarshalResponseBytes(cdc, []byte(getblockInfo), trustBlockInfo)
		//1119====err = json.Unmarshal(content, &resQuerySign)
		if nil != err {
			fmt.Println("resp=%s,url=%s,err=%v", string(getblockInfo), UrlVerify, err.Error())
			fmt.Println("UnmarshalResponseBytes err!!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
			return 0, nil,err
		}else{
			//trustBlockInfo.Block.ValidatorsHash
			fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Header.Height is:%s,Header.time is:%v", trustBlockInfo.BlockMeta.Header.Height,trustBlockInfo.Block.Header.Time)
			fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Header.total_txs is:%d", trustBlockInfo.Block.Header.TotalTxs)
			fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Block.ValidatorsHash is:%s", trustBlockInfo.Block.ValidatorsHash)
			fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Block.trust_data is:%v", trustBlockInfo.Block.TrustData)

		}
		return 1,trustBlockInfo,nil
	}
	reqBody, err := json.Marshal(&reqInfo)
	if nil != err {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return 0, nil,nil
	}

	reader := bytes.NewReader(reqBody)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, reader) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0, nil,err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return 0, nil,err
	}
	fmt.Println("post sendTransactionPostForm success-----55,get res is :%v", string(content))
	//		if _, err = n.rpcClient.Call("validators", nil, vals); err != nil {


	cdc :=amino.NewCodec()
	//cdc.RegisterInterface((*TMcrypteCur.PubKey)(nil), nil)
	cdc.RegisterInterface((*cryptoloc.PubKey)(nil), nil)
	//cdc.RegisterInterface((*TMcrypte.ed25519.PubKey)(nil), nil)
	//cdc.RegisterConcrete(&tendermint.PubKeyEd25519{}, "amino_test/MyStruct", nil)
	cdc.RegisterConcrete(cryptoed255519.PubKeyEd25519{},
		PubKeyAminoName, nil)

	//gettrustvals
	err = json.Unmarshal(content, &resQuerySign)
	if err != nil {
		fmt.Println("error content form rpc response to resQuerySign,err is: %v", err)
	}
	err = json.Unmarshal(resQuerySign.Result, trustInfo)
	if err !=nil {
		fmt.Println("resQuerySign.Result ummarsshal resQuerySign.Result is:%s,err=%v", string(resQuerySign.Result), err.Error())
	}
	_,err = client_http.UnmarshalResponseBytes(cdc, content, trustInfo)
	//1119====err = json.Unmarshal(content, &resQuerySign)
	if nil != err {
		fmt.Println("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		fmt.Println("json.Unmarshal err!!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
		return 0, nil,err
	}

	fmt.Println("json.Unmarshal succ!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
	//1209PM，tmp to get tendermint's key:
	//trustInfo.Validators[0].PubKey =
	for key,val := range trustInfo.Validators {
		//fmt.Println("get cur trustInfo,val.PubKey is %s",val.PubKey)
		fmt.Println("get cur trustInfo id:%d,Validator.PubKey is:%s",key,string(val.PubKey.Bytes()))//02x
		pubkeyAddr :=val.PubKey.Address()
		fmt.Println("get cur trustInfo id:%d,Validator.PubKey'Address--444, len is:%d,Address is:%s",key,len(pubkeyAddr),val.PubKey.Address())
		envcodeStr :=hex.EncodeToString([]byte(pubkeyAddr))
		fmt.Println("Validator.PubKey'Address--555 is: %s",envcodeStr)
		strval :=[]byte("1B6011C07B")
		envcodetest :=hex.EncodeToString([]byte(strval))
		fmt.Println("my test=444 envcodetest is: %s",envcodetest)


	}
	return 1,trustInfo,nil
}