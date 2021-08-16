package models

type EcsDetail struct {
	Ecs
	// 扩展信息
}

type RdsDetail struct {
	Rds
	// 扩展信息
}

type KvDetail struct {
	Kv
	// 扩展信息
}

type SlbDetail struct {
	Slb
	// 扩展信息
}

type OtherDetail struct {
	OtherRes
	// 扩展信息
}