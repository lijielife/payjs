package native

import (
	"encoding/json"
	"fmt"
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/util"
)

const getPayQrcodeURL = "https://payjs.cn/api/native"

// Native struct
type Native struct {
	*context.Context
}

// PayQrcodeRequest 请求参数
type PayQrcodeRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	NotifyUrl  string `json:"notify_url"`   //N	接收微信支付异步通知的回调地址。必须为可直接访问的URL，不能带参数、session验证、csrf验证。留空则不通知
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// PayQrcodeResponse PayJS返回参数
type PayQrcodeResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功，0:请求失败
	Status       int    `json:"status"`         //N	return_code为0时有status参数为0
	Msg          string `json:"msg"`            //N	return_code为0时返回的错误消息
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo   string `json:"out_trade_no"`   //Y	用户生成的订单号原样返回
	TotalFee     int    `json:"total_fee"`      //Y	金额。单位：分
	Qrcode       string `json:"qrcode"`         //Y	二维码图片地址
	CodeUrl      string `json:"code_url"`       //Y	可将该参数生成二维码展示出来进行扫码支付
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewNative init
func NewNative(context *context.Context) *Native {
	native := new(Native)
	native.Context = context
	return native
}

// GetPayQrcode 请求PayJS获取支付二维码
func (native *Native) GetPayQrcode(totalFeeReq int, bodyReq, outTradeNoReq, attachReq string) (outTradeNoResp string, totalFeeResp int, qrcodeResp, codeUrlResp, payJSOrderIDResp string, err error) {
	payQrcodeRequest := PayQrcodeRequest{
		MchID:      native.MchID,
		TotalFee:   totalFeeReq,
		OutTradeNo: outTradeNoReq,
		Body:       bodyReq,
		Attach:     attachReq,
		NotifyUrl:  native.NotifyUrl,
	}
	sign := util.Signature(payQrcodeRequest, native.Key)
	payQrcodeRequest.Sign = sign
	response, err := util.PostJSON(getPayQrcodeURL, payQrcodeRequest)
	if err != nil {
		return
	}

	payQrcodeResponse := PayQrcodeResponse{}
	err = json.Unmarshal(response, &payQrcodeResponse)
	if err != nil {
		return
	}
	if payQrcodeResponse.ReturnCode != 1 {
		err = fmt.Errorf("GetPayQrcode Error , errcode=%v , errmsg=%s, errmsg=%s", payQrcodeResponse.ReturnCode, payQrcodeResponse.Msg, payQrcodeResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := payQrcodeResponse.Sign
	msgSignatureGen := util.Signature(payQrcodeResponse, native.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}

	outTradeNoResp = payQrcodeResponse.OutTradeNo
	totalFeeResp = payQrcodeResponse.TotalFee
	qrcodeResp = payQrcodeResponse.Qrcode
	codeUrlResp = payQrcodeResponse.CodeUrl
	payJSOrderIDResp = payQrcodeResponse.PayJSOrderID
	return
}
