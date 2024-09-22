package model

import (
	"strconv"

	"jcourse_go/model/po"
)

type SettingType = string

const (
	SettingTypeString SettingType = "string"
	SettingTypeInt    SettingType = "int"
	SettingTypeBool   SettingType = "bool"
)

type Setting interface {
	GetKey() string
	GetValue() any
	ToPO() po.SettingPO
	FromPO(po.SettingPO) error
}

type StringSetting struct {
	Key   string
	Value string
}

func (s *StringSetting) GetKey() string {
	return s.Key
}

func (s *StringSetting) GetValue() any {
	return s.Value
}

func (s *StringSetting) ToPO() po.SettingPO {
	return po.SettingPO{
		Key:   s.Key,
		Value: s.Value,
		Type:  SettingTypeString,
	}
}

func (s *StringSetting) FromPO(po po.SettingPO) error {
	s.Key = po.Key
	s.Value = po.Value
	return nil
}

type IntSetting struct {
	Key   string
	Value int64
}

func (s *IntSetting) GetKey() string {
	return s.Key
}

func (s *IntSetting) GetValue() any {
	return s.Value
}

func (s *IntSetting) ToPO() po.SettingPO {
	return po.SettingPO{
		Key:   s.Key,
		Value: strconv.FormatInt(s.Value, 10),
		Type:  SettingTypeInt,
	}
}

func (s *IntSetting) FromPO(po po.SettingPO) (err error) {
	s.Key = po.Key
	s.Value, err = strconv.ParseInt(po.Value, 10, 64)
	return
}

type BoolSetting struct {
	Key   string
	Value bool
}

func (s *BoolSetting) GetKey() string {
	return s.Key
}

func (s *BoolSetting) GetValue() any {
	return s.Value
}

func (s *BoolSetting) ToPO() po.SettingPO {
	return po.SettingPO{
		Key:   s.Key,
		Value: strconv.FormatBool(s.Value),
		Type:  SettingTypeBool,
	}
}

func (s *BoolSetting) FromPO(po po.SettingPO) (err error) {
	s.Key = po.Key
	s.Value, err = strconv.ParseBool(po.Value)
	return
}
