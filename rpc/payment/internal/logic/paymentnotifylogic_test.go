package logic

import (
	"context"
	"encoding/base64"
	"fmt"
	"lxtian-blog/common/pkg/initdb"
	"net/url"
	"strings"
	"testing"

	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/rpc/payment/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestVerifySign 本地调试 verifySign 方法
//
// 使用方法：
//  1. 修改下面的 testData 和 testSign 为你的实际数据
//  2. 如果使用真实客户端，设置 useRealClient = true 并配置 alipayConfig
//  3. 运行测试:
//     go test -v -run TestVerifySign ./rpc/payment/internal/logic
//     或者
//     go test -v -run TestVerifySign rpc/payment/internal/logic
//
// 注意事项：
// - 如果 useRealClient = false，只会测试 buildSignContent 函数的逻辑（不进行真实签名验证）
// - 如果 useRealClient = true，需要配置正确的支付宝公钥等信息
// - testData 应该包含支付宝通知的原始数据（URL编码格式）
// - testSign 应该从支付宝通知中提取的 sign 参数值
func TestVerifySign(t *testing.T) {
	// ========== 配置区域：修改下面的数据为你的实际测试数据 ==========

	// 支付宝通知数据（原始数据，可能包含 sign 和 sign_type，这些会被自动过滤）
	testData := `gmt_create=2025-11-11+23%3A46%3A36&charset=utf-8&gmt_payment=2025-11-11+23%3A46%3A43&notify_time=2025-11-13+00%3A15%3A25&subject=%E4%BC%9A%E5%91%98%E5%A5%97%E9%A4%90%E8%B4%AD%E4%B9%B0+-+%E6%9C%88%E5%BA%A6%E4%BC%9A%E5%91%98&sign=PRI0mL9vNaR0XekR952yzu4JKL6lWVFojax%2Bl6PR4%2FCiWyD%2F5j3QDTCEcoECMkaLR7jQgsWQeMEBXou1fvmTjOgrpfKUn2P6WqJDH%2BjCE7WPxG6Lai%2FrMp%2Fj3vi7IoT3DpEHyVPqpCRTqg22X1K53C0v9jF9IkmgLHQ675sbJzCRwIw7iZaa%2Bk77gpgYPK97nO4KePskrQPEA724szeCfcBW847yoJ8Kdo6lcpRPcHz6%2BR4MR8aNUOuvhhDVGr1c2clVqFUvxp4znQqBkq9muHsHRy%2FDnl5NHy7hY6UV6yaqztaWCOBLc%2FVJH6oepikmdEHT6zRny8vMCteo%2BaVufQ%3D%3D&buyer_id=2088222260960342&body=%E6%9C%88%E5%BA%A6%E4%BC%9A%E5%91%98+-+%E9%80%82%E5%90%88%E7%9F%AD%E6%9C%9F%E4%BD%93%E9%AA%8C%E7%94%A8%E6%88%B7&invoice_amount=0.01&version=1.0&notify_id=2025111101222234643060341413047947&fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D&notify_type=trade_status_sync&out_trade_no=NO77356367081373703&total_amount=0.01&trade_status=TRADE_SUCCESS&trade_no=2025111122001460341440512041&auth_app_id=2016110102454851&receipt_amount=0.01&point_amount=0.00&buyer_pay_amount=0.01&app_id=2016110102454851&sign_type=RSA2&seller_id=2088012021942135`
	// 签名（从支付宝通知中提取的 sign 参数值）
	//testSign := "rysqMtOkc6HpJdFS3JEZ/3DKzdclldm/KpjD RQv6i7B/2/mVlNlMM4V7q56HH4LVWCHl9rWH8KjyRnyQTxQyFszyq0TVVxXLAqCUdP/x3I8TxA5MKhDPULssOwi5qVzFoGoi3JrHQO9xesZrWVzf1Gs0ct4qHt0y/rNDbfepJiyuuUmy9x5i5UBgNKMuIizVjTHvUbMm2e4MBSafPbrPzPjAFJOJxEzCNopXfZzN417j31du0F0QFp07jrT6ydgeQpx9BZjHTO2eQXtV/62fpO21pqX5Dfv9e5Q5tEzQOgktANhxUNhdysSMjf99HdxcDvjNHMgsmK80QWxvIYkeQ=="
	var testSign string
	// 是否使用真实的支付宝客户端（需要配置下面的 alipayConfig）
	useRealClient := true

	// 支付宝客户端配置（如果需要使用真实客户端，请填写下面的配置）
	var alipayConfig *alipay.AlipayConfig
	if useRealClient {
		alipayConfig = &alipay.AlipayConfig{
			AppId:           "2016110102454851",
			AppPrivateKey:   "MIIEpAIBAAKCAQEAuTV3koMcBir0K3QRCyuiLTVISu+1yxV5r4SFo43iT1DoRRQYOsNY5rx52xmfC4hOCJ1xdFjKjhzoS6ma+9eHjOKZTo8CtnmAMmJhh0Xu9sdmlQVWLbdptN7kfPP2qCiPwCYdX7UGslhER17l4z/FxsvNUdlZa5NoN7Ugyf+DnfM4QSIuYi5mfuhmQZHS1wcgAnyNlyjRgJbahvyp1XCRQCUKPONrTtpat9V64PX/M6uddOH7wT15347lNHSkJasCUAyGh9FFKAdoHb+/B+RN6HYWABR4izJILoCrQP6NdbUc6sjRz7voc55dekDqFZIcM5ixxotd2rqub8zzK3VTHQIDAQABAoIBAEV4Hp+f+fT+S5O492OfPDeE0tb4ztGb/oatSIsufwKNMHIotWXlPAVgELz0AUoMGGj21UV0wJVJloA6390y3K8ll6d10Ois7j730+kvfBCofnvLAqYnM8kaCc4snAo7HKBQK5hoiFiA7ytuFwCEPSTx8NOQmQ/WvLKYh/H8m1u06bsVyEZfTYm6xlhmccynNpgw3RugEy02pVaf/uVJq15TwqdTadlzO57iVOSTHyKErTGF5l9VCYm0ujCp/5mID2b0ImLoQX7NS5epcMBeEcgExH49TYxy10UyArRTMbp3BGSW8YLnaEaW9bGcreKZrynIiCEf3gu60IvUHRCEF+kCgYEA3ANHAlQO8Ry6P1fs/dkU0+tqgEsA9x2l58Tvqnf+zwVXQma8zDfKgapzlEAjUHq29OR8rZBXBZH9Amt0DeIS0d9fWVfzzrUUdDYP1hkAl8hw3rNU8AYQiCmo/1H4JpPf+8/SE21r0i6CmFHeNUJof945GyBOutfkTPBq228V/tMCgYEA14DS85vnuRsxwSeJhiqczhszPRfdTOhZFXHTyY1ilsLc7vDPcWedQrEiddrJsA8qcKNb3w86IhYPUI5uaVz816ja5sF9ULfR35Znka+Fe54b3y+R5pqUtEAl4frAmxok9H7M8Lo9eZ7oMkq05lynf8M+3kKtiNWMQA/B+gc3kE8CgYEAyU3GXv8CXOJoFyUgFnPVdsFjxNbbnz9lWVb74wHABzNfz8Wo4UH67AFFl1PH/A8L765P1Y7H0LTuxpQCr+E2TwkOePTcgzlz6ZC9lKtzu20OuPVktekWnz9e/Z3Ga6XJvuE72cK4cKtVmoDty9VjP/vYTFWXM6XtoegoHXbarTkCgYEAjWOSBA7H66Sx8h50ljgjBP7HkU+0/B59RBqYb2Z5xpw2w/XuxGLMxNLe3yAar45js98aCbE93NtIVPv96Nb/dKbuZ/OOuoTAB8fwT58vHrnPY5EcUoYdBl4H/Mm90IVItbjz0QUADGl7wnNNWM51ftekycJJhLtG90jfZaGSjPUCgYBaocx1OCnlYMZaLKt+tZTiUHUnfNfMi91VHQSea+8Fp2YRFLmd6vab9ZVhmuqU/rjsFu0rCUAUUKq9ahx4Jc1ZMHFiqhJ1nfSQhIbhB2/OTjjjZWL+irPbnJxBO2lVCqiuLaAp2AwJcA6P2OXLs2DWMA6PSmIH3N1JY1mswoQ/bA==",
			AlipayPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuTV3koMcBir0K3QRCyuiLTVISu+1yxV5r4SFo43iT1DoRRQYOsNY5rx52xmfC4hOCJ1xdFjKjhzoS6ma+9eHjOKZTo8CtnmAMmJhh0Xu9sdmlQVWLbdptN7kfPP2qCiPwCYdX7UGslhER17l4z/FxsvNUdlZa5NoN7Ugyf+DnfM4QSIuYi5mfuhmQZHS1wcgAnyNlyjRgJbahvyp1XCRQCUKPONrTtpat9V64PX/M6uddOH7wT15347lNHSkJasCUAyGh9FFKAdoHb+/B+RN6HYWABR4izJILoCrQP6NdbUc6sjRz7voc55dekDqFZIcM5ixxotd2rqub8zzK3VTHQIDAQAB",
			GatewayUrl:      "https://openapi.alipay.com/gateway.do",
			NotifyUrl:       "https://gw.100txy.com/api/payment/notify",
			ReturnUrl:       "https://www.100txy.com/api/payment/return",
			IsProd:          true, // 测试环境用 false，生产环境用 true
			SignType:        "RSA2",
			Charset:         "utf-8",
			Format:          "JSON",
			Version:         "1.0",
			Timeout:         "30m", // 注意：Timeout 是 string 类型
		}
	}
	// ========== 配置区域结束 ==========

	// 初始化日志
	logx.DisableStat()

	// 初始化数据库
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"lxt_blog",
		"YXDpHdrPLpFTBnp6",
		"111.229.27.82",
		"3306",
		"lxt_blog",
	)
	// 初始化GORM数据库
	db := initdb.InitDB(dataSource)

	// 如果使用真实客户端
	if useRealClient && alipayConfig != nil {
		alipayClient, err := alipay.NewAlipayClient(alipayConfig)
		if err != nil {
			t.Fatalf("Failed to create alipay client: %v", err)
		}

		// 创建 ServiceContext
		svcCtx := &svc.ServiceContext{
			AlipayClient: alipayClient,
			DB:           db,
		}

		// 创建 context
		ctx := context.Background()

		// 创建 PaymentNotifyLogic
		logic := NewPaymentNotifyLogic(ctx, svcCtx)

		// 如果 testSign 为空，尝试从 testData 中提取 sign
		if testSign == "" || testSign == "替换为实际的签名值" {
			sign := extractSignFromData(testData)
			if sign != "" {
				testSign = sign
				t.Logf("从数据中提取到签名: %s (前50字符)", safeSubstring(testSign, 50))
			}
		}

		// 如果仍然没有签名，报错
		if testSign == "" || testSign == "替换为实际的签名值" {
			t.Fatal("testSign 未设置，请手动设置 testSign 或在 testData 中包含 sign 参数")
		}

		t.Logf("=== 开始测试 verifySign (使用真实客户端) ===")
		t.Logf("测试数据长度: %d", len(testData))
		t.Logf("测试签名长度: %d", len(testSign))

		// 先测试 buildSignContent，查看构建的签名字符串
		signContent, err := buildSignContent(testData)
		if err != nil {
			t.Fatalf("buildSignContent 失败: %v", err)
		}
		t.Logf("构建的签名字符串: %s", signContent)
		t.Logf("构建的签名字符串长度: %d", len(signContent))

		// 打印签名的详细信息
		t.Logf("签名长度: %d", len(testSign))
		if len(testSign) > 100 {
			t.Logf("签名（前100字符）: %s", testSign[:100])
			t.Logf("签名（后100字符）: %s", testSign[len(testSign)-100:])
		} else {
			t.Logf("完整签名: %s", testSign)
		}

		// 验证签名是否为有效的Base64字符串
		_, err = base64.StdEncoding.DecodeString(testSign)
		if err != nil {
			t.Logf("警告: 签名不是有效的Base64字符串: %v", err)
			t.Logf("提示: 签名可能需要URL解码")
		}

		// 检查公钥格式和完整性
		t.Logf("公钥长度: %d 字符", len(alipayConfig.AlipayPublicKey))
		if strings.Contains(alipayConfig.AlipayPublicKey, "BEGIN") {
			t.Logf("公钥格式: 包含PEM头尾")
		} else {
			t.Logf("公钥格式: 纯Base64编码（将自动添加PEM头尾）")
		}

		// 检查公钥长度（完整的RSA公钥通常有294-450个字符）
		if len(alipayConfig.AlipayPublicKey) < 200 {
			t.Logf("⚠️  警告: 公钥长度可能不完整！")
			t.Logf("   完整的RSA公钥（Base64编码）通常有294-450个字符")
			t.Logf("   当前公钥只有 %d 个字符，可能缺少部分内容", len(alipayConfig.AlipayPublicKey))
			t.Logf("   请确认这是完整的支付宝公钥（不是应用公钥）")
		}

		// 尝试解析公钥，验证格式是否正确
		testClient, err := alipay.NewAlipayClient(alipayConfig)
		if err != nil {
			t.Fatalf("创建支付宝客户端失败（可能是公钥格式错误）: %v", err)
		}
		_ = testClient // 避免未使用变量警告

		// 验证公钥是否能正确解析
		t.Logf("✓ 公钥格式验证通过（可以创建AlipayClient）")

		// 调用 verifySign
		err = logic.verifySign(testData, testSign)
		if err != nil {
			t.Errorf("verifySign 失败: %v", err)
		} else {
			t.Log("✓ verifySign 成功")
		}
	} else {
		// 不使用真实客户端，只测试 buildSignContent 逻辑和 verifySign 的参数验证
		t.Log("=== 测试 buildSignContent 函数和 verifySign 参数验证（mock 模式）===")

		// 测试 buildSignContent
		testBuildSignContent(t, testData)

		// 测试 verifySign 的参数验证逻辑（这些测试不会调用 AlipayClient）
		// 创建一个简单的 ServiceContext，用于测试参数验证
		// 注意：我们不设置 AlipayClient，因为我们会先测试参数验证
		svcCtx := &svc.ServiceContext{
			AlipayClient: nil, // 不设置，因为我们只测试参数验证
			DB:           db,
		}
		ctx := context.Background()
		logic := NewPaymentNotifyLogic(ctx, svcCtx)

		// 测试1: 空签名（会在调用 AlipayClient 之前返回错误）
		t.Run("空签名验证", func(t *testing.T) {
			err := logic.verifySign(testData, "")
			if err == nil {
				t.Error("空签名应该返回错误")
			} else if !strings.Contains(err.Error(), "sign is empty") {
				t.Errorf("期望错误包含 'sign is empty'，但得到: %v", err)
			} else {
				t.Logf("✓ 空签名验证正确返回错误: %v", err)
			}
		})

		// 测试2: 空数据（会在调用 AlipayClient 之前返回错误）
		t.Run("空数据验证", func(t *testing.T) {
			err := logic.verifySign("", "test_sign")
			if err == nil {
				t.Error("空数据应该返回错误")
			} else {
				t.Logf("✓ 空数据验证正确返回错误: %v", err)
			}
		})

		// 测试3: 只有空白字符的数据
		t.Run("空白数据验证", func(t *testing.T) {
			err := logic.verifySign("   ", "test_sign")
			if err == nil {
				t.Error("空白数据应该返回错误")
			} else {
				t.Logf("✓ 空白数据验证正确返回错误: %v", err)
			}
		})

		t.Log("注意: 由于未配置 AlipayClient，无法测试完整的签名验证流程")
		t.Log("提示: 要测试完整的签名验证，请设置 useRealClient = true 并配置 alipayConfig")
	}
}

