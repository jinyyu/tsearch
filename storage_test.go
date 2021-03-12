package tsearch

import "testing"

func Test_RedisStorage(t *testing.T) {
	s, err := NewRedisStorage("127.0.0.1:6379")
	if err != nil {
		t.Errorf("NewRedisStorage error")
	}

	kvs := []*KeyValue{
		{
			"abc",
			"def",
		},
		{
			"123",
			"456",
		},
	}

	err = s.MultiSet(kvs...)
	if err != nil {
		t.Errorf("MultiSet error")
	}

	values, err := s.MultiGet("abc", "kkk", "123")
	if err != nil {
		t.Errorf("MultiGet error %s", err)
	}
	if len(values) != 3 || values[0] != "def" || values[1] != "" || values[2] != "456" {
		t.Errorf("MultiGet error %s", err)
	}
}
