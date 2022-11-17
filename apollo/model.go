package apollo

type Response struct {
	AppId          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

type NotificationResponse struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int    `json:"notificationId"`
	Messages       struct {
		Details map[string]int
	} `json:"messages"`
}

type NotificationRequestPayload struct {
	NamespaceName  string `json:"namespaceName"`
	Cluster        string `json:"cluster"`
	NotificationId int    `json:"notificationId"`
}

type Option struct {
	Server     string   `json:"server"`
	Cluster    string   `json:"cluster"`
	AppId      string   `json:"appId"`
	Namespaces []string `json:"namespaces"`
	CacheDir   string   `json:"cacheDir"`
}
