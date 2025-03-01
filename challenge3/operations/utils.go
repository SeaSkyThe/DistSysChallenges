package operations

import (
	"crypto/rand"
	"math/big"
	"sync"
)

var (
	NEIGHBOR_NODES []string = []string{}

	seenMessages = struct {
		sync.RWMutex
		msgs map[int]bool
	}{msgs: make(map[int]bool)}

	messageStorage = struct {
		sync.RWMutex
		msgs []int
	}{msgs: []int{}}
)

func generateRandomMessageID() (int, error) {
	idBig, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return 0, err
	}
	return int(idBig.Int64()) + 1, nil
}

func wasSeen(messageValue int) bool {
	seenMessages.RLock()
	defer seenMessages.RUnlock()

	_, seen := seenMessages.msgs[messageValue]
	return seen
}

func storeMessage(messageValue int) {
	messageStorage.Lock()
	seenMessages.Lock()
	defer seenMessages.Unlock()
	defer messageStorage.Unlock()

	messageStorage.msgs = append(messageStorage.msgs, messageValue)
	seenMessages.msgs[messageValue] = true
}

func storeMessages(messageValues []int) {
	messageStorage.Lock()
	seenMessages.Lock()
	defer messageStorage.Unlock()
	defer seenMessages.Unlock()

	for _, msgValue := range messageValues {
		if !seenMessages.msgs[msgValue] {
			messageStorage.msgs = append(messageStorage.msgs, msgValue)
			seenMessages.msgs[msgValue] = true
		}
	}
}

func readMessages() []int {
	// Lets make a copy
	// So, we can avoid using the original messageStorage.msg
	// and to risk a problem/state change when we are sending the response
	// Think about if we modify the messageStorage in between preparing the repsonse and sending it
	//
	messageStorage.RLock()
	defer messageStorage.RUnlock()

	return append([]int(nil), messageStorage.msgs...)
}

func clearStorages() {
	// Clear seenMessages
	seenMessages.Lock()
	seenMessages.msgs = make(map[int]bool)
	seenMessages.Unlock()

	// Clear messageStorage
	messageStorage.Lock()
	messageStorage.msgs = []int{}
	messageStorage.Unlock()
}
