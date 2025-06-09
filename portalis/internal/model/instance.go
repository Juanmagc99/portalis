package model

import "time"

type Instance struct {
	ServiceName string            `json:"serviceName`
	InstanceID  string            `json:"instanceID"`
	Host        string            `json:"host"`
	Port        int               `json:"host"`
	Metadata    map[string]string `json:"metadata"`
	LastSeen    time.Time         `json:"lastSeen"`
}
