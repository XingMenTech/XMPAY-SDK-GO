package xmpay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"xmpay/pb"

	"github.com/sirupsen/logrus"
)

var receiveParam = &ReceiveParam{
	OrderParam: OrderParam{
		Ip:        "192.168.1.116",
		Uid:       "1906",
		Name:      "Jack",
		Phone:     "9988776654",
		Email:     "Yukdkskkssls@gmail.com",
		IdNum:     "2503131209376517",
		Pid:       10001,
		NotifyUrl: "https://123456.abc.com/gateway/notify/receive",
		Amount:    10000,
		Subject:   "testReceive",
		Body:      "testReceive",
	},
	ReturnUrl: "",
}

var outParam = &OutParam{
	OrderParam: OrderParam{
		Ip:        "192.168.1.116",
		Uid:       "1907",
		Name:      "Jack kjh",
		Phone:     "9988776654",
		Email:     "Yukdkskkssls@gmail.com",
		IdNum:     "2503131209376517",
		Pid:       10002,
		NotifyUrl: "https://123456.abc.com/gateway/notify/out",
		Amount:    10000,
		Subject:   "testCash",
		Body:      "testCash",
	},
	BankNo:   "8946536458965423",
	BankCode: "37006",
	BankName: "test",
}

var virtualParam = &OrderParam{
	Ip:        "192.168.1.116",
	Uid:       "1906",
	Name:      "Jack kgj",
	Phone:     "9988766754",
	Email:     "Yukdkskls@gmail.com",
	IdNum:     "2503131209376517",
	Pid:       10001,
	NotifyUrl: "https://123456.abc.com/gateway/notify/virtual",
}

var httpClient *HttpClient

func init() {
	config := &Config{
		ApiUrl:    HttpApiUrl,
		AccessId:  AppKey,
		AccessKey: AppSecret,
	}
	httpClient = NewHttpClient(config, logrus.WithField("model", "HttpClient"))
}

func TestHttpClient_CreateVirtual(t *testing.T) {
	virtualParam.OrderNo = fmt.Sprintf("%s%d", "virtual", time.Now().UnixMilli())
	login, err := httpClient.CreateVirtual(virtualParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(login)
}

func TestHttpClient_CreateReceive(t *testing.T) {
	receiveParam.OrderNo = fmt.Sprintf("%s%d", "receive", time.Now().UnixMilli())
	token, err := httpClient.CreateReceive(receiveParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestHttpClient_QueryReceive(t *testing.T) {
	menu, err := httpClient.QueryReceive("2503331339977720", "SKD2810amj0U9i")
	if err != nil {
		t.Error(err)
	}
	t.Log(menu)
}
func TestHttpClient_CreateOut(t *testing.T) {
	outParam.OrderNo = fmt.Sprintf("%s%d", "cash", time.Now().UnixMilli())
	token, err := httpClient.CreateOut(outParam)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestHttpClient_QueryOut(t *testing.T) {
	menu, err := httpClient.QueryOut("2533131279386546", "DFD2910asQwPlKB")
	if err != nil {
		t.Error(err)
	}
	t.Log(menu)
}

func TestHttpClient_Channel(t *testing.T) {

	resp, err := httpClient.Channel(pb.ORDER_TYPE_RECEIVE)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestHttpClient_Balance(t *testing.T) {
	resp, err := httpClient.Balance()
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

type Transfer struct {
	Mode       string `json:"mode"`
	OrderNo    string `json:"orderNo"`
	MerchantId int64  `json:"merchantId"`
	Remark     string `json:"remark"`
	IsRefund   bool   `json:"isRefund"`
}

func TestName(t *testing.T) {
	ts := Transfer{
		Mode:       "3",
		OrderNo:    "DFD1405ayQt6ruB",
		MerchantId: 8,
		Remark:     "",
		IsRefund:   false,
	}
	marshal, err := json.Marshal(ts)

	fmt.Println(string(marshal), err)
}

func TestNotifyVirtual(t *testing.T) {
	param := "{\"amount\":\"100.00\",\"merchant_id\":\"MSF20241226035117000\",\"partner_trx_id\":\"VA_1928012086812151808\",\"payment_element\":\"646687207183486759\",\"plat_fee\":\"7.00\",\"plat_trx_id\":\"VA_1928012086812151808\",\"related_trx_id\":\"VAD2905aabfThd\",\"status\":\"SUCCESS\",\"status_msg\":\"Payment Succeed\",\"success_time\":\"2025-05-29 02:54:37\",\"type\":1,\"signature\":\"doOkPF+Xd31MtPYg0/A7v9xlX5uL3L8xpaBuelHFeajfEU2G8FMhvW9nYZbFyva1fOZ+hYd1KRrKAYhLEaSQPA==\"}"
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9001/gateway/notify/virtual/10004", strings.NewReader(param))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	t.Error(err)
	t.Log(resp)
}

func TestNotifyIn(t *testing.T) {
	param := "{\"amount\":\"50.00\",\"merchant_id\":\"MSF20241226035117000\",\"partner_trx_id\":\"VA_1900512459364073472\",\"payment_element\":\"646310283688173152\",\"plat_fee\":\"7.00\",\"plat_trx_id\":\"VA_1900512459364073472\",\"related_trx_id\":\"VAD1403aYSUHgfD\",\"status\":\"SUCCESS\",\"status_msg\":\"Payment Succeed\",\"success_time\":\"2025-03-14 05:40:55\",\"type\":1,\"signature\":\"amwOp7K8rGEHabbB1g12KPo7j1/nko+IcnHekQH2T5OxV0v3TTBwh9W9Yq33Tld/tAq0+w6v+/H3ARq9igSkbw==\"}"

	req, err := http.NewRequest(http.MethodPost, "http://xmpay-interface.xmtest.in/notify/virtual/10004", strings.NewReader(param))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	t.Error(err)
	t.Log(resp)
}
