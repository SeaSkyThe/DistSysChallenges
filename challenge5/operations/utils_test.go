package operations

import (
	"reflect"
	"testing"
)

func Test_nodeState_addLog(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int
		want  log
	}{
		{
			name:  "Check if log is added correctly",
			key:   "k1",
			value: 10,
			want:  log{offset: 0, commited: false, value: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNodeState()
			got := n.addLog(tt.key, tt.value)
			if tt.want != got {
				t.Errorf("addLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeState_getLog(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		offset  int
		want    log
		wantErr bool
	}{
		{
			name:    "Check if log is consulted correctly",
			key:     "k1",
			offset:  0,
			want:    log{offset: 0, commited: false, value: 100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNodeState()
			n.addLog(tt.key, 100)
			got, gotErr := n.getLog(tt.key, tt.offset)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getLog() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getLog() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("getLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeState_getLogs(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		offset  int
		want    []log
		wantErr bool
	}{
		{
			name:   "Try to get all logs after specific offset",
			key:    "k1",
			offset: 1,
			want: []log{
				{offset: 1, commited: false, value: 2},
				{offset: 2, commited: false, value: 3},
				{offset: 3, commited: false, value: 4},
				{offset: 4, commited: false, value: 5},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNodeState()
			n.addLog("k1", 1)
			n.addLog("k1", 2)
			n.addLog("k1", 3)
			n.addLog("k1", 4)
			n.addLog("k1", 5)
			got, gotErr := n.getLogs(tt.key, tt.offset)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getLogs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getLogs() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeState_commitOffset(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		upToOffset int
		want       []log
		wantErr    bool
	}{
		{
			name:       "Check commit up until offset",
			key:        "k1",
			upToOffset: 3,
			want: []log{
				{offset: 0, commited: true, value: 1},
				{offset: 1, commited: true, value: 2},
				{offset: 2, commited: true, value: 3},
				{offset: 3, commited: true, value: 4},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNodeState()
			n.addLog("k1", 1)
			n.addLog("k1", 2)
			n.addLog("k1", 3)
			n.addLog("k1", 4)
			n.addLog("k1", 5)

			gotErr := n.commitOffset(tt.key, tt.upToOffset)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("commitOffset() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("commitOffset() succeeded unexpectedly")
			}

			got := n.logs[tt.key][:tt.upToOffset+1]
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeState_getHighestCommitedOffset(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		want  int
		want2 bool
	}{
		{
			name:  "Check highest commited offset",
			key:   "k1",
			want:  4,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNodeState()
			n.addLog("k1", 1)
			n.addLog("k1", 2)
			n.addLog("k1", 3)
			n.addLog("k1", 4)
			n.addLog("k1", 5)
			n.addLog("k1", 6)
			n.addLog("k1", 7)
			n.addLog("k1", 8)
			n.addLog("k1", 9)
			n.addLog("k1", 10)
			n.addLog("k1", 11)
			n.addLog("k1", 12)
			n.addLog("k1", 13)
			n.addLog("k1", 14)
			n.commitOffset(tt.key, 4)

			got, got2 := n.getHighestCommitedOffset(tt.key)
			if tt.want != got {
				t.Errorf("getHighestCommitedOffset() = %v, want %v", got, tt.want)
			}
			if !tt.want2 {
				t.Errorf("getHighestCommitedOffset() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
