package ip

import (
	"encoding/json"
	"fmt"
	"go_utils/logger"
	"io"
	"net/http"
	"time"
)

const (
	IPIP_NET_URL = "https://myip.ipip.net/json"
)

// sharedHTTPClient 避免每次请求都创建新的 http.Client, 并设置合理超时.
var sharedHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// GetPublicIP 从 myip.ipip.net 获取当前公网 IP.
// 如果请求或解析失败, 返回具体的错误信息以方便排查.
func GetPublicIP() (string, error) {
	return queryIpIpNet()
}

// queryIpIpNet 使用 ipip.net 来查询公网地址
func queryIpIpNet() (string, error) {
	// ipIpNetResponse 定义 myip.ipip.net 返回的数据结构.
	// {"ret":"ok","data":{"ip":"49.82.159.197","location":["中国","江苏","淮安","","电信"]}}
	type ipIpNetResponse struct {
		Ret  string `json:"ret"`
		Data struct {
			IP       string   `json:"ip"`
			Location []string `json:"location"`
		} `json:"data"`
	}

	resp, err := sharedHTTPClient.Get(IPIP_NET_URL)
	if err != nil {
		logger.SugaredLogger.Errorw("访问时发生错误", "url", IPIP_NET_URL, "err", err)
		return "", fmt.Errorf("访问%s时发生错误: %s", IPIP_NET_URL, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.SugaredLogger.Errorw("读取返回的内容时发生错误", "url", IPIP_NET_URL, "err", err)
		return "", fmt.Errorf("读取%s返回的内容时发生错误: %s", IPIP_NET_URL, err.Error())
	}

	var ipResp ipIpNetResponse
	if err = json.Unmarshal(body, &ipResp); err != nil {
		logger.SugaredLogger.Errorw("反序列化返回的内容时发生错误",
			"url", IPIP_NET_URL, "body", string(body), "err", err)
		return "", fmt.Errorf("反序列化myip.ipip.net返回的内容时发生错误: %w", err)
	}

	if ipResp.Data.IP == "" {
		return "", fmt.Errorf("%s返回的内容未找到IP信息", IPIP_NET_URL)
	} else {
		return ipResp.Data.IP, nil
	}
}
