package global

import (
	"fmt"
	"time"
)

type (
	// 平台表映射结构
	Platform struct {
		ID       int64     `json:"id"`
		Name     string    `json:"name"`
		Link     string    `json:"link"`
		Owner    string    `json:"owner"`
		CreateAt time.Time `json:"create_at"`
		UpdateAt time.Time `json:"update_at"`
	}

	// 平台列表
	PlatformList []Platform
)

const (
	MySQLPlatformTableName = "static_platform"
)

var ()

func (p *Platform) MySQLTableName() string {
	return MySQLPlatformTableName
}

func (pl PlatformList) QueryOrder() string {
	return fmt.Sprintf("select * from %s", MySQLPlatformTableName)
}
