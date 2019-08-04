package risk_control

type BaseControl struct {
	IsOn int64		// 开关 1：表示开  0：表示关
	IsInitConf bool	// 是否初始化配置
	KillRate int64			// 表示百分比，   例如 40 表示 40%
}
