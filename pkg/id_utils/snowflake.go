package id_utils

import "github.com/bwmarrin/snowflake"

type snowFlake struct {
	generator *snowflake.Node
}

func NewSnowFlake(nodeID int64) IDGenerator {
	generator, _ := snowflake.NewNode(nodeID)

	return &snowFlake{
		generator: generator,
	}
}

func (s *snowFlake) String() string {
	return s.generator.Generate().String()
}

func (s *snowFlake) Int64() int64 {
	return s.generator.Generate().Int64()
}
