package other

import (
	"encoding/json"
	"fmt"
	"go_utils/logger"
	"io"
	"net/http"
	"time"
)

// ipipNetResponse 定义 myip.ipip.net 返回的数据结构.
// {"ret":"ok","data":{"ip":"49.82.159.197","location":["中国","江苏","淮安","","电信"]}}
type ipipNetResponse struct {
	Ret  string `json:"ret"`
	Data struct {
		IP       string   `json:"ip"`
		Location []string `json:"location"`
	} `json:"data"`
}

// sharedHTTPClient 避免每次请求都创建新的 http.Client, 并设置合理超时.
var sharedHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// GetPublicIPByIPIPNet 从 myip.ipip.net 获取当前公网 IP.
// 如果请求或解析失败, 返回具体的错误信息以方便排查.
func GetPublicIPByIPIPNet() (string, error) {
	resp, err := sharedHTTPClient.Get("https://myip.ipip.net/json")
	if err != nil {
		logger.SugaredLogger.Errorw("访问ipip.net时发生错误", "url", "err", err)
		return "", fmt.Errorf("访问myip.ipip.net时发生错误: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.SugaredLogger.Errorw("读取 ipip.net 返回的内容时发生错误", "err", err)
		return "", fmt.Errorf("读取myip.ipip.net返回的内容时发生错误: %w", err)
	}

	var ipResp ipipNetResponse
	if err := json.Unmarshal(body, &ipResp); err != nil {
		logger.SugaredLogger.Errorw("反序列化 ipip.net 返回的内容时发生错误", "body", string(body), "err", err)
		return "", fmt.Errorf("反序列化myip.ipip.net返回的内容时发生错误: %w", err)
	}

	if ipResp.Data.IP == "" {
		return "", fmt.Errorf("myip.ipip.net返回的内容未找到IP信息")
	}
	return ipResp.Data.IP, nil
}
