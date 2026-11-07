package alipay

import "context"

// OauthTokenReq exchanges an authorization code for tokens.
type OauthTokenReq struct {
	Code string // authorization code from user consent
}

// OauthTokenResp is the OAuth token response.
type OauthTokenResp struct {
	AlipaySystemOauthToken struct {
		BaseResp
		SubCode      string `json:"sub_code,omitempty"`
		SubMsg       string `json:"sub_msg,omitempty"`
		UserId       string `json:"user_id,omitempty"`       // Alipay user ID
		AccessToken  string `json:"access_token,omitempty"`  // access token
		ExpiresIn    int    `json:"expires_in,omitempty"`    // access token TTL in seconds
		RefreshToken string `json:"refresh_token,omitempty"` // refresh token
		ReExpiresIn  int    `json:"re_expires_in,omitempty"` // refresh token TTL in seconds
	} `json:"alipay_system_oauth_token_response"`
	Sign string `json:"sign"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/02qkj4
func (o *OauthTokenReq) ApiName() string {
	return "alipay.system.oauth.token"
}

func (o *OauthTokenReq) ApiParam() map[string]string {
	return map[string]string{
		"code":       o.Code,
		"grant_type": "authorization_code",
	}
}

func (c *Client) SystemOauthToken(ctx context.Context, req *OauthTokenReq) (*OauthTokenResp, error) {
	out := &OauthTokenResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipaySystemOauthToken.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