// testBuildSignContent 测试 buildSignContent 函数的逻辑
func testBuildSignContent(t *testing.T, rawData string) {
	signContent, err := buildSignContent(rawData)
	if err != nil {
		t.Errorf("buildSignContent 失败: %v", err)
		return
	}

	t.Logf("原始数据长度: %d", len(rawData))
	t.Logf("构建的签名字符串长度: %d", len(signContent))
	t.Logf("构建的签名字符串: %s", signContent)

	// 验证 sign 和 sign_type 是否被过滤
	if strings.Contains(signContent, "sign=") {
		t.Error("错误: 签名字符串中不应该包含 sign 参数")
	}
	if strings.Contains(signContent, "sign_type=") {
		t.Error("错误: 签名字符串中不应该包含 sign_type 参数")
	}

	// 验证参数是否按字母顺序排序
	params := strings.Split(signContent, "&")
	if len(params) > 1 {
		for i := 0; i < len(params)-1; i++ {
			key1 := strings.Split(params[i], "=")[0]
			key2 := strings.Split(params[i+1], "=")[0]
			if key1 > key2 {
				t.Errorf("错误: 参数未按字母顺序排序。%s 应该在 %s 之后", key1, key2)
			}
		}
	}

	t.Log("✓ buildSignContent 测试通过")
}

