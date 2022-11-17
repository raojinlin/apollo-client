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
