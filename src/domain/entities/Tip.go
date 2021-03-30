package entities

import "gorm.io/gorm"

type Tip struct {
	gorm.Model
	Name string
}

type Tips []*Tip
