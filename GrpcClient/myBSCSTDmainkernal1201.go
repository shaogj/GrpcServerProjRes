
//package handlereq
///package main//1123bef testing req totalproc good
//1124newint
//package  handlereq
package main

import (
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/client_http"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/handlereq"
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
	cryptoloc "20210810BFLProj/grpcSimpleService1017/crypto"
	cryptoed255519 "20210810BFLProj/grpcSimpleService1017/crypto/ed25519"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-amino"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"errors"
)

var cdc = amino.NewCodec()
const (
	PrivKeyAminoName = "tendermint/PrivKeyEd25519"
	PubKeyAminoName  = "tendermint/PubKeyEd25519"
)

type TrustScore struct{
	TrustPubData string
	Score int
}
type TrustTask struct {
	TrustNodeMap map[string]TrustScore
	//1127update
	//TrustNodeMap map[string]int
	lock sync.Mutex
	NodeTrustInfo  TrustConfig
}
type TrustConfig struct {
	RequestInterval int
	TrustNodeNum	int
	PrivKey  	string
	PublicKey	string
	TrustUrlVerify string
	TrustRankingUrl string
	GRPCListenAddress string
}

//1129add
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

//end1129
func DefaultTrustConfig(trustLocalIp string) *TrustConfig {
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip.String())
	trustLocalIp = ip.String()
	if trustLocalIp == ""{
		trustLocalIp = "127.0.0.1"
	}//1206localkng

	trustUrlVerifyReal := fmt.Sprintf("http://%s:46657/tri_block_validators",trustLocalIp)
	fmt.Println("in DefaultTrustConfig(),cur trustUrlVerifyReal is:%s",trustUrlVerifyReal)

	return &TrustConfig{
		RequestInterval:4,
		TrustNodeNum:3,
		//192.168.1.204
		TrustUrlVerify:"http://192.168.1.221:46657/tri_block_validators",
		//TrustUrlVerify:trustUrlVerifyReal,

		TrustRankingUrl: "http://192.168.1.221:46657/tri_block_info",
		//TrustUrlVerify:"http://http://49.233.196.33:46657/tri_bc_tx_commit",
		//GRPCListenAddress:"tcp://0.0.0.0:36658"
	}
}

//1111.newadd
var inturn = 0
var addressTrusts1 = []string{"keyaddr1","keyaddr2","keyaddr3"}
var addressTrusts2 = []string{"keyaddr4","keyaddr5","keyaddr6","keyaddr7"}


func (this *TrustTask) RequestTrustInfoOld(reqNodeNum int) (int, *proto.TrustDataRespInfo,error) {
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
}
//1127add
func (this *TrustTask) GetTrustNodeScore(trustpubkey string,newscore int) (istrust bool,oldscore,newscoreval int){

	var scoreold int
	var newTrustNodeItem  TrustScore
	for address,trustinfo := range this.TrustNodeMap{
		if trustinfo.TrustPubData == trustpubkey {
			fmt.Println("cur trustpubkey is trustnode,address1201 is :%s ,trustpubkey is :%s,cur trust score is :%d,to set new score is:%d",address,trustpubkey,oldscore,newscore)
			scoreold = trustinfo.Score
			trustinfo.Score = newscore
			//1130update,
			newTrustNodeItem.TrustPubData = trustinfo.TrustPubData
			newTrustNodeItem.Score = newscore
			this.TrustNodeMap[address] = newTrustNodeItem
			fmt.Println("cur trustpubkey TrustNodeMap addr:%s,newscoreval:",address,this.TrustNodeMap[address].Score)
			return true,scoreold,trustinfo.Score
		}
	}
	fmt.Println("cur trustnode list len :%d, no contain cur trustpubkey is:%s ",len(this.TrustNodeMap),trustpubkey)
	return false,0,0

}
//1201add,获取当前节点的可信score
func (this *TrustTask) GetCurNodeScore(trustbscaddr string) int{
	if _, ok := this.TrustNodeMap[trustbscaddr]; ok{
		//get里面的score值是否>0
		curScore := this.TrustNodeMap[trustbscaddr].Score
		fmt.Println("cur address :%s is trusted,cur tmpubkey score is :%V",trustbscaddr,curScore)
		return curScore
	}
	return 0

}
//1201end

