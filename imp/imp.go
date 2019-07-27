package imp

import "awesomeProject/stract"

func (t *stract.Class)GetName()string{
	return t.name
}
func NewClass() stract.Class{
	return &stract.Class{}
}
