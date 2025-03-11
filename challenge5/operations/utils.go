package operations

import (
	"context"
	"fmt"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type log struct {
	Offset   int  `json:"offset"`
	Commited bool `json:"commited"`
	Value    int  `json:"value"`
}

type globalState struct {
	Node   *maelstrom.Node
	lin_kv *maelstrom.KV
}

func NewGlobalState() *globalState {
	n := maelstrom.NewNode()
	return &globalState{
		Node:   n,
		lin_kv: maelstrom.NewLinKV(n),
	}
}

// GLOBAL STATE VARIABLE
var (
	GlobalState        = NewGlobalState()
	CONTEXT_EXPIRES_IN = 2
	MAX_RETRIES        = 10
)

func isKeyNotExistError(err error) bool {
	if rpcErr, ok := err.(*maelstrom.RPCError); ok && rpcErr.Code == maelstrom.KeyDoesNotExist {
		return true
	}
	return false
}

func (n *globalState) getAndConvertLogList(key string, ctx context.Context) ([]log, any, error) {
	var logList []log

	logListAny, err := n.lin_kv.Read(ctx, key)
	if err != nil {
		if isKeyNotExistError(err) {
			return []log{}, nil, nil
		}
		return nil, nil, err

	}
	if logListAny != nil {
		typedList, ok := logListAny.([]any)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected type for key %s: %T", key, logListAny)
		}

		// Convert from generic interface{} slice to our log slice
		logList = make([]log, len(typedList))
		for i, item := range typedList {
			if m, ok := item.(map[string]any); ok {
				// Extract values with proper type conversion
				offset, _ := m["offset"].(float64)
				committed, _ := m["commited"].(bool)
				val, _ := m["value"].(float64)

				logList[i] = log{
					Offset:   int(offset),
					Commited: committed,
					Value:    int(val),
				}
			}
		}

	} else {
		logList = []log{}
	}

	return logList, logListAny, nil
}

func (n *globalState) addLog(key string, value int) (log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CONTEXT_EXPIRES_IN)*time.Second)
	defer cancel()

	for retry := range MAX_RETRIES {
		logList, oldValue, err := n.getAndConvertLogList(key, ctx)
		if err != nil {
			return log{}, err
		}

		newLog := log{
			Offset:   len(logList),
			Commited: false,
			Value:    value,
		}

		newLogList := append(logList, newLog)

		err = n.lin_kv.CompareAndSwap(ctx, key, oldValue, newLogList, true)
		if err == nil {
			return newLog, nil
		}

		time.Sleep(50 * time.Millisecond * time.Duration(retry+1))
	}
	return log{}, fmt.Errorf("failed to add log after %d attempts: conflict detected", MAX_RETRIES)
}

func (n *globalState) getLog(key string, offset int) (log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CONTEXT_EXPIRES_IN)*time.Second)
	defer cancel()

	logList, _, err := n.getAndConvertLogList(key, ctx)
	if err != nil {
		return log{}, err
	}

	if offset >= len(logList) {
		return log{}, fmt.Errorf("the offset %d is out of range", offset)
	}

	return logList[offset], nil
}

func (n *globalState) commitOffset(key string, upToOffset int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CONTEXT_EXPIRES_IN)*time.Second)
	defer cancel()

	logList, oldValue, err := n.getAndConvertLogList(key, ctx)
	if err != nil {
		return fmt.Errorf("unexpected type for key %s", key)
	}

	if upToOffset >= len(logList) {
		return fmt.Errorf("there is no value in the log for the given offset: %d", upToOffset)
	}

	for i := 0; i <= upToOffset; i++ {
		logList[i].Commited = true
	}

	return n.lin_kv.CompareAndSwap(ctx, key, oldValue, logList, false)
}

func (n *globalState) getHighestCommitedOffset(key string) (int, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CONTEXT_EXPIRES_IN)*time.Second)
	defer cancel()

	logList, _, err := n.getAndConvertLogList(key, ctx)
	if err != nil {
		return 0, false
	}

	highestOffset := -1
	left, right := 0, len(logList)-1

	for left <= right {
		mid := (left + right) / 2

		if logList[mid].Commited {
			highestOffset = logList[mid].Offset
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
