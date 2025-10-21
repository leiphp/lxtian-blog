# OAuth 社会化登录模块

这是一个通用的 OAuth 2.0 社会化登录客户端库，支持多个主流平台。

## 支持的平台

- ✅ QQ 登录
- ✅ 新浪微博登录
- ✅ GitHub 登录
- ✅ 微信扫码登录

## 快速使用

### 1. 基本用法

```go
package main

import (
    "fmt"
    "lxtian-blog/common/pkg/oauth"
)

func main() {
    // 创建 QQ OAuth 客户端
    qqConfig := oauth.DefaultQQConfig(
        "your_app_id",
        "your_app_secret",
        "http://your-domain.com/callback",
    )
    qqClient := oauth.NewQQClient(qqConfig)
    
    // 获取授权 URL
    authURL := qqClient.GetAuthURL("random_state")
    fmt.Println("请访问:", authURL)
    
    // 用户授权后获取 code，然后获取 access_token
    accessToken, err := qqClient.GetAccessToken("auth_code_from_callback")
    if err != nil {
        panic(err)
    }
    
    // 获取用户信息
    userInfo, err := qqClient.GetUserInfo(accessToken)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("用户信息: %+v\n", userInfo)
}
```

### 2. 在 RPC 服务中使用

已集成到 `rpc/user/internal/logic/user/loginlogic.go`，直接调用 Login RPC 接口即可。

```go
// 调用示例
resp, err := userRpc.Login(ctx, &user.LoginReq{
    LoginType: 1, // 1=QQ, 2=微博, 3=微信, 5=GitHub
    Code:      "oauth_code",
})
```

## 接口说明

### OAuthClient 接口

所有平台客户端都实现了统一的 `OAuthClient` 接口：

```go
type OAuthClient interface {
    // 获取授权URL
    GetAuthURL(state string) string
    
    // 通过code获取access_token
    GetAccessToken(code string) (string, error)
    
    // 获取用户信息
    GetUserInfo(accessToken string) (*OAuthUserInfo, error)
    
    // 刷新token（可选）
    RefreshToken(refreshToken string) (string, error)
}
```

### OAuthUserInfo 结构

```go
type OAuthUserInfo struct {
    OpenID      string `json:"openid"`       // 用户唯一标识
    Nickname    string `json:"nickname"`     // 昵称
    HeadImg     string `json:"head_img"`     // 头像
    Email       string `json:"email"`        // 邮箱（可选）
    UnionID     string `json:"unionid"`      // 联合ID（微信专用）
    AccessToken string `json:"access_token"` // 访问令牌
}
```

## 各平台特点

### QQ

```go
config := oauth.DefaultQQConfig(appID, appSecret, redirectURL)
client := oauth.NewQQClient(config)
```

- 需要两次请求（先获取 OpenID，再获取用户信息）
- 不支持 refresh_token
- Scope: `get_user_info`

### 微博

```go
config := oauth.DefaultWeiboConfig(appID, appSecret, redirectURL)
client := oauth.NewWeiboClient(config)
```

- 需要先获取 UID
- 支持多种头像尺寸
- Scope: `email`

### GitHub

```go
config := oauth.DefaultGithubConfig(clientID, clientSecret, redirectURL)
client := oauth.NewGithubClient(config)
```

- 使用 Bearer Token 认证
- 邮箱需要额外请求 `/user/emails` 接口
- Scope: `user:email`

### 微信

```go
config := oauth.DefaultWechatConfig(appID, appSecret, redirectURL)
client := oauth.NewWechatClient(config)
```

- 支持 UnionID 机制
- 支持 refresh_token
- Scope: `snsapi_login`

## 配置说明

详细配置请参考: `docs/oauth_login_guide.md`

环境变量配置:
```bash
export QQ_APP_ID="..."
export QQ_APP_SECRET="..."
export QQ_REDIRECT_URL="..."

export WEIBO_APP_ID="..."
export WEIBO_APP_SECRET="..."
export WEIBO_REDIRECT_URL="..."

export GITHUB_CLIENT_ID="..."
export GITHUB_CLIENT_SECRET="..."
export GITHUB_REDIRECT_URL="..."

export WECHAT_APP_ID="..."
export WECHAT_APP_SECRET="..."
export WECHAT_REDIRECT_URL="..."
```

## 完整文档

查看完整使用指南: [OAuth Login Guide](../../../docs/oauth_login_guide.md)

## 目录结构

```
common/pkg/oauth/
├── oauth.go      # OAuth 客户端接口定义
├── config.go     # OAuth 配置定义
├── qq.go         # QQ OAuth 客户端实现
├── weibo.go      # 微博 OAuth 客户端实现
├── github.go     # GitHub OAuth 客户端实现
├── wechat.go     # 微信 OAuth 客户端实现
└── README.md     # 本文档
```

## License

MIT

