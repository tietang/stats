package stats

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counters map[string]Counter

const defaultCounterName = "defaultCounter"

func init() {
	counters = make(map[string]Counter)
}

type Metric struct {
	Count  int64 //次数
	AvgDur int64 //平均耗时
	sum    int64 //总累计时间
	MaxDur int64 //最大
	MinDur int64 //最小
}

type Counter struct {
	lock sync.RWMutex
	m    map[string]Metric
	name string
}

func NewCounter(names ...string) *Counter {
	c := new(Counter)
	c.lock = sync.RWMutex{}
	c.m = make(map[string]Metric)
	if len(names) != 0 {
		c.name = names[0]
		counters[names[0]] = *c
	} else {
		counters[defaultCounterName] = *c
	}
	return c
}
func (c *Counter) Increment(key string) *Metric {
	return c.Incr(key, 1)
}
func (c *Counter) Incr(key string, duration int64) *Metric {
	c.lock.Lock()
	defer c.lock.Unlock()
	metric := c.m[key]
	metric.Count++
	if metric.MaxDur < duration {
		metric.MaxDur = duration
	}
	if metric.MinDur == 0 || metric.MinDur > duration {
		metric.MinDur = duration
	}
	metric.sum += duration
	metric.AvgDur = metric.sum / metric.Count
	c.m[key] = metric
	return &metric
}

func (c *Counter) Get(key string) *Metric {
	c.lock.RLock()
	defer c.lock.RUnlock()
	v := c.m[key]
	return &v
}

func (c *Counter) GetAll() map[string]Metric {
	return c.m
}

func (c *Counter) HttpIncr(req *http.Request, duration int64) {
	path := req.RequestURI
	c.Incr(path, duration)
}

func (c *Counter) httpHandle(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(c.m)
	if err == nil {
		w.Write(b)
	} else {
		io.WriteString(w, "{}")
	}
}

func (c *Counter) HttpServer(port int) {
	http.HandleFunc("/metrics", c.httpHandle)
	log.Println("http://127.0.0.1:" + strconv.Itoa(port) + "/metrics")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
