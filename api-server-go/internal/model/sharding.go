package model

import (
	"hash/fnv"
	"strconv"
)

const UnionidMappingShardCount = 10

func GetUnionidMappingTableName(unionid string) string {
	h := fnv.New32a()
	h.Write([]byte(unionid))
	shard := h.Sum32() % UnionidMappingShardCount
	return "mc_work_unionid_external_userid_mapping_" + strconv.Itoa(int(shard))
}

func GetUnionidMappingTableNames() []string {
	names := make([]string, UnionidMappingShardCount)
	for i := 0; i < UnionidMappingShardCount; i++ {
		names[i] = "mc_work_unionid_external_userid_mapping_" + strconv.Itoa(i)
	}
	return names
}
