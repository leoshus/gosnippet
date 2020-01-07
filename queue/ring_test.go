/**
 * Created by shangyindong
 * Date 2020/1/2 2:13 下午
 **/
package queue

import (
	"idgenerator/util"
	"testing"
)

func BenchmarkRBGet(b *testing.B) {
	rb := NewRingBuffer(uint64(b.N))

	node, _ := util.NewNode(1)
	for i := 0; i < b.N; i++ {
		rb.Put(node.Generate().Int64())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rb.Get()
	}

}

func BenchmarkSnowflake(b *testing.B) {
	node, _ := util.NewNode(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Generate().Int64()
	}

}