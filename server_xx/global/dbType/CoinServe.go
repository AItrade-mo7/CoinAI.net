package dbType

// 一个虚拟主机的数据库模板 AIFundServer
type CoinServeTable struct {
	CoinServeID string `bson:"CoinServeID"` // 这个ID是唯一的，其值为 IP+端口 的字符串
	OkxKeyID    string `bson:"OkxKeyID"`    // ApiKey 的 ID 一一对应
	UserID      string `bson:"UserID"`      // UserID 的 ID 一对多
	Host        string `bson:"Host"`        // 服务的 Host  , 与 OkxKeys 的相同
	Port        string `bson:"Port"`        // 需要请求的服务的端口  Host + Port 必须唯一
	Note        string `bson:"Note"`        // 备注
	CreateTime  int64  `bson:"CreateTime"`  // 创建时间
}
