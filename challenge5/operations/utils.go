package operations

import (
	"fmt"
	"sync"
)

type log struct {
	offset   int
	commited bool
	value    int
}

type nodeState struct {
	mu   sync.RWMutex
	logs map[string][]log
}

func NewNodeState() *nodeState {
	return &nodeState{
		logs: make(map[string][]log),
	}
}

var state = NewNodeState()

func (n *nodeState) addLog(key string, value int) log {
	n.mu.Lock()
	defer n.mu.Unlock()

	logList, exists := n.logs[key]
	if !exists {
		logList = make([]log, 0)
	}

	newLog := log{
		offset:   len(logList),
		commited: false,
		value:    value,
	}

	logList = append(logList, newLog)
	n.logs[key] = logList

	return newLog
}

func (n *nodeState) getLog(key string, offset int) (log, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	logList, exists := n.logs[key]
	if !exists {
		return log{}, fmt.Errorf("the key %s does not exist", key)
	}

	if offset >= len(logList) {
		return log{}, fmt.Errorf("the offset %d is out of range", offset)
	}

	return logList[offset], nil
}

func (n *nodeState) getLogs(key string, offset int) ([]log, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	logs, exists := n.logs[key]
	if !exists {
		return nil, fmt.Errorf("there is no log for the given key: %s", key)
	}

	if offset >= len(logs) {
		return nil, fmt.Errorf("there is no value in the log for the given offset: %d", offset)
	}

	return logs[offset:], nil
}

func (n *nodeState) commitOffset(key string, upToOffset int) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	logs, exists := n.logs[key]
	if !exists {
		return fmt.Errorf("there is no log for the given key: %s", key)
	}

	if upToOffset >= len(logs) {
		return fmt.Errorf("there is no value in the log for the given offset: %d", upToOffset)
	}

	for i := 0; i <= upToOffset; i++ {
		logs[i].commited = true
	}

	return nil
}

func (n *nodeState) getHighestCommitedOffset(key string) (int, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	logs, exists := n.logs[key]
	if !exists {
		return 0, false
	}

	highestOffset := -1
	left, right := 0, len(logs)-1

	for left <= right {
		mid := (left + right) / 2

		if logs[mid].commited {
			highestOffset = logs[mid].offset
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if highestOffset == -1 {
		return 0, false
	}

	return highestOffset, true
}
