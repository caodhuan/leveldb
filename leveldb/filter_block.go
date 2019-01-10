package leveldb


const (
	// Generate new filter every 2KB of data
	kFilterBaseLg = 11
	kFilterBase = 1 << kFilterBaseLg
)

type FilterBlockReader struct {
	policy FilterPolicy
	
}

type FilterBlockBuilder struct {
	policy FilterPolicy
	keys string 	// Flattened key contents
	start []uint 	// Starting index in keys of each key
	result string	// Filter data computed so far
	tmpKeys []string	// policy->CreateFilter() argument
	filterOffsets []int
}


func newFliterBlockBuilder(policy FilterPolicy) *FilterBlockBuilder {
	return &FilterBlockBuilder{
		policy: policy,
	}
}

func (this *FilterBlockBuilder) StartBlock(blockOffset uint64) {
	filterIndex := blockOffset / kFilterBase
	for filterIndex > uint64(len(this.filterOffsets) ) {
		this.generateFilter()
	}
}


func (this *FilterBlockBuilder) generateFilter() {
	numKeys := len(this.start)
	if (numKeys == 0) {
		// Fast path if there are no keys for this filter
		this.filterOffsets = append(this.filterOffsets, len(this.result) )
	}
}