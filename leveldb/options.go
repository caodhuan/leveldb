package leveldb

type CompressionType int32

const (
	NoCompression = iota
	SnappyCompression
)

type Options struct {
	Compartor 			Compartor
	CreateIfMissing 	bool
	ErrorIfExists		bool
	ParanoidChecks		bool
	Env					Env
	InfoLog				Logger
	WriteBufferSize		uint32
	MaxOpenFiles		int32
	BlockCache			Cache
	BlockSize			uint32
	BlockRestartInterval int32
	Compression			CompressionType
	FilterPolicy		FilterPolicy
}

func NewOptions() *Options {
	return &Options {
		Compartor: BytewiseComparator(),
		CreateIfMissing: false,
		ErrorIfExists: false,
		ParanoidChecks: false,
		Env: DefaultEnv(),
		InfoLog: nil,
		WriteBufferSize: 4 << 20,
		MaxOpenFiles: 1000,
		BlockCache: nil,
		BlockSize: 4096,
		BlockRestartInterval: 16,
		Compression: NoCompression,
		FilterPolicy: nil,
	}
}