// extractSignFromData 从数据字符串中提取 sign 参数
func extractSignFromData(data string) string {
	pairs := strings.Split(data, "&")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "sign=") {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				// URL 解码
				sign, err := url.QueryUnescape(parts[1])
				if err != nil {
					// 如果解码失败，返回原始值
					return parts[1]
				}
				return sign
			}
		}
	}
	return ""
}

// safeSubstring 安全截取字符串
//func safeSubstring(s string, length int) string {
//	if len(s) <= length {
//		return s
//	}
//	return s[:length] + "..."
//}

// TestVerifySignWithRealData 使用真实数据测试 verifySign（便于快速调试）
// 直接在这里填入你的支付宝通知数据，然后运行这个测试
func TestVerifySignWithRealData(t *testing.T) {
	// ========== 在这里填入你的真实数据 ==========
	// 支付宝通知的完整数据（包含所有参数，包括 sign）
	notifyData := "在这里填入支付宝通知的原始数据"

	// 支付宝公钥（用于验证签名）
	alipayPublicKey := "在这里填入支付宝公钥"

	// 支付宝 AppID（可选，主要用于日志）
	appId := "your_app_id"
	// ========== 数据配置结束 ==========

	// 如果数据未配置，跳过测试
	if notifyData == "在这里填入支付宝通知的原始数据" {
		t.Skip("请先配置 notifyData 和 alipayPublicKey")
		return
	}

	// 从通知数据中提取 sign
	testSign := extractSignFromData(notifyData)
	if testSign == "" {
		// 如果数据中没有 sign，尝试从其他地方获取
		t.Log("警告: 无法从 notifyData 中提取 sign，请手动设置 testSign")
		t.Skip("请手动设置 testSign")
		return
	}

	// 移除 notifyData 中的 sign 和 sign_type（如果存在）
	// verifySign 会自动处理，但为了测试，我们可以先查看构建的内容
	signContent, err := buildSignContent(notifyData)
	if err != nil {
		t.Fatalf("构建签名字符串失败: %v", err)
	}

	t.Logf("=== 使用真实数据测试 ===")
	t.Logf("AppID: %s", appId)
	t.Logf("通知数据长度: %d", len(notifyData))
	t.Logf("签名字符串长度: %d", len(signContent))
	t.Logf("签名字符串: %s", signContent)
	t.Logf("签名长度: %d", len(testSign))
	t.Logf("签名（前50字符）: %s", safeSubstring(testSign, 50))

	// 初始化日志
	logx.DisableStat()

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// 验证配置参数
	if alipayPublicKey == "" || alipayPublicKey == "在这里填入支付宝公钥" {
		t.Fatal("alipayPublicKey 未配置，请设置支付宝公钥")
	}
	if appId == "" || appId == "your_app_id" {
		t.Log("警告: appId 未配置，使用默认值")
		appId = "test_app_id"
	}

	// 创建支付宝客户端配置
	// 注意：虽然验证签名不需要私钥，但 AlipayConfig.Validate() 要求必须有 AppPrivateKey
	// 所以这里使用一个占位符私钥（实际验证签名时不会使用）
	alipayConfig := &alipay.AlipayConfig{
		AppId:           appId,
		AlipayPublicKey: alipayPublicKey,
		SignType:        "RSA2",
		Charset:         "utf-8",
		Format:          "JSON",
		Version:         "1.0",
		IsProd:          false, // 根据实际情况设置
		GatewayUrl:      "https://openapi.alipay.com/gateway.do",
		NotifyUrl:       "",
		ReturnUrl:       "",
		// 验证签名时不需要私钥，但配置验证要求必须有，所以使用占位符
		// 注意：这里使用一个格式正确的占位符私钥（实际不会用于签名）
		AppPrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEjWT2bt\n-----END RSA PRIVATE KEY-----",
		Timeout:       "30m",
	}

	alipayClient, err := alipay.NewAlipayClient(alipayConfig)
	if err != nil {
		t.Fatalf("创建支付宝客户端失败: %v", err)
	}

	// 创建 ServiceContext
	svcCtx := &svc.ServiceContext{
		AlipayClient: alipayClient,
		DB:           db,
	}

	// 创建 context
	ctx := context.Background()

	// 创建 PaymentNotifyLogic
	logic := NewPaymentNotifyLogic(ctx, svcCtx)

	// 调用 verifySign
	err = logic.verifySign(notifyData, testSign)
	if err != nil {
		t.Errorf("签名验证失败: %v", err)
		t.Logf("提示: 请检查:")
		t.Logf("1. alipayPublicKey 是否正确")
		t.Logf("2. notifyData 是否完整")
		t.Logf("3. testSign 是否正确")
	} else {
		t.Log("✓ 签名验证成功！")
	}
}

