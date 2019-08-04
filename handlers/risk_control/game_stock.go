package risk_control

type GameStock struct {
	BaseControl
	RoundCount int64
}

func (t *GameStock)Check()int64{
	return 0
}

func (t *GameStock)CheckContinuousTriggering()int64{
	return 0
}
func (t *GameStock)setSelfConf(appId int64)int64{
	return 0
}