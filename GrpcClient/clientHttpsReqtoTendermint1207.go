package main

import (
	auth "20210810BFLProj/grpcSimpleService1017/GrpcClient/KeyStore"
	//"20210810BFLProj/grpcSimpleService1017/GrpcClient/client_http"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/handlereq"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/handleencryreq"
	pb "20210810BFLProj/grpcSimpleService1017/proto"
	"context"
	"encoding/json"
	"errors"
	log "fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/go-amino"
	"google.golang.org/grpc"
	"strconv"
	"testing"
	"time"

	//1117add,sgj
	pbtrust "20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
	"fmt"

)
const (
	extraVanity      = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal        = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
	nextForkHashSize = 4  // Fixed number of extra-data suffix bytes reserved for nextForkHash.
	validatorBytesLength = 20
)
// routeList 调用服务端RouteList方法
func routeList() {
	//调用服务端RouteList方法，获流
	stream, err := streamClient.RouteList(context.Background())
	if err != nil {
		log.Println("Upload list err: %v", err)
	}
	for n := 0; n < 5; n++ {
		//向流中发送消息
		err := stream.Send(&pb.StreamRequest{StreamData: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Println("stream request err: %v", err)
		}

	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("RouteList get response err: %v", err)
	}
	log.Println(res)
}

// Address 连接地址
const Address string = ":8000"
const MaximumExtraDataSize  uint64 = 32
var streamClient pb.StreamClientClient

func ParseValidators(validatorsBytes []byte) ([]common.Address, error) {
	if len(validatorsBytes)%validatorBytesLength != 0 {
		return nil, errors.New("invalid validators bytes")
	}
	n := len(validatorsBytes) / validatorBytesLength
	result := make([]common.Address, n)
	for i := 0; i < n; i++ {
		address := make([]byte, validatorBytesLength)
		copy(address, validatorsBytes[i*validatorBytesLength:(i+1)*validatorBytesLength])
		result[i] = common.BytesToAddress(address)
	}
	return result, nil
}
type Request_Type int32

const (
	Request_SEND_MESSAGE Request_Type = 0
	Request_UPDATE_PEER  Request_Type = 1
)


var cdc = amino.NewCodec()

type TrustGrpcClient struct {
	client pbtrust.TrustClientClient
}

var Client *TrustGrpcClient

func NewTrustGrpc(trustAddr string) error {

	conn, err := grpc.Dial(trustAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc init is failed, err: %v", err)
	}

	trustclient := pbtrust.NewTrustClientClient(conn)

	Client = &TrustGrpcClient{client: trustclient}
	return nil
}

func GetTrustGrpcClient() (*TrustGrpcClient, error) {
	if Client == nil {
		return nil, fmt.Errorf("get lgl grpc client is nil")
	}
	return Client, nil
}
var grpcTimeout = time.Second *3
func (c *TrustGrpcClient) SendTrustRequestMsg(req *pbtrust.TrustRequest) (*pbtrust.TrustResponse, error) {
	ctx, _ := context.WithTimeout(context.Background(), grpcTimeout)
	// ErrorCode:200成功 ErrorDesc:错误描述
	//context.Background()
	gettrustRes,err := c.client.TrustRequestData(ctx,req)
	if err != nil || gettrustRes == nil {
		fmt.Errorf("RPC.TrustRequestData() err!,get gettrustRes is :%v,err: %v", gettrustRes,err)
		return nil, fmt.Errorf("sendGrpcMsg TrustRequestData failed, resp is nil")
	}else{
		log.Println("RPC.TrustRequestData() succ!,get gettrustRes is :%v,err: %v", gettrustRes,err)
		log.Println("gettrustRes sub info GetRanking.Address is :%s,GetRanking.Score is: %d,err: %v", gettrustRes.GetRanking.Address,*gettrustRes.GetRanking.Score,err)
	}
	//if resp.GetErrorCode() != conf.ERROR_CODE_SUCCESS {

	return gettrustRes, nil
}

//1130testing
func main() {
	/*1205bef--testingGRPCgood!
	now := time.Now().Unix()
	//sgj 1118doing,grpc request:

		err := NewTrustGrpc(Address)
	if err != nil {
		log.Println("NewTrustGrpc init failed. err: %v", err)
		return
	}
	curClient,_ :=GetTrustGrpcClient()
	var sendtype = pbtrust.TrustRequest_SEND_MESSAGE
	curGrpcRequest := &pbtrust.TrustRequest{
		Type:	&sendtype,
		SendMessage: &pbtrust.SendMessage{
			Id:    []byte("IdSendMessage"),
			Data:    []byte("msg"),
			Created: &now,
		},
	}
	resp, error :=curClient.SendTrustRequestMsg(curGrpcRequest)
	if error != nil || resp == nil {
		log.Println("RPC.TrustRequestData() err!,get gettrustRes is :%v,err: %v", resp, err)
		return
	}
	//1121testing
	return
	end1205*/
	cdc :=amino.NewCodec()
	cdc.RegisterInterface((*pbtrust.PubKey)(nil), nil)
	//1207,,test https request
	UrlVerify :="http://49.233.196.33:46657/tri_block_validators"

	retval,resultinfo,err :=handlereq.RequestTrustInfo(UrlVerify)
	if err != nil  || retval !=1 {
		log.Println("RequestTrustInfo(),return value is err: %v", err)
	}

	trustInfo := resultinfo.(*pbtrust.ResultValidators1117)
	for _,val := range trustInfo.Validators {
		//fmt.Println("get cur trustInfo,val.PubKey is %s",val.PubKey)
		fmt.Println("after handlereq.RequestTrustInfo(),get Validator.Address len is:%d,Address is:%s,VotingPower is:%d", len(val.PubKey.Address()), val.PubKey.Address(),val.VotingPower)
	}
	//1207doing
	UrlhttpsReq := "http://192.168.1.224:46657/tri_block_validators"
	retval,resultinfo,err =handleencryreq.RequestEncryTrustInfo(UrlhttpsReq)
	if err != nil  || retval !=1 {
		log.Println("RequestTrustInfo(),return value is err: %v", err)
	}

	trustInfo2 := resultinfo.(*pbtrust.ResultValidators1117)
	for _,val := range trustInfo2.Validators {
		//fmt.Println("get cur trustInfo,val.PubKey is %s",val.PubKey)
		fmt.Println("after handlereq.RequestTrustInfo(),get Validator.Address len is:%d,Address is:%s,VotingPower is:%d", len(val.PubKey.Address()), val.PubKey.Address(),val.VotingPower)
	}
	//1210doing:
	//to 请求score data
	trustRankingUrl := "http://192.168.1.224:46657/tri_block_info"
	retval,resultinfo,err =handleencryreq.RequestNodeTrustScore(trustRankingUrl)
	if err != nil  || retval !=1 {
		log.Println("RequestNodeTrustScore(),resultinfo is:%v,return value is err: %v", resultinfo,err)
	}
	return
	/*
	return
	stringinfo := resultinfo.(proto.ResultValidators1117{})
	var getValidator pbtrust.ResultValidators
	getValidator3,err := client_http.UnmarshalResponseTest(cdc,[]byte(stringinfo),&getValidator)
	log.Println("RequestTrustInfo(),return getValidator3 is: %v", getValidator3)

	return
	*/
	//end 1118,,testing
	// 连接服务器
	/*1204doing grcp checking
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Println("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	streamClient = pb.NewStreamClientClient(conn)
	routeList()
	//1117add
	trustclient := pbtrust.NewTrustClientClient(conn)

	//var code int32= 44
	var curtype = pbtrust.TrustRequest_SEND_MESSAGE
	curRequest := &pbtrust.TrustRequest{
		Type:	&curtype,
		SendMessage: &pbtrust.SendMessage{
			Id:    []byte("IdSendMessage"),
			Data:    []byte("msg"),
			Created: &now,
		},
	}
	gettrustRes,err := trustclient.TrustRequestData(context.Background(),curRequest)
	if err != nil {
		log.Println("RPC.TrustRequestData() err!,get gettrustRes is :%v,err: %v", gettrustRes,err)
	}else{
		log.Println("RPC.TrustRequestData() succ!,get gettrustRes is :%v,err: %v", gettrustRes,err)
		log.Println("gettrustRes sub info GetRanking.Address is :%s,GetRanking.Score is: %d,err: %v", gettrustRes.GetRanking.Address,*gettrustRes.GetRanking.Score,err)

	}
	end 1205befing*/
	//1117add
	return
	//var Extra []byte
	//Extra := []byte("0x00000000000000000000000000000000000000000000000000000000000000002a7cdd959bfe8d9487b2a43b33565295a698f7e26488aa4d1955ee33403f8ccb1d4de5fb97c7ade29ef9f4360c606c7ab4db26b016007d3ad0ab86a0ee01c3b1283aa067c58eab4709f85e99d46de5fe685b1ded8013785d6623cc18d214320b6bb6475978f3adfc719c99674c072166708589033e2d9afec2be4ec20253b8642161bc3f44")
	//Extra :=[]byte("0x00000000000000000000000000000000000000000000000000000000000000002a7cdd959bfe8d9487b2a43b33565295a698f7e26488aa4d1955ee33403f8ccb1d4de5fb97c7ade29ef9f4360c606c7ab4db26b016007d3ad0ab86a0ee01c3b1283aa067c58eab4709f85e99d46de5fe685b1ded8013785d6623cc18d214320b6bb6475978f3adfc719c99674c072166708589033e2d9afec2be4ec20253b8642161bc3f444f53679c1f3d472f7be8361c80a4c1e7e9aaf001d0877f1cfde218ce2fd7544e0b2cc94692d4a704debef7bcb61328b8f7166496996a7da21cf1f1b04d9b3e26a3d0772d4c407bbe49438ed859fe965b140dcf1aab71a96bbad7cf34b5fa511d8e963dbba288b1960e75d64430b3230294d12c6ab2aac5c2cd68e80b16b581ea0a6e3c511bbd10f4519ece37dc24887e11b55d7ae2f5b9e386cd1b50a4550696d957cb4900f03a82012708dafc9e1b880fd083b32182b869be8e0922b81f8e175ffde54d797fe11eb03f9e3bf75f1d68bf0b8b6fb4e317a0f9d6f03eaf8ce6675bc60d8c4d90829ce8f72d0163c1d5cf348a862d55063035e7a025f4da968de7e4d7e4004197917f4070f1d6caa02bbebaebb5d7e581e4b66559e635f805ff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	//sgj1117add,cfgvalue:
	var ExtraData []byte
	//0x4bae4a1a12f26e26cac5ce7fba5dce201d24c6ef
	//var faucet = "4bae4a1a12f26e26cac5ce7fba5dce201d24c6ef"

	faucet :=common.HexToAddress("4bae4a1a12f26e26cac5ce7fba5dce201d24c6ef")	//("12345")
	log.Println("HexToAddress 12345 is:%v",faucet)
	ExtraData = append(append(make([]byte, 32), faucet[:]...), make([]byte, crypto.SignatureLength)...)

	log.Println(len(ExtraData))
	log.Println(ExtraData)

	//1117add
	ExtraData2 := append(append(make([]byte, 32), faucet[:]...))
	addressval :=DAODrainList()
	for _,valaddr := range addressval{
		log.Println("cur valaddr  is:%v",valaddr)
		ExtraData2 = append(ExtraData,valaddr[:]...)
	}
	ExtraData2 = append(ExtraData,make([]byte, crypto.SignatureLength)...)
	log.Println(len(ExtraData2))
	log.Println(ExtraData2)

	Extra :=ExtraData2
	//under using
	//Extra := []byte("52657370656374206d7920617574686f7269746168207e452e436172746d616e42eb768f2244c8811c63729a21a3569731535f067ffc57839b00206d1ad20c69a1981b489f772031b279182d99e65703f0076e4812653aab85fca0f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	//("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa")
	//log.Println(len(Extra))
	validatorBytes := Extra[extraVanity : len(Extra)-extraSeal]
	// get validators from headers and use that for new validator set
	newValArr, err := ParseValidators(validatorBytes)
	if err != nil {
		log.Println("ParseValidators err: %v", err)
	}
	log.Println("ParseValidators newValArr: %v", newValArr)
	//1108doing:
	curMySuite := auth.MySuite{}
	curMySuite.TestSign()
	curMySuite.TestKAuth()
	return
	UrlVerify ="http://192.168.1.221:46657/tri_block_validators"

	handlereq.RequestTrustInfo(UrlVerify)
}

func DAODrainList() []common.Address {
	return []common.Address{
		common.HexToAddress("0xd4fe7bc31cedb7bfb8a345f31e668033056b2728"),
		common.HexToAddress("0xb3fb0e5aba0e20e5c49d252dfd30e102b171a425"),
		common.HexToAddress("0x2c19c7f9ae8b751e37aeb2d93a699722395ae18f"),
	}
}

//1119doing:

func NewRPCSuccessResponse(cdc *amino.Codec, id string, res interface{}) pbtrust.RPCResponse {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			log.Println("Error marshalling response,err is:%v",err)
			//return RPCInternalError(id, errors.Wrap(err, "Error marshalling response"))
		}
		rawMsg = json.RawMessage(js)
	}

	return pbtrust.RPCResponse{JSONRPC: "2.0", ID: id, Result: rawMsg, CODE: 0}
}
type SampleResult struct {
	Value string
}
func TestUnmarshallResponses(t *testing.T) {
	//assert := assert.New(t)
	cdc := amino.NewCodec()
	//for _, tt := range responseTests {
		response := &pbtrust.RPCResponse{}
		err := json.Unmarshal([]byte(log.Sprintf(`{"jsonrpc":"2.0","id":%v,"result":{"Value":"hello"}}`, 33)), response)
		log.Println("cur err info is:%s",err)
		a := NewRPCSuccessResponse(cdc, "33", &SampleResult{"hello"})
		log.Println("get response info is:%v",*response, a)
	//}
	response = &pbtrust.RPCResponse{}
	err = json.Unmarshal([]byte(`{"jsonrpc":"2.0","id":true,"result":{"Value":"hello"}}`), response)
	log.Println(err)
}

//results []interface{}
//unmarshalResponseBytesArray
func UnmarshalResponse(cdc *amino.Codec, responseBytes []byte,results interface{}) (interface{}, error) {
	var (
		err       error
		responses pbtrust.RPCResponse
	)
	err = json.Unmarshal(responseBytes, &responses)
	if err != nil {
		return nil, log.Errorf("error unmarshalling rpc response: %v", err)
	}
	if err := cdc.UnmarshalJSON(responses.Result, results); err != nil {
		return nil, log.Errorf("error unmarshalling rpc response result: %v", err)
	}
	return results, nil
}

/*
	bsj := BlockStoreStateJSON{}
	err := cdc.UnmarshalJSON(bytes, &bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal bytes: %X", bytes))
	}*/