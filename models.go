package xmpay

import (
	"encoding/json"
	"net/http"
	"xmpay/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Config struct {
	ApiUrl       string `yaml:"api_url" json:"apiUrl" comment:"API地址"`
	AccessId     string `yaml:"access_id" json:"accessId" comment:"accessId"`
	AccessKey    string `yaml:"access_key" json:"accessKey" comment:"accessKey"`
	InId         string `yaml:"in_id" json:"inId" comment:"收款通道ID"`
	OutId        string `yaml:"out_id" json:"outId" comment:"代付通道ID"`
	InNotifyUrl  string `yaml:"notify_in" json:"inNotifyUrl" comment:"收款回调地址"`
	OutNotifyUrl string `yaml:"notify_out" json:"outNotifyUrl" comment:"代付回调地址"`
}

type PayClient interface {
	CreateVirtual(param *OrderParam) (*pb.VirtualResp, error)
	CreateReceive(param *OrderParam) (*pb.ReceiveResp, error)
	QueryReceive(orderNo, trxNo string) (*pb.OrderQueryResp, error)
	CreateOut(param *OutParam) (*pb.OutResp, error)
	QueryOut(orderNo, trxNo string) (*pb.OrderQueryResp, error)
	Channel() (*pb.ChannelQueryResp, error)
	Balance() (*pb.MerchantBalanceResp, error)
}

type OrderParam struct {
	OrderNo   string `json:"orderNo" validate:"required" comment:"订单号"`
	Ip        string `json:"ip" validate:"required" comment:"用户IP地址"`
	Uid       string `json:"uid"  comment:"用户ID"`
	Name      string `json:"name" validate:"required" comment:"用户姓名"`
	Phone     string `json:"phone" validate:"required" comment:"用户手机号"`
	Email     string `json:"email" validate:"required,email" comment:"用户邮箱"`
	IdNum     string `json:"idNum" validate:"required" comment:"用户证件号码"`
	Pid       int32  `json:"pid" validate:"required" comment:"支付通道ID"`
	NotifyUrl string `json:"notifyUrl" validate:"required,url" comment:"回调地址"`
	Amount    int64  `json:"amount" validate:"required" comment:"交易金额（分）"`
	Subject   string `json:"subject" comment:"商品标题"`
	Body      string `json:"body" comment:"商品描述"`
}
type ReceiveParam struct {
	OrderParam
	ReturnUrl string `json:"returnUrl,omitempty" comment:"付款成功后跳转地址"`
}
type OutParam struct {
	OrderParam
	BankNo   string `json:"bankNo" validate:"required" comment:"银行卡号"`
	BankCode string `json:"bankCode" validate:"required" comment:"银行编号"`
	BankName string `json:"bankName" comment:"银行名称"`
	Mode     string `json:"mode" comment:"付款方式"`
}

type PayClientImpl struct {
	*Config
	aes      *AES
	log      *logrus.Entry
	accessId string
}

type GrpcClient struct {
	PayClientImpl
	conn   *grpc.ClientConn
	client pb.PayServiceClient
}

type HttpClient struct {
	PayClientImpl
	apiUrl string
	client *http.Client
}

func (c *PayClientImpl) Decrypt(body []byte) string {
	decrypt, err := c.aes.Decrypt(body)
	if err != nil {
		return ""
	}

	return string(decrypt)
}

func (c *PayClientImpl) encrypt(param interface{}) *pb.PayRpcParam {

	result := &pb.PayRpcParam{
		AppKey: c.accessId,
	}

	if param == nil {
		return result
	}
	marshal, _ := json.Marshal(param)

	encrypt, err := c.aes.Encrypt(marshal)
	if err != nil {
		return nil
	}
	result.Data = encrypt
	return result
}
