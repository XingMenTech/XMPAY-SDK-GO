package xmpay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"xmpay/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	fmt.Printf("Starting RPC %s \n", method)             // 在调用之前记录日志
	err := invoker(ctx, method, req, reply, cc, opts...) // 调用原始方法
	if err != nil {                                      // 处理错误和在调用之后记录日志
		fmt.Printf("Error in RPC %s: %v \n", method, err)
		return err
	}
	fmt.Printf("Finished RPC %s \n", method) // 在调用之后记录日志
	return nil
}

// NewGrpcClient 创建一个新的gRPC客户端
func NewGrpcClient(config *Config, log *logrus.Entry) (*GrpcClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(GrpcClientInterceptor),
	}

	conn, err := grpc.NewClient(config.ApiUrl, opts...)

	if err != nil {
		return nil, err
	}

	if log == nil {
		log = logrus.WithField("model", "HttpClient")
		log.Level = logrus.DebugLevel
	}

	return &GrpcClient{
		PayClientImpl: PayClientImpl{
			accessId: config.AccessId,
			aes:      NewAES([]byte(config.AccessId), []byte(config.AccessKey)),
			log:      log,
		},
		conn:   conn,
		client: pb.NewPayServiceClient(conn),
	}, nil
}

// Close 关闭gRPC连接
func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

// SayHello 调用SayHello RPC方法
func (c *GrpcClient) CreateVirtual(param *OrderParam) (*pb.VirtualResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.VirtualParam{
		OrderNo:   param.OrderNo,
		Uid:       param.Uid,
		Ip:        param.Ip,
		Email:     param.Email,
		Phone:     param.Phone,
		Name:      param.Name,
		Pid:       param.Pid,
		IdNum:     param.IdNum,
		NotifyUrl: param.NotifyUrl,
	}

	if param.NotifyUrl == "" {
		req.NotifyUrl = c.InNotifyUrl
	}
	if param.Pid <= 0 {
		req.Pid = StringToInt32(c.InId)
	}

	rpcParam := c.encrypt(req)
	resp, err := c.client.VirtualAccount(ctx, rpcParam)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}

	respData := c.Decrypt([]byte(resp.Data))
	var user pb.VirtualResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *GrpcClient) CreateReceive(param *ReceiveParam) (*pb.ReceiveResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.ReceiveParam{
		OrderNo:   param.OrderNo,
		Amount:    param.Amount,
		Uid:       param.Uid,
		Ip:        param.Ip,
		Email:     param.Email,
		Phone:     param.Phone,
		Name:      param.Name,
		Pid:       param.Pid,
		IdNum:     param.IdNum,
		NotifyUrl: param.NotifyUrl,
		ReturnUrl: param.ReturnUrl,
		Subject:   param.Subject,
		Body:      param.Body,
	}
	if param.NotifyUrl == "" {
		req.NotifyUrl = c.InNotifyUrl
	}
	if param.Pid <= 0 {
		req.Pid = StringToInt32(c.InId)
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.Receive(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.Decrypt([]byte(resp.Data))
	var user pb.ReceiveResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *GrpcClient) QueryReceive(orderNo, trxNo string) (*pb.OrderQueryResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.OrderQueryParam{
		OrderNo:    trxNo,
		MerchantNo: orderNo,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.ReceiveQuery(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.Decrypt([]byte(resp.Data))
	var user pb.OrderQueryResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *GrpcClient) CreateOut(param *OutParam) (*pb.OutResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	rpcParam := c.encrypt(req)
	resp, err := c.client.Out(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.Decrypt([]byte(resp.Data))
	var menu *pb.OutResp
	if err := json.Unmarshal([]byte(respData), &menu); err != nil {
		return nil, err
	}
	return menu, nil
}
func (c *GrpcClient) QueryOut(orderNo, trxNo string) (*pb.OrderQueryResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.OrderQueryParam{
		OrderNo:    trxNo,
		MerchantNo: orderNo,
	}

	rpcParam := c.encrypt(req)
	resp, err := c.client.OutQuery(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.Decrypt([]byte(resp.Data))
	var menu *pb.OrderQueryResp
	if err := json.Unmarshal([]byte(respData), &menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (c *GrpcClient) Channel(orderType pb.ORDER_TYPE) ([]*pb.ChannelQueryResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcParam := c.encrypt(&pb.ChannelQueryParam{
		OrderType: orderType,
	})
	resp, err := c.client.ChannelQuery(ctx, rpcParam)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}

	respData := c.Decrypt([]byte(resp.Data))
	var user []*pb.ChannelQueryResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *GrpcClient) Balance() (*pb.MerchantBalanceResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcParam := c.encrypt(nil)
	resp, err := c.client.MerchantBalance(ctx, rpcParam)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}

	respData := c.Decrypt([]byte(resp.Data))
	var user pb.MerchantBalanceResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
