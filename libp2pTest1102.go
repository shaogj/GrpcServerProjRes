package main

import (
	//"20210810BFLProj/CaTriasServerHP0818/mspcert0922/m0ca"
	multiaddr "github.com/multiformats/go-multiaddr"
		//multiaddr "go-multiaddr"

	"fmt"
)

func main() {
	cases := []string{
		"/ip4/1.2.3.4",
		"/ip4/0.0.0.0",
		"/ip6/::1",
		"/ip6/2601:9:4f81:9700:803e:ca65:66e8:c21",
		"/ip6/2601:9:4f81:9700:803e:ca65:66e8:c21/udp/1234/quic",
	}

		for _, a := range cases {
			m, err := multiaddr.NewMultiaddr(a)
			if err != nil {
				fmt.Println(err)
			}
			curinfo := fmt.Sprintf("cur cases addr is:%s:,get m is:%s",a,m)
			fmt.Println(curinfo)
		}
		
		//1103add
			a,err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/utp/tcp/5555/udp/1234/utp/ipfs/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP")
			if err != nil {
				fmt.Println(err)
			}
		curval,err := a.ValueForProtocol(multiaddr.P_IP4)//P_IP6

		curinfo := fmt.Sprintf("cur 'ip4/127.0.0.1/utp',get curval is:%v,err is:%v", curval,err)
		fmt.Println(curinfo)

		//fmt.Println("cur InitMSPManager(),get serializedIdentity is:%v,err is:%v", "serializedIdentity", "err")
	/*
	//1)test good!
	//LoadLocaMSP()
	//ca stdtest good 1
	//MyTestMSPDSSetupNoCryptoConf()
	//curmsppath := "/Users/gejians/20210921DATA-MAC-PC/orderer.example.com"
	//curmsppath := "/Users/gejians/go/src/20210810BFLProj/CaTriasServerHP0818/mspcert0922/sampleconfig/msp"

	//curmsppathtrias := "/Users/gejians/go/src/20210810BFLProj/CaTriasServerHP0818/mspcert0922/trias.orderer.example.config0924/msp"
	//0925sgj
	//curmsppathtriasUserAdmin := "/Users/gejians/go/src/20210810BFLProj/CaTriasServerHP0818/mspcert0922/trias.0923checkingUp_ordererOrganizations_users_Admin/msp"

	//trias.0923checkingUp_peer0.org1.example.com
	curmsppathtriasPeer0Org1 := "/Users/gejians/go/src/20210810BFLProj/CaTriasServerHP0818/mspcert0922/trias.0923checkingUp_peer0.org1.example.com/msp"

	///mspMgr, serializedIdentity, err := m0ca.InitMSPManager(curmsppathtrias, "LocalMSPID007")
	mspMgr, serializedIdentity, err := m0ca.InitMSPManager(curmsppathtriasPeer0Org1, "LocalMSPID037")

	//mspMgr, serializedIdentity, err := m0ca.InitMSPManager(curmsppath,"LocalMSPID")
	if err != nil {
		fmt.Println("Failed to InitMSPManager(),get serializedIdentity is:%v,err is:%v", serializedIdentity, err)
	}
	m0ca.HandleValide(mspMgr, serializedIdentity)

	curmsppathtriasPeer1Org1 := "/Users/gejians/go/src/20210810BFLProj/CaTriasServerHP0818/mspcert0922/trias.0923checkingUp_peer1.org1.example.com/msp"

	///mspMgr, serializedIdentity, err := m0ca.InitMSPManager(curmsppathtrias, "LocalMSPID007")
	mspMgr2, serializedIdentity2, err2 := m0ca.InitMSPManager(curmsppathtriasPeer1Org1, "LocalMSPID038")

	//mspMgr, serializedIdentity, err := m0ca.InitMSPManager(curmsppath,"LocalMSPID")
	if err2 != nil {
		fmt.Println("Failed to InitMSPManager(),get serializedIdentity is:%v,err is:%v", serializedIdentity2, err2)
	}
	//part二.有效性校验
	m0ca.HandleValide(mspMgr2, serializedIdentity2)
	*/
}
