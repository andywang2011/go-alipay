package alipay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleClient_SystemOauthToken() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"u1","access_token":"tok"}}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.Gateway = srv.URL
	api, err := NewClient(cfg)
	if err != nil {
		panic(err)
	}
	c := api.(*Client)
	c.gateway = srv.URL

	out, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "auth-code"})
	if err != nil {
		panic(err)
	}
	fmt.Println(out.AlipaySystemOauthToken.UserId)
	// Output: u1
}
