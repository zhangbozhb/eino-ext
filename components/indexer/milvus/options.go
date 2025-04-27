package milvus

import "github.com/cloudwego/eino/components/indexer"

type ImplOptions struct {
	Partition string
}

func WithPartition(partition string) indexer.Option {
	return indexer.WrapImplSpecificOptFn(func(o *ImplOptions) {
		o.Partition = partition
	})
}
