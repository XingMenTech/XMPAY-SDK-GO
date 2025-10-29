package xmpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/XingMenTech/XMPAY-SDK-GO/pb"
	"github.com/sirupsen/logrus"
)

const (
	CreateVirtual = "/gateway/api/order/virtual"
	CreateReceive = "/gateway/api/order/receive"
	QueryReceive  = "/gateway/api/order/receive/query"
	CreateOut     = "/gateway/api/order/out"
	QueryOut      = "/gateway/api/order/out/query"
	Channel       = "/gateway/api/channel/query"
	Balance       = "/gateway/api/merchant/balance"
)

func NewHttpClient(config *Config, log *logrus.Entry) *HttpClient {

	if log == nil {
		log = logrus.WithField("model", "HttpClient")
		log.Level = logrus.DebugLevel
	}

	return &HttpClient{
		PayClientImpl: PayClientImpl{
			Config:   config,
			accessId: config.AccessId,
			aes:      NewAES([]byte(config.AccessId), []byte(config.AccessKey)),
			log:      log,
		},
		apiUrl: config.ApiUrl,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
func (c *HttpClient) CreateVirtual(param *OrderParam) (data *pb.VirtualResp, err error) {
	req := &pb.VirtualParam{
		OrderNo:   param.OrderNo,
		Uid:       param.Uid,
		Ip:        param.Ip,
		Email:     param.Email,
		Phone:     param.Phone,
		Name:      param.Name,
		IdNum:     param.IdNum,
		Pid:       param.Pid,
		NotifyUrl: param.NotifyUrl,
	}
	if param.NotifyUrl == "" {
		req.NotifyUrl = c.InNotifyUrl
	}
	if param.Pid <= 0 {
		req.Pid = StringToInt32(c.InId)
	}
	err = c.doRequest(CreateVirtual, http.MethodPost, req, &data)

	return
}

func (c *HttpClient) CreateReceive(param *ReceiveParam) (data *pb.ReceiveResp, err error) {
	req := &pb.ReceiveParam{
		OrderNo:   param.OrderNo,
		Amount:    param.Amount,
		Uid:       param.Uid,
		Ip:        param.Ip,
		Email:     param.Email,
		Phone:     param.Phone,
		Name:      param.Name,
		Subject:   param.Subject,
		Body:      param.Body,
		IdNum:     param.IdNum,
		Pid:       param.Pid,
		NotifyUrl: param.NotifyUrl,
	}
	if param.NotifyUrl == "" {
		req.NotifyUrl = c.InNotifyUrl
	}
	if param.Pid <= 0 {
		req.Pid = StringToInt32(c.InId)
	}
	err = c.doRequest(CreateReceive, http.MethodPost, req, &data)
	return
}

func (c *HttpClient) QueryReceive(orderNo, trxNo string) (data *pb.OrderQueryResp, err error) {

	req := &pb.OrderQueryParam{
		OrderNo:    trxNo,
		MerchantNo: orderNo,
	}
	err = c.doRequest(QueryReceive, http.MethodPost, req, &data)

	return
}

func (c *HttpClient) CreateOut(param *OutParam) (data *pb.OutResp, err error) {

	req := &pb.OutParam{
		OrderNo:   param.OrderNo,
		Amount:    param.Amount,
		Uid:       param.Uid,
		Ip:        param.Ip,
		Email:     param.Email,
		Phone:     param.Phone,
		Name:      param.Name,
		IdNum:     param.IdNum,
		Pid:       param.Pid,
		BankNo:    param.BankNo,
		BankCode:  param.BankCode,
		BankName:  param.BankName,
		Mode:      param.Mode,
		NotifyUrl: param.NotifyUrl,
		Subject:   param.Subject,
		Body:      param.Body,
	}

	if param.NotifyUrl == "" {
		req.NotifyUrl = c.OutNotifyUrl
	}
	if param.Pid <= 0 {
		req.Pid = StringToInt32(c.OutId)
	}
	err = c.doRequest(CreateOut, http.MethodPost, req, &data)

	return
}

func (c *HttpClient) QueryOut(orderNo, trxNo string) (data *pb.OrderQueryResp, err error) {
	req := &pb.OrderQueryParam{
		OrderNo:    trxNo,
		MerchantNo: orderNo,
	}
	err = c.doRequest(QueryOut, http.MethodPost, req, &data)

	return
}
func (c *HttpClient) Channel(orderType pb.ORDER_TYPE) (data []*pb.ChannelQueryResp, err error) {

	param := &pb.ChannelQueryParam{
		OrderType: orderType,
	}

	err = c.doRequest(Channel, http.MethodPost, param, &data)

	return
}

func (c *HttpClient) Balance() (data *pb.MerchantBalanceResp, err error) {

	err = c.doRequest(Balance, http.MethodPost, nil, &data)

	return
}

func (c *HttpClient) doRequest(path, method string, params interface{}, result interface{}) error {
	url := c.apiUrl + path

	requestParam := c.encrypt(params)
	reqParam, _ := json.Marshal(requestParam)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqParam))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Errorf("pay center http request failed , err: %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("network error: (%d)", resp.StatusCode)
		c.log.Error(msg)
		return errors.New(msg)
	}
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("response body read error")
		return errors.New("response body read error")
	}

	c.log.Debug("解码前响应数据：", string(bodyByte))
	var res *pb.PayRpcResp
	err = json.Unmarshal(bodyByte, &res)
	if err != nil {
		c.log.Error("response body unmarshal error")
		return errors.New("response body unmarshal error")
	}
	if res.Code != http.StatusOK {
		c.log.Error(res.Message)
		return errors.New(res.Message)
	}

	decrypt := c.Decrypt([]byte(res.Data))

	c.log.Debug("解码后响应数据：", decrypt)
	err = json.Unmarshal([]byte(decrypt), result)
	if err != nil {
		c.log.Error("response body unmarshal error")
		return err
	}
	return nil
}
