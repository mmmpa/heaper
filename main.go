package heaper

import (
	"time"
	"runtime/pprof"
	"regexp"
	"fmt"
	"strings"
	"encoding/json"
)

type Each struct {
	Alloc        int
	TotalAlloc   int
	Sys          int
	Lookups      int
	Mallocs      int
	Frees        int
	HeapAlloc    int
	HeapSys      int
	HeapIdle     int
	HeapInuse    int
	HeapReleased int
	HeapObjects  int
	StackTotal   int
	MSpanTotal   int
	MCacheTotal  int
	StackPart    int
	MSpanPart    int
	MCachePart   int
	BuckHashSys  int
	NextGC       int
	PauseNs      []int
	NumGC        int
	DebugGC      bool
	Time         int
}

type Process struct {
	Sec   int
	Size  int
	Head  int
	Pos   int
	Stack []Each
}

var (
	opened bool
	q chan struct{}
	process Process
)

func Run(sec, size int) {
	Stop()

	q = open()
	go StartNewProcess(sec, size)
	<-q
}

func open() chan struct{} {
	Stop()
	opened = true
	return make(chan struct{})
}

func Stop() {
	if opened {
		opened = false
		q <- struct{}{}
	}
}

func Read() []Each {
	return process.Read()
}

func StartNewProcess(sec, size int) (*Process) {
	process = Process{
		Sec: sec,
		Size: size,
	}
	process.Stack = make([]Each, size, size)
	process.Run()
	return &process
}

func (p *Process) Run() {
	for {
		time.Sleep(time.Duration(p.Sec) * time.Second)
		p.Stock()
	}
}

func (p *Process) Stock() {
	now := &Each{}
	pprof.WriteHeapProfile(now)
	p.Stack[p.Pos % p.Size] = *now

	p.Pos += 1
}

func (p *Process) Read() []Each {
	if p.Pos < p.Size {
		return p.Stack
	}

	stack := make([]Each, p.Size, p.Size)

	for i := 0; i < p.Size; i++ {
		stack[i] = p.Stack[(i + p.Pos) % p.Size]
	}

	return stack
}

var picker = regexp.MustCompile(`([a-zA-z]+) = ([^\n]+)`)

func (e *Each) Write(p []byte) (int, error) {
	j := "{"
	for _, m := range picker.FindAllStringSubmatch(string(p), -1) {
		if m[1] == "PauseNs" {
			j += fmt.Sprintf(`"%vPart":%v`, m[1], "[]")
		} else if strings.Index(m[2], " / ") != -1 {
			splat := strings.Split(m[2], " / ")
			j += fmt.Sprintf(`"%vPart":%v,`, m[1], splat[0])
			j += fmt.Sprintf(`"%vTotal":%v`, m[1], splat[1])
		} else {
			j += fmt.Sprintf(`"%v":%v`, m[1], m[2])
		}
		j += ","
	}
	j += fmt.Sprintf(`"Time":%v`, time.Now().Unix())
	j += "}"

	json.Unmarshal([]byte(j), e)

	return len(p), nil
}
