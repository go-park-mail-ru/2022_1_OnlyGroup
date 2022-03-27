package redis

import (
	"github.com/alicebob/miniredis"
	"testing"
)

func TestSessionsRedisRepo(t *testing.T) {
	_, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis start failed %v", err)
		return
	}
}
