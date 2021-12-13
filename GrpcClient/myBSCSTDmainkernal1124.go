
//package handlereq
///package main//1123bef testing req totalproc good
//1124newint
//package  handlereq
package main
import (
	auth "20210810BFLProj/grpcSimpleService1017/GrpcClient/KeyStore"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/handlereq"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type TrustTask struct {
	TrustNodeMap map[string]int
	lock sync.Mutex
	NodeTrustInfo  TrustConfig
}
type TrustConfig struct {
	RequestInterval int
	TrustNodeNum	int
	PrivKey  	string
	PublicKey	string
	TrustUrlVerify string
	GRPCListenAddress string
}

func DefaultTrustConfig() *TrustConfig {
	return &TrustConfig{
		RequestInterval:4,
		TrustNodeNum:3,
		//TrustUrlVerify:"http://192.168.1.221:46657/tri_bc_tx_commit",
		//1123doing
		TrustUrlVerify:"http://http://49.233.196.33:46657/tri_bc_tx_commit",
		//GRPCListenAddress:"tcp://0.0.0.0:36658"
	}
}
//1111.newadd
var inturn = 0
var addressTrusts1 = []string{"keyaddr1","keyaddr2","keyaddr3"}
var addressTrusts2 = []string{"keyaddr4","keyaddr5","keyaddr6","keyaddr7"}

