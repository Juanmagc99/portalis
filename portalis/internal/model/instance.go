package model

import "time"

type Instance struct {
	ServiceName string            `json:"serviceName validate:"required"`
	InstanceID  string            `json:"instanceID" validate:"required"`
	Host        string            `json:"host" validate:"required"`
	Port        int               `json:"host" validate:"required"`
	Metadata    map[string]string `json:"metadata" validate:"required"`
	LastSeen    time.Time         `json:"lastSeen" validate:"required"`
}
