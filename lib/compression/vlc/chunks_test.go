package vlc

import (
	"reflect"
	"testing"
)

func Test_splitByChanks(t *testing.T) {
	type args struct {
		bStr      string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want BinaryChunks
	}{
		{
			name: "Test 1",
			args: args{
				bStr:      "001000100110100101",
				chunkSize: 8,
			},
			want: BinaryChunks{
				"00100010",
				"01101001",
				"01000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByChanks(tt.args.bStr, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByChanks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryChunks_Join(t *testing.T) {
	tests := []struct {
		name string
		bcs  BinaryChunks
		want string
	}{
		{
			name: "Test 1",
			bcs:  BinaryChunks{"00100010", "01101001", "01000000"},
			want: "001000100110100101000000",
		},
		{
			name: "Test 2",
			bcs:  BinaryChunks{"00100010", "011010"},
			want: "00100010011010",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcs.Join(); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBinChunks(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want BinaryChunks
	}{
		{
			name: "Test 1",
			data: []byte{20, 30, 60, 18},
			want: BinaryChunks{"00010100", "00011110", "00111100", "00010010"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBinChunks(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBinChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}
