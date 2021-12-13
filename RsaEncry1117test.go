package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)
var(
	publicKey = "c65f1a13709292eaad6615ff885d3023f4772b5810c112853bfee20b9565d87132c519201c5368a24c452c78c0e29634f4ed4f6bae8f8bae1892858294edc8a2"
	privateKey = "796c823671b118258b53ef6056fd1f9fc96d125600f348f75f397b2000267fe8"
	publicKey2 = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDq8BXSyd1jFgwbXEvy93XlAHUn6Qfi5FlCOA68ZDxdAKKQcBht5eLOSxCQv4hTTFw/WocfMSGiAcQSO/dPAXF8qnMgjy0TKDLhjVgb2q0AsPABZpPLrZ4Gk9BKklY6eDjJzxGoaI+w1fcnnY4lRadDyW+nC2LHfGf3GybkBDQhI18ciWuwRQ6iY32dVgCyXF2ksEPkA2C1pZQy5bsXBNqJvhXAHzNC0/0QZsxUD02oc27dRVtHjHJB8v3x7455gPwyPrMZVep1wGGc7DlrWVsE2OqJrAC4xdJTmodMZ6obFj5mEUkeOthOSCp2tuo6MQ/Tha1Lal+ZcMIkr2hfiwyrhaZ0WGMzazEJC60KPzX1Dq46nB4diphpmI58C1uXGT7RN3ervaTNMlcpGNx9BX7GYc2ET1mLJB2c4R7ICOy5zfcR6VdoL+EjLKyVwt891UYyV9a09L+C8BmXwRa6RU2IgzocZMjq6O+5eJLSiNpcB8f7Yk+UopzItrU9N3YZbcU= gj_shao@126.com`
	privateKey2 = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`
)
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey2))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey2))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
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
	//1117 opencode
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("from CreateKey,get publickey is :%s",hexutil.Encode(publicKeyBytes)[4:])
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priv, address
}
//1117add
func GetPrivPubKey(curpriKeyHash string){
	priKeyHash :=curpriKeyHash
	//priKeyHash := "796c823671b118258b53ef6056fd1f9fc96d125600f348f75f397b2000267fe8"
	// 创建私钥对象，上面私钥没有钱哦
	priKey, err := crypto.HexToECDSA(priKeyHash)
	if err != nil {
		panic(err)
	}
	priKeyBytes := crypto.FromECDSA(priKey)
	fmt.Println("私钥为: %s", hex.EncodeToString(priKeyBytes))

	pubKey := priKey.Public().(*ecdsa.PublicKey)
	// 获取公钥并去除头部0x04
	pubKeyBytes := crypto.FromECDSAPub(pubKey)[1:]
	fmt.Println("公钥为: %s", hex.EncodeToString(pubKeyBytes))

	// 获取地址
	addr := crypto.PubkeyToAddress(*pubKey)
	fmt.Println("地址为: %s", addr.Hex())

}
func main() {

	getPriv,getAddr :=CreateKey()
	fmt.Println("after CreateKey() get getPriv is:%s,getAddr is:%s,len(getAddr) is:%d", getPriv,getAddr,len(getAddr))
	//2teimes:
	GetPrivPubKey(getPriv)
	//1207,,to check http.sign request
		return
	/**/

	data, err := RsaEncrypt([]byte("test"))
	if err != nil {
		panic(err)
	}
	origData, err := RsaDecrypt(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}