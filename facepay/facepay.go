package facepay

import (
	"encoding/json"
	"fmt"
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/util"
)

const getFacepayURL = "https://payjs.cn/api/facepay"

// Facepay struct
type Facepay struct {
	*context.Context
}

// FacepayRequest 请求参数
type FacepayRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	Openid     string `json:"openid"`       //Y	OPENID
	FaceCode   string `json:"face_code"`    //Y	人脸支付识别码
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// FacepayResponse PayJS返回参数
type FacepayResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功，0:请求失败
	Status       int    `json:"status"`         //N	return_code为0时有status参数为0
	Msg          string `json:"msg"`            //N	return_code为0时返回的错误消息
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo   string `json:"out_trade_no"`   //Y	用户生成的订单号原样返回
	TotalFee     string `json:"total_fee"`      //Y	金额。单位：分
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewFacepay init
func NewFacepay(context *context.Context) *Facepay {
	facepay := new(Facepay)
	facepay.Context = context
	return facepay
}

// GetFacepay
func (facepay *Facepay) GetFacepay(facepayRequest *FacepayRequest) (facepayResponse FacepayResponse, err error) {
	sign := util.Signature(facepayRequest, facepay.Context.Key)
	facepayRequest.Sign = sign
	response, err := util.PostJSON(getFacepayURL, facepayRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &facepayResponse)
	if err != nil {
		return
	}
	if facepayResponse.ReturnCode == 0 {
		err = fmt.Errorf("GetPayQrcode Error , errcode=%d , errmsg=%s", facepayResponse.Status, facepayResponse.Msg)
		return
	}
	// 检测sign
	msgSignature := facepayResponse.Sign
	msgSignatureGen := util.Signature(facepayResponse, facepay.Context.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
