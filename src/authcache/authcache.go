package authcache

import "time"

type Entry struct {
	stamp time.Time
}

type Cache struct {
	cache map[string]Entry
}

func New() *Cache {
	c := &Cache{make(map[string]Entry)}
	go c.cleanup()
	return c
}

func (c *Cache) Add(key string) {
	c.cache[key] = Entry{time.Now()}
}

func (c *Cache) Check(key string) bool {
	_, ok := c.cache[key]
	return ok
}

func (c *Cache) cleanup() {
	for {
		toDelete := []string{}
		for key, entry := range c.cache {
			if time.Since(entry.stamp) > 10*time.Minute {
				toDelete = append(toDelete, key)
			}
		}
		for _, key := range toDelete {
			delete(c.cache, key)
		}
		time.Sleep(1 * time.Minute)
	}
}
