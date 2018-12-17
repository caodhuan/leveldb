package leveldb


type FilterPolicy interface {
	Name() string
	CreateFilter(keys []string, dst *string)
	KeyMayMatch(key string, filter string) bool
}

type internalFilterPolicy struct {
	userPolicy FilterPolicy
}

func (this *internalFilterPolicy) Name() string {
	return this.userPolicy.Name()
}

func (this *internalFilterPolicy) CreateFilter(keys []string, dst *string) {
	realKeys := make([]string, len(keys))
	for index, value := range keys {
		realKeys[index] = extractUserKey(value)
		 // TODO(sanjay): Suppress dups?
	}

	this.userPolicy.CreateFilter(realKeys, dst)
}

func (this *internalFilterPolicy) KeyMayMatch(key string, filter string) bool {
	return this.userPolicy.KeyMayMatch(key, filter)
}