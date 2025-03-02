package operations

import (
	"crypto/rand"
	"math/big"
	"sync"
)

var (
	NEIGHBORHOOD []string = []string{}

	messageStorage = struct {
		sync.RWMutex
		msgs map[int]bool
	}{msgs: make(map[int]bool)}

	neighborhoodKnows = struct {
		sync.RWMutex
		neighborhood map[string]map[int]struct{}
	}{neighborhood: make(map[string]map[int]struct{})}
)

func generateRandomMessageID() (int, error) {
	idBig, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return 0, err
	}
	return int(idBig.Int64()) + 1, nil
}

func storeMessage(messageValue int) {
	messageStorage.Lock()
	defer messageStorage.Unlock()

	messageStorage.msgs[messageValue] = true
}

func storeMessages(messageValues []int) {
	messageStorage.Lock()
	defer messageStorage.Unlock()

	for _, msgValue := range messageValues {
		if _, exists := messageStorage.msgs[msgValue]; !exists {
			messageStorage.msgs[msgValue] = true
		}
	}
}

func readMessages() []int {
	// Lets make a copy
	// So, we can avoid using the original messageStorage.msg
	// and to risk a problem/state change when we are sending the response
	// Think about if we modify the messageStorage in between preparing the repsonse and sending it

	messageStorage.RLock()
	defer messageStorage.RUnlock()

	keys := make([]int, 0, len(messageStorage.msgs))
	for k := range messageStorage.msgs {
		keys = append(keys, k)
	}

	return keys
}

// Returns the one I know and the neighbor dont
func neighborAck(messages []int, src string) []int {
	neighborhoodKnows.Lock()
	defer neighborhoodKnows.Unlock()

	if _, exists := neighborhoodKnows.neighborhood[src]; !exists {
		neighborhoodKnows.neighborhood[src] = make(map[int]struct{})
	}

	unknownMessages := []int{}

	// Add each message to the source's set
	for _, msg := range messages {
		if _, exists := neighborhoodKnows.neighborhood[src][msg]; !exists {
			unknownMessages = append(unknownMessages, msg)
		}
		neighborhoodKnows.neighborhood[src][msg] = struct{}{}
	}

	return unknownMessages
}
