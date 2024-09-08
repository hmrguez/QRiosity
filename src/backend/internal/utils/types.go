package utils

import "encoding/json"

type AppSyncEvent struct {
	TypeName  string            `json:"parentTypeName"`
	FieldName string            `json:"fieldName"`
	Arguments json.RawMessage   `json:"arguments"`
	Headers   map[string]string `json:"headers"`
}
