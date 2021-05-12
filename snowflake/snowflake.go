package snowflake

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var node = NewNode()

func GenerateID() (ID, error) {
	return node.NextId()
}

func NewID() ID {
	return Must(node.NextId())
}

func Must(id ID, err error) ID {
	if err != nil {
		panic(err)
	}
	return id
}

type ID int64

type Node struct {
	workerId           int64
	workerIdBits       int64
	workerIdShift      int64
	maxWorkerId        int64
	dataCenterId       int64
	dataCenterIdBits   int64
	dataCenterIdShift  int64
	maxDataCenterId    int64
	sequence           int64
	sequenceBits       int64
	sequenceMask       int64
	startTimestamp     int64
	lastTimestamp      int64
	timestampLeftShift int64

	mu sync.RWMutex
}

type Option func(n *Node)

func WorkerIdOption(workerId int64) Option {
	return func(n *Node) {
		n.workerId = workerId
	}
}

func DataCenterIdOption(dataCenterId int64) Option {
	return func(n *Node) {
		n.dataCenterId = dataCenterId
	}
}

func NewNode(opts ...Option) *Node {
	n := &Node{
		workerId:         0,
		dataCenterId:     0,
		sequence:         0,
		sequenceMask:     4095,
		sequenceBits:     12,
		workerIdBits:     5,
		dataCenterIdBits: 5,
		maxWorkerId:      31,
		maxDataCenterId:  31,
		startTimestamp:   1288834974657,
	}

	n.workerIdShift = n.sequenceBits
	n.dataCenterIdShift = n.sequenceBits + n.workerIdBits
	n.timestampLeftShift = n.sequenceBits + n.workerIdBits + n.dataCenterIdBits

	for _, opt := range opts {
		opt(n)
	}

	_ = n.init()

	return n
}

func (n *Node) NextId() (ID, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	timestamp := time.Now().UTC().UnixNano() / 1e6
	if timestamp < n.lastTimestamp {
		return ID(0), fmt.Errorf("clock moved backwards.  refusing to generate id for %d milliseconds", n.lastTimestamp-timestamp)
	}
	n.sequence = (n.sequence + 1) & n.sequenceMask
	if n.sequence == 0 {
		timestamp = n.tillNextMillis(n.lastTimestamp)
	}
	n.lastTimestamp = timestamp

	l := ((n.lastTimestamp - n.startTimestamp) << n.timestampLeftShift) | (n.dataCenterId << n.dataCenterIdShift) | (n.workerId << n.workerIdShift) | n.sequence
	id, err := ParseInt(l)
	if err != nil {
		return ID(0), err
	}
	return id, nil
}

func (n *Node) tillNextMillis(lastTimestamp int64) int64 {
	ts := time.Now().UTC().UnixNano() / 1e6
	for {
		if ts > lastTimestamp {
			break
		}
		ts = time.Now().UTC().UnixNano() / 1e6
	}
	return ts
}

func (n *Node) init() error {
	hn, err := os.Hostname()
	if err != nil {
		return err
	}
	hnSplited := strings.Split(hn, "-")
	if len(hnSplited) > 0 {
		hnId, err := strconv.ParseInt(hnSplited[len(hnSplited)-1], 10, 64)
		if err != nil {
			return err
		}
		if hnId > n.maxWorkerId {
			return fmt.Errorf("worker Id can't be greater than %d or less than 0", n.maxWorkerId)
		}
		n.workerId = hnId
		if hnId > n.maxDataCenterId {
			return fmt.Errorf("data center Id can't be greater than %d or less than 0", n.maxDataCenterId)
		}
		n.dataCenterId = hnId
	} else {
		n.workerId, n.dataCenterId = 0, 0
	}
	return nil
}

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) String() string {
	return strconv.FormatInt(id.Int64(), 10)
}

func ParseInt(i int64) (ID, error) {
	return ID(i), nil
}

func ParseString(s string) (ID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ID(0), err
	}
	return ID(i), nil
}
