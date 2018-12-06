package leveldb

type Env interface {
}


type defaultEnv struct {
}

func DefaultEnv() Env {
	return defaultEnv{}
}