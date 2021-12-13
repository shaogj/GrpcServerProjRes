package client_http

import (
	"20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tendermint/go-amino"
	"io/ioutil"
	"net/http"
)

func UnmarshalResponseBytes(cdc *amino.Codec, responseBytes []byte, result interface{}) (interface{}, error) {
	// Read response.  If rpc/core/types is imported, the result will unmarshal
	// into the correct type.
	// log.Notice("response", "response", string(responseBytes))
	var err error
	response := &proto.RPCResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshalling rpc response: %v", err))
	}
	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("response error: %v", response.Error))
	}

	// Unmarshal the RawMessage into the result.
	err = cdc.UnmarshalJSON(response.Result, result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshalling rpc response result: %v", err))
	}
	return result, nil
}

//11.21add
//11119add
func UnmarshalResponseTest(cdc *amino.Codec, responseBytes []byte,results interface{}) (interface{}, error) {

	return nil, nil
}


func ReqGetUrl(vertify string)(respstr string,err error){
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