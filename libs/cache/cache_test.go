package cache

import (
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	c := NewMemoryCache()

	// Test Set & Get
	c.Set("foo", "bar", 0)
	value, err := c.Get("foo")
	if err != nil || value.(string) != "bar" {
		t.Errorf("Expected 'bar', got %v", value)
	}

	// Test TTL Expiration
	c.Set("temp", "expire", 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	_, err = c.Get("temp")
	if err == nil {
		t.Errorf("Expected 'temp' to expire, but found in cache")
	}

	// Test Delete
	c.Set("delete_me", "gone", 0)
	c.Delete("delete_me")
	_, err = c.Get("delete_me")
	if err == nil {
		t.Errorf("Expected 'delete_me' to be deleted")
	}

	// Test Flush
	c.Set("key1", "val1", 0)
	c.Set("key2", "val2", 0)
	c.Flush()
	_, err1 := c.Get("key1")
	_, err2 := c.Get("key2")
	if err1 == nil || err2 == nil {
		t.Errorf("Expected all keys to be flushed")
	}
}
