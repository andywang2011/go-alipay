package alipay

// MessageTpl is the biz_content template for subscribe messages.
type MessageTpl struct {
	ToUserId       string      `json:"to_user_id"`       // recipient user_id
	UserTemplateId string      `json:"user_template_id"` // template ID
	Page           string      `json:"page"`             // mini-program page to open on tap
	Data           interface{} `json:"data"`             // template variables; see *Tpl types
}

// TestOrderPayTpl is test data for order-paid notification template variables.
type TestOrderPayTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
	Keyword3 struct {
		Value string `json:"value"`
	} `json:"keyword3"`
	Keyword4 struct {
		Value string `json:"value"`
	} `json:"keyword4"`
}

// OrderPayTpl holds order-paid notification template variables.
type OrderPayTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
	Keyword3 struct {
		Value string `json:"value"`
	} `json:"keyword3"`
	Keyword4 struct {
		Value string `json:"value"`
	} `json:"keyword4"`
}

// PickWarnTpl holds pickup reminder template variables.
type PickWarnTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
	Keyword3 struct {
		Value string `json:"value"`
	} `json:"keyword3"`
}

// TestOrderChangeTpl is test data for order-status change template variables.
type TestOrderChangeTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
}

// OrderChangeTpl holds order-status change template variables.
type OrderChangeTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
	Keyword3 struct {
		Value string `json:"value"`
	} `json:"keyword3"`
	Keyword4 struct {
		Value string `json:"value"`
	} `json:"keyword4"`
}

// RefundDoneTpl holds refund-completed template variables.
type RefundDoneTpl struct {
	Keyword1 struct {
		Value string `json:"value"`
	} `json:"keyword1"`
	Keyword2 struct {
		Value string `json:"value"`
	} `json:"keyword2"`
	Keyword3 struct {
		Value string `json:"value"`
	} `json:"keyword3"`
	Keyword4 struct {
		Value string `json:"value"`
	} `json:"keyword4"`
	Keyword5 struct {
		Value string `json:"value"`
	} `json:"keyword5"`
}
