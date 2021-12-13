package auth

import (
	//"backend/support/libraries/auth"
	"fmt"
	"io/ioutil"

	//. "gopkg.in/check.v1"
)
type MySuite struct{

}
// func KAuth(addressStr string, body []byte)
func (s *MySuite) TestKAuth() {
	keyfilepath := "./KeyStore/auth.json"
	keyjson, err := ioutil.ReadFile(keyfilepath)
	fmt.Printf("Succeed to read the keyfile at '%s' \n",
		string(keyjson))

	if err != nil {
		fmt.Printf("Failed to read the keyfile at '%s': %v \n",
			keyfilepath, err)
	}

	err = KAuth("0x5ffF03b858D5Fd1633f92d7874DCe20c0D44fAC8",
		"8cb7f6e50c72e39363a14b5edfb2710da53f5ed4ae37d774eb60e40db1638a5425f01dbbb438c86d00c5dc4495313243b7f5f4f387e30c86a71e5e58d4c1a51b00",
		keyjson)
	fmt.Printf("TestKAuth finished ,err is:%v\n",err)
	return
}

func (s *MySuite) TestSign() {
	//keyfilepath := "./KeyStore/auth.json"
	keyfilepath := "/Users/gejians/go/src/20210810BFLProj/grpcSimpleService1017/GrpcClient/KeyStore/auth.json"
	keyjson, _ := ioutil.ReadFile(keyfilepath)
	privk := "8ccced7682da5ca8d94d6f96ecc822ee5b32db047a9e28a51308178ad95ad0e1"
	signStr, _ := KSign(keyjson, privk)
	fmt.Printf("TestSign finished:%s \n", signStr)
}
