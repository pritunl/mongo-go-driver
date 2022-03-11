package driver

import (
	"testing"

	"github.com/pritunl/mongo-go-driver/internal/testutil/assert"
)

func TestBatchCursor(t *testing.T) {
	t.Run("setBatchSize", func(t *testing.T) {
		var size int32
		bc := &BatchCursor{
			batchSize: size,
		}
		assert.Equal(t, size, bc.batchSize, "expected batchSize %v, got %v", size, bc.batchSize)

		size = int32(4)
		bc.SetBatchSize(size)
		assert.Equal(t, size, bc.batchSize, "expected batchSize %v, got %v", size, bc.batchSize)
	})
}
