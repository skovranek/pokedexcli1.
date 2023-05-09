package pokecache

import (
    "testing"
    "time"
    "fmt"
	"internal/pokecache"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://website.com",
			val: []byte("characters"),
		},
		{
			key: "https://website.com/path",
			val: []byte("stringcharacters"),
		},
		{
			key: "https://website.com/pathname",
			val: []byte("yetstringcharacters"),
		},
		{
			key: "https://website.com/pathnameis",
			val: []byte("alsostringcharacters"),
		},
		{
			key: "https://website.com/pathnameisnot",
			val: []byte("superplusstringcharacters"),
		},
		{
			key: "https://website.com/pathnameisnotthis",
			val: []byte("getstringcharacters"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to get key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to get value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://website.com", []byte("characters"))
	cache.Add("https://website.com/path", []byte("stringcharacters"))
	cache.Add("https://website.com/pathname", []byte("yetstringcharacters"))

	_, ok := cache.Get("https://website.com")
	if !ok {
		t.Errorf("expected to get key")
		return
	}
	_, ok = cache.Get("https://website.com/path")
	if !ok {
		t.Errorf("expected to get key")
		return
	}
	_, ok = cache.Get("https://website.com/pathname")
	if !ok {
		t.Errorf("expected to get key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://website.com")
	if ok {
		t.Errorf("expected to not get key")
		return
	}
	_, ok = cache.Get("https://website.com/path")
	if ok {
		t.Errorf("expected to not get key")
		return
	}
	_, ok = cache.Get("https://website.com/pathname")
	if ok {
		t.Errorf("expected to not get key")
		return
	}
}
