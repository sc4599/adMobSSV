package model

import "awesomeProject/handlers/risk_control"

var Rc risk_control.RiskControl

func init(){
	Rc = risk_control.NewRiskControl(101)
}
