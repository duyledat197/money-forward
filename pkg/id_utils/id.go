package id_utils

type IDGenerator interface {
	String() string
	Int64() int64
}
