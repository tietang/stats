package stats

import (
	"github.com/tietang/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCounterIncr(t *testing.T) {
	c := NewCounter()
	key := "test1"
	channel := make(chan *Metric)
	size := 1000
	for i := 0; i < size; i++ {
		go func() {
			v := c.Incr(key, rand.Int63()%10000)
			if v.Count >= int64(size) {
				channel <- v
			}

		}()

	}
	s := <-channel
	t.Log(s)

	v := c.Get(key)
	assert.Equals(t, size, v.Count)
	t.Log(v)
	//	time.Sleep(1000 * time.Millisecond)
}

func TestCounterIncrAll(t *testing.T) {
	c := NewCounter()
	key := "/v1/test1"
	size := 10000
	for i := 0; i < size; i++ {
		c.Incr(key+strconv.Itoa(int(rand.Int63()%10)), rand.Int63()%10000)
	}

	vs := c.GetAll()
	for k, v := range vs {
		t.Log(k, v)
	}

	//	time.Sleep(1000 * time.Millisecond)
}

func TestCounterIncr2(t *testing.T) {
	c := NewCounter()
	key := "test1"
	size := 10000
	for i := 0; i < size; i++ {
		c.Incr(key, rand.Int63()%10000)
	}

	v := c.Get(key)
	assert.Equals(t, size, v.Count)
	t.Log(v)
	//	time.Sleep(1000 * time.Millisecond)
}
func TestCounterIncrHttpServer(t *testing.T) {
	c := NewCounter()

	key := "test1"
	size := 10000
	for i := 0; i < size; i++ {
		c.Incr(key, rand.Int63()%10000)
	}

	v := c.Get(key)
	assert.Equals(t, size, v.Count)
	t.Log(v)
	go func() {
		time.Sleep(5000 * time.Millisecond)
		os.Exit(0)
	}()
	c.HttpServer(8888)

}
