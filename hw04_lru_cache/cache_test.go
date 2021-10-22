package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(1)

		c.Set("test", 123)
		val, _ := c.Get("test")
		require.Equal(t, val, 123)

		c.Clear()
		_, ok := c.Get("test")
		require.False(t, ok)
	})

	t.Run("values of any type", func(t *testing.T) {
		c := NewCache(5)

		testsCases := []struct {
			t string
			v interface{}
		}{
			{t: "int", v: 10},
			{t: "float", v: 10.123},
			{t: "string", v: "test"},
			{t: "bool", v: true},
			{t: "struct", v: struct{ field int }{field: 1}},
		}

		for _, tc := range testsCases {
			tc := tc

			t.Run(tc.t, func(t *testing.T) {
				key := Key(tc.t)
				c.Set(key, tc.v)
				value, _ := c.Get(key)
				require.Equal(t, tc.v, value)
			})
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