// TestBuildSignContent 单独测试 buildSignContent 函数
func TestBuildSignContent(t *testing.T) {
	testCases := []struct {
		name    string
		rawData string
		wantErr bool
	}{
		{
			name:    "正常数据",
			rawData: "gmt_create=2024-01-01+12%3A00%3A00&charset=utf-8&out_trade_no=ORDER123&total_amount=100.00",
			wantErr: false,
		},
		{
			name:    "包含sign和sign_type",
			rawData: "gmt_create=2024-01-01+12%3A00%3A00&charset=utf-8&sign=test_sign&sign_type=RSA2&out_trade_no=ORDER123",
			wantErr: false,
		},
		{
			name:    "空数据",
			rawData: "",
			wantErr: true,
		},
		{
			name:    "URL编码数据",
			rawData: "subject=%E6%B5%8B%E8%AF%95%E5%95%86%E5%93%81&charset=utf-8&out_trade_no=ORDER123",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := buildSignContent(tc.rawData)
			if tc.wantErr {
				if err == nil {
					t.Errorf("buildSignContent 应该返回错误，但没有返回")
				}
				return
			}

			if err != nil {
				t.Errorf("buildSignContent 不应该返回错误: %v", err)
				return
			}

			if result == "" {
				t.Error("buildSignContent 返回的结果不应该为空")
				return
			}

			// 验证 sign 和 sign_type 被过滤
			if strings.Contains(result, "sign=") || strings.Contains(result, "sign_type=") {
				t.Error("结果中不应该包含 sign 或 sign_type 参数")
			}

			t.Logf("输入: %s", tc.rawData)
			t.Logf("输出: %s", result)
		})
	}
}
