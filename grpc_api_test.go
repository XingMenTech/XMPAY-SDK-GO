package xmpay

import (
	"fmt"
	"testing"
	"time"

	"xmpay/pb"
)

var grpcClient *GrpcClient

const (
	GrpcApiUrl = "xmpay.xmtest.in:9201"
	HttpApiUrl = "https://xmpay.xmtest.in"
	AppKey     = "pmO14m7sSjHJg7ov"
	AppSecret  = "LrJyqOFlkPwYpP"
)

func init() {
	config := &Config{
		ApiUrl:    GrpcApiUrl,
		AccessId:  AppKey,
		AccessKey: AppSecret,
	}
	var err error
	grpcClient, err = NewGrpcClient(config, nil)
	if err != nil {
		panic(err)
	}
}
func TestGrpcClient_CreateVirtual(t *testing.T) {
	virtualParam.OrderNo = fmt.Sprintf("%s%d", "virtual", time.Now().UnixMilli())
	login, err := grpcClient.CreateVirtual(virtualParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(login)
}

func TestGrpcClient_CreateReceive(t *testing.T) {
	receiveParam.OrderNo = fmt.Sprintf("%s%d", "receive", time.Now().UnixMilli())
	token, err := grpcClient.CreateReceive(receiveParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestGrpcClient_QueryReceive(t *testing.T) {
	menu, err := grpcClient.QueryReceive("receive1761709212", "SKD2910aYSNf0SC")
	if err != nil {
		t.Error(err)
	}
	t.Log(menu)
}
func TestGrpcClient_CreateOut(t *testing.T) {
	outParam.OrderNo = fmt.Sprintf("%s%d", "cash", time.Now().Unix())
	token, err := grpcClient.CreateOut(outParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestGrpcClient_QueryOut(t *testing.T) {
	menu, err := grpcClient.QueryOut("cash1761709269", "DFD2910aaSKyR0D")
	if err != nil {
		t.Error(err)
	}
	t.Log(menu)
}

func TestGrpcClient_Channel(t *testing.T) {
	resp, err := grpcClient.Channel(pb.ORDER_TYPE_RECEIVE)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGrpcClient_Balance(t *testing.T) {
	resp, err := grpcClient.Balance()
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