func (this *TrustTask) RequestTrustInfo(reqNodeNum int) (int, *proto.TrustDataRespInfo,error) {
	var reqInfo proto.TrustQueryReq
	//sgj1115trying
	curRespInfo := &proto.TrustDataRespInfo{}
	if inturn == 0 {
		inturn = 1
		curRespInfo.Address = addressTrusts1
	}else{
		inturn = 0
		curRespInfo.Address = addressTrusts2
	}

	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo is:%v,RespInfo' addrs is %v", this.NodeTrustInfo.TrustUrlVerify, reqInfo,curRespInfo.Address)

	return len(curRespInfo.Address),curRespInfo,nil

	//end 1115
	reqInfo.TxSignReq = "signstr123===004"
	//UrlVerify := "http://192.168.1.221:46657/tri_bc_tx_commit"

	resQuerySign := proto.RPCResponse{}
	trustInfo := &proto.TrustDataRespInfo{}
	//resQuerySign.Data = trustInfo
	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo is:%v", this.NodeTrustInfo.TrustUrlVerify, reqInfo)

	reqBody, err := json.Marshal(&reqInfo)
	if nil != err {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return 0, nil,nil
	}

	//PrivKey :="56b37a7ecc60166b00c46de086722c0e9f52b5437cd6991a19014b2528deb601"
	//AccessKeyAddr := "0x3D9C84cDe940c63B99722630EE4F1FCA5917e0b9"
	var signData string
	if signData, err = auth.KSign(reqBody, this.NodeTrustInfo.PrivKey); err != nil {
		fmt.Println("In trustQuery(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return 0,nil, nil
	}else{
		fmt.Println("In trustQuery(),auth.KSign succ!,get signData is :%v,err is:%v", signData, err)
	}
	/*
		err = auth.KAuth(AccessKeyAddr, signData, reqBody)
		fmt.Println("In trustQuery(),auth.KAuth finished!,get err is:%v", err)
		return 0,nil, nil
	*/
	reader := bytes.NewReader(reqBody)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", this.NodeTrustInfo.TrustUrlVerify, reader) // URL-encoded payload
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
	err = json.Unmarshal(content, &resQuerySign)
	if nil != err {
		fmt.Println("resp=%s,url=%s,err=%v", string(content), this.NodeTrustInfo.TrustUrlVerify, err.Error())
		return 0, nil,err
	}
	fmt.Println("json.Unmarshal succ!,cur get resQuerySign :%v,get resQuerySign.Result is :%v", resQuerySign, resQuerySign.Result)
	/*1122try request interval routing
	err = json.Unmarshal(resQuerySign.Result, trustInfo)
	//err = cdc.UnmarshalJSON(resQuerySign.Result, trustInfo)

	if nil != err {
		fmt.Println("resp's result is:=%s,err=%v", string(resQuerySign.Result), err.Error())
		return 0, nil,err
	}
	fmt.Println("json.Unmarshal succ!,cur get trustInfo info is :%v",trustInfo)
	*/
	return 1,trustInfo,nil

}
func (this *TrustTask) StartRequest() {
			iserion :=0
			for {
				time.Sleep(time.Second * time.Duration(this.NodeTrustInfo.RequestInterval))
				trustNum,getTrustInfo,err :=this.RequestTrustInfo(this.NodeTrustInfo.TrustNodeNum)
				//1124doing:
				retval,resultinfo,err :=handlereq.RequestNodeTrustScore("tri_block_info")
				if err != nil  || retval !=1 {
					log.Println("RequestNodeTrustScore(),resultinfo is:%v,return value is err: %v", resultinfo,err)
				}
				//if err != nil  || trustNum != this.NodeTrustInfo.TrustNodeNum{
				//sgj 1115to checking:
				if err != nil  || trustNum < this.NodeTrustInfo.TrustNodeNum{
					fmt.Println("cur RequestTrustInfo() failed,get trustNum is:%d,get err is:%v",trustNum,err)
					continue
				}
				fmt.Println("cur RequestTrustInfo(),get trustNum is%d,get getTrustInfo is:%v,err is:%v",trustNum,getTrustInfo,err)
				trustnodeaddr := make(map[string]int,trustNum)
				for _,nodeaddr:= range getTrustInfo.Address{
					trustnodeaddr[nodeaddr] = iserion
				}
				iserion ++

				this.lock.Lock()
				this.TrustNodeMap = trustnodeaddr
				this.lock.Unlock()

			}
}

func NewTrustTask(curinfo *TrustConfig) *TrustTask{
	curTrustTask := &TrustTask{}
	curTrustTask.NodeTrustInfo =  *curinfo
	curTrustTask.TrustNodeMap = make(map[string]int,curinfo.TrustNodeNum)
	return curTrustTask
}

func (this *TrustTask) GetTrustData() {

	t:=time.NewTicker(time.Second *2)

	for {
		select {
		case <-t.C:
			this.lock.Lock()
			getnum := len(this.TrustNodeMap)
			this.lock.Unlock()
			fmt.Println("time interval,GetTrustData' maplen is:%d",getnum)
			keyaddr2status := this.IsTrustNode("keyaddr2")
			keyaddr5status := this.IsTrustNode("keyaddr5")
			fmt.Println("from GetTrustData(), node2 truststatus is:%v,node5 truststatus is:%v",keyaddr2status,keyaddr5status)

		}
	}
}
func (this *TrustTask) IsTrustNode(nodeaddress string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if score, ok := this.TrustNodeMap[nodeaddress]; ok{
		fmt.Println("cur address :%s is trusted,score is :%V",nodeaddress,score)
		return true
	} else {
		fmt.Println("cur address :%s is no trusted,score is :%V",nodeaddress,score)
		//fmt.Println("cur address :%s is no trusted")
		return false
	}
}
//if number%p.config.Epoch == 0 {
func (this *TrustTask) GetCurrentTrustValidators(validatornum int) (nodeaddrs []string,err error){
	//add by scores
	var validatorAddr []string
	this.lock.Lock()
	defer this.lock.Unlock()
	for addrkey,_ := range this.TrustNodeMap{
		validatorAddr = append(validatorAddr,addrkey)
	}
	return validatorAddr,nil
}

/**/
//1129doing pm Last Temping:
func main() {
	stopchan := make(chan int,1)
	curTrustTask := Start()
	go curTrustTask.GetTrustData()
	<-stopchan
}


func Start() *TrustTask{
	//stopchan := make(chan int,1)
	crustConfig :=DefaultTrustConfig()
	curTrustTask:= NewTrustTask(crustConfig)
	go curTrustTask.StartRequest()
	//go curTrustTask.GetTrustData()
	//<-stopchan
	return curTrustTask
}
/*
1.初始配置替换掉ParseValidators():
	validatorBytes := checkpoint.Extra[extraVanity : len(checkpoint.Extra)-extraSeal]
				// get validators from headers
				validators, err := ParseValidators(validatorBytes)
				
2.
替换掉getCurrentValidators
func (p *Parlia) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
if number%p.config.Epoch == 0 {
		newValidators, err := p.getCurrentValidators(header.ParentHash)
		if err != nil {
			return err
		}
				
*/
