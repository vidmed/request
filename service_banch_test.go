package request

import (
	"testing"
)

var Result Request

func BenchmarkGetRequest(b *testing.B) {
	b.ReportAllocs()
	srvc := NewRequestService()
	for i := 0; i < b.N; i++ {
		Result = srvc.GetRequest()
	}
}
