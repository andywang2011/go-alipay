package alipay

import "context"

// LocalitemQueryReq is the request to query a local mini-program item.
type LocalitemQueryReq struct {
	BizContent string // biz_content JSON; see LocalitemQueryParam
}

// LocalitemQueryParam holds biz_content for local item query.
type LocalitemQueryParam struct {
	ItemId      string `json:"item_id"`                 // Alipay item ID
	OutItemId   string `json:"out_item_id,omitempty"`   // merchant item ID
	NeedEditSpu string `json:"need_edit_spu,omitempty"` // query draft/pending SPU when "1"
}

// LocalitemQueryResp is the response for local item query.
type LocalitemQueryResp struct {
	AlipayOpenAppLocalitemQueryResponse struct {
		BaseResp
		CategoryId           string   `json:"category_id"`             // platform category ID
		ItemType             string   `json:"item_type"`               // item type
		OutItemId            string   `json:"out_item_id"`             // merchant item ID
		ItemId               string   `json:"item_id"`                 // Alipay item ID
		IsOnline             bool     `json:"is_online"`               // whether this is the online version
		Title                string   `json:"title"`                   // item title
		MerchantName         string   `json:"merchant_name"`           // merchant name
		HeadImg              string   `json:"head_img"`                // cover image URL
		ImageList            []string `json:"image_list"`              // image URLs
		Path                 string   `json:"path"`                    // item detail page path
		ItemDetailsPageModel string   `json:"item_details_page_model"` // consumer detail page mode
		StockNum             int      `json:"stock_num"`               // stock quantity
		SpuStatus            string   `json:"spu_status"`              // SPU status
		Attrs                []struct {
			AttrKey   string `json:"attr_key"`   // attribute name
			AttrValue string `json:"attr_value"` // attribute value
		} `json:"attrs"`
	} `json:"alipay_open_app_localitem_query_response"`
	Sign string `json:"sign"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/dbbe065d_alipay.open.app.localitem.query
func (o *LocalitemQueryReq) ApiName() string {
	return "alipay.open.app.localitem.query"
}

func (o *LocalitemQueryReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) LocalItemQuery(ctx context.Context, req *LocalitemQueryReq) (*LocalitemQueryResp, error) {
	out := &LocalitemQueryResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipayOpenAppLocalitemQueryResponse.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