//1125add,从TM节点，获取节点的validators列表信息
func (this *TrustTask) RequestTrustInfo(urlTMValidator string) (int, interface{},error) {
	var reqInfo proto.TrustQueryReq
	//1205doing
	//reqInfo.NodeNum = 5
	//UrlVerify := config.GbConf.SettleApiQuery
	var UrlVerify string
	if urlTMValidator == "" {
		//1205--UrlVerify ="http://49.233.196.33:46657/tri_block_validators"
		UrlVerify = "http://192.168.1.221:46657/tri_block_validators"

	}else{
		UrlVerify = urlTMValidator
	}
	resQuerySign := proto.RPCResponse{}
	trustInfo := &proto.ResultValidators{}
	//resQuerySign.Result = &trustInfo

	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo====KKKKK is:%v", UrlVerify, reqInfo)

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

	cdc :=amino.NewCodec()
	cdc.RegisterInterface((*cryptoloc.PubKey)(nil), nil)
	cdc.RegisterConcrete(cryptoed255519.PubKeyEd25519{},
		PubKeyAminoName, nil)

	err = json.Unmarshal(content, &resQuerySign)
	if err != nil {
		fmt.Println("error content form rpc response to resQuerySign,err is: %v", err)
	}
	_,err = client_http.UnmarshalResponseBytes(cdc, content, trustInfo)
	if nil != err {
		fmt.Println("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		fmt.Println("json.Unmarshal err!!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
		return 0, nil,err
	}

	fmt.Println("json.Unmarshal succ!,cur get resQuerySign.Result to gettrustvals.validators is:%v", trustInfo.Validators)
	for key,val := range trustInfo.Validators {
		pubkeyAddr :=val.PubKey.Address()
		//fmt.Println("get cur trustInfo id:%d,Validator.PubKey'Address--444, len is:%d,Address is:%s",key,len(pubkeyAddr),val.PubKey.Address())
		envcodeStr :=hex.EncodeToString([]byte(pubkeyAddr))
		fmt.Println("get cur validators's id:%d, addr len is:%d,Validator.PubKey is:%v,Address--555 is: %s",key,len(pubkeyAddr),val.PubKey,envcodeStr)

	}
	return len(trustInfo.Validators),trustInfo,nil

}

//1125，获取本地TM节点的kernalscore排名分数
func (this *TrustTask) RequestNodeTrustScore(defUrlVerify string) (uint, interface{},error) {
	UrlVerify := defUrlVerify
	if UrlVerify == "" {
		UrlVerify = "http://192.168.1.221:46657/tri_block_info"
	}
	//trustInfo := &proto.ResultValidators{}
	fmt.Println("trustQuery.UrlVerify is:%s,reqInfo is:%v", UrlVerify, "reqInfo")
	trustBlockInfo :=&handlereq.ResultBlock{}
	getblockInfo,err :=client_http.ReqGetUrl(UrlVerify)
	fmt.Println("RequestNodeTrustScore.reqGetUrl(),get getblockInfo is:%s,err is:%v", getblockInfo, err)
	//1125update//client_http
	_,err = client_http.UnmarshalResponseBytes(cdc, []byte(getblockInfo), trustBlockInfo)
	//1119====err = json.Unmarshal(content, &resQuerySign)
	if nil != err {
		fmt.Println("resp=%s,url=%s,err=%v", string(getblockInfo), UrlVerify, err.Error())
		return 0, nil,err
	}else{
		fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Header.Height is:%s,Header.time is:%v,TotalTxs is:%d", trustBlockInfo.BlockMeta.Header.Height,trustBlockInfo.Block.Header.Time,trustBlockInfo.Block.Header.TotalTxs)
		fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Block.ValidatorsHash is:%s", trustBlockInfo.Block.ValidatorsHash)
		fmt.Println("UnmarshalResponseBytes succ!!,cur get trustBlockInfo.Block.trust_data is:%v", trustBlockInfo.Block.TrustData)
		//11.27doing:
		if trustBlockInfo.Block.TrustData == nil {
			fmt.Println("in RequestNodeTrustScore(), cur block height is :%d,get TrustData is nil", trustBlockInfo.BlockMeta.Header.Height)
			return 0,nil,errors.New("cur ResultBlock's TrustData in nil")

		}
		for keyid,pubKey := range *trustBlockInfo.Block.TrustData {
			curTrustPubKey := pubKey.PubKey()
			curscore :=strings.Split(string(pubKey), "/")[1]
			icurscore,_ :=strconv.Atoi(curscore)
			fmt.Println("watching===cur get trustBlockInfo.Block.trust_data's id is:%d,TrustData' pubkey is:%v", keyid,curTrustPubKey)
			isexist,oldscore,newscore := this.GetTrustNodeScore(curTrustPubKey,icurscore)
			if isexist {
				fmt.Println("cur TrustData' pubkey is exist in trustData list,oldscore is:%d,newscore is:%d", oldscore,newscore)
			}else{
				fmt.Println("cur TrustData' pubkey :%s,is exist no in trustData list,it trust newscore is:%d", curTrustPubKey,newscore)

			}
	}



	}
	return 1,trustBlockInfo,nil

}
//end1127
func (this *TrustTask) StartRequest() {
			iserion :=0
			for {
				time.Sleep(time.Second * time.Duration(this.NodeTrustInfo.RequestInterval))
				trustNum,getTrustInfo,err :=this.RequestTrustInfo(this.NodeTrustInfo.TrustUrlVerify)
				//sgj 1115to checking:
				if err != nil  || trustNum < this.NodeTrustInfo.TrustNodeNum{
					fmt.Println("cur RequestTrustInfo() num no matched!,get trustNum is:%d,get err is:%v",trustNum,err)
					//1126,tmp
					//continue
				}
				fmt.Println("cur RequestTrustInfo(),get trustNum is%d,get getTrustInfo is:%v,err is:%v",trustNum,getTrustInfo,err)
				trustnodeaddr := make(map[string]TrustScore,trustNum)

				//1127fix
				if getTrustInfo == nil {
					continue
				}
				getValidators := getTrustInfo.(*proto.ResultValidators)
				var trustpubkey string
				var score int
				for id,nodeaddr:= range getValidators.Validators{
					pubkeyAddr :=nodeaddr.PubKey.Address()
					envcodeStr :=hex.EncodeToString([]byte(pubkeyAddr))
					//1206,to update addr to big byte
					//1206skip--
					srcAddr :=nodeaddr.Address
					decodesrcAddr,decodeerr := hex.DecodeString(envcodeStr)//(string(srcAddr))
					fmt.Println("get cur trustInfo id:%d,Validator.PubKey' Address--999 is :%v, after envcodeStr is:%s,decodesrcAddr is:%s,decodeerr is:%v",id,srcAddr,envcodeStr,decodesrcAddr,decodeerr)

					//1129skip--fmt.Println("get cur trustInfo id:%d,Validator.PubKey is:%s,Address--999 is:%s",id,nodeaddr.PubKey,envcodeStr)
					//1127doing:
					pkCurValidator := fmt.Sprintf("%s", nodeaddr.PubKey)
					if !(strings.HasPrefix(pkCurValidator, "PubKeyEd25519") && len(pkCurValidator) == 79) {
						err = fmt.Errorf("The obtained public key format is incorrect. pubkey=%s", pkCurValidator)
					}
					trustdatapubkey :=fmt.Sprintf("%s/%d", pkCurValidator[14:78], 300)
					parts := bytes.Split([]byte(trustdatapubkey), []byte("/"))
					if len(parts) == 2 {
						trustpubkey = string(parts[0])
						score,_ = strconv.Atoi(string(parts[1]))
					}
					var curTrustScore TrustScore
					curTrustScore.TrustPubData = trustpubkey
					curTrustScore.Score= score
					//1201add
					curscore :=this.GetCurNodeScore(envcodeStr)
					if curscore >0 {
						curTrustScore.Score= curscore
					}
					trustnodeaddr[envcodeStr] = curTrustScore
					fmt.Println("id:%d,cur Validator len(pkCurValidator) is: %d,cur curTrustScore.TrustPubData is :%v",id,len(pkCurValidator),curTrustScore)

				}

				iserion ++
				this.lock.Lock()
				//11.25,返回的可信node列表，比现有的列表少，则设置不可信的节点score为-1，及为不可信的节点
				this.TrustNodeMap = trustnodeaddr
				this.lock.Unlock()
				//to 请求score data：请求应在每个节点出块流程commit后
				retval,resultinfo,err :=this.RequestNodeTrustScore(this.NodeTrustInfo.TrustRankingUrl)
				if err != nil  || retval !=1 {
					log.Println("RequestNodeTrustScore(),resultinfo is:%v,return value is err: %v", resultinfo,err)
				}
			}
}

func NewTrustTask(curinfo *TrustConfig) *TrustTask{
	curTrustTask := &TrustTask{}
	curTrustTask.NodeTrustInfo =  *curinfo
	curTrustTask.TrustNodeMap = make(map[string]TrustScore,curinfo.TrustNodeNum)
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
			//"keyaddr2"
			validaddr1 := "b055c3f23dd27601721b21610b9a7786ce94c893"
			keyaddr2status := this.IsTrustNode(validaddr1)
			keyaddr5status := this.IsTrustNode("keyaddr5")
			fmt.Println("from GetTrustData(), node2 truststatus is:%v,node5 truststatus is:%v",keyaddr2status,keyaddr5status)

		}
	}
}
func (this *TrustTask) IsTrustNode(nodeaddress string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if score, ok := this.TrustNodeMap[nodeaddress]; ok{
		//1125，to 判断里面的score值是否-1，确定是否可信的节点
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

func main() {
	/*
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ip.String())
	return
	*/

	stopchan := make(chan int,1)
	curTrustTask := Start()
	go curTrustTask.GetTrustData()
	<-stopchan
}
func Start() *TrustTask{
	//stopchan := make(chan int,1)
	crustConfig :=DefaultTrustConfig("192.168.1.221")
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
