package go_opus

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestOpusDecoder(t *testing.T) {
	channels := 1
	sampleRate := 16000
	decoder, err := NewDecoder(sampleRate, channels)
	if err != nil {
		t.Errorf("Create decoder error %v", err)
	}

	inData := []byte{
		120, 0, 178, 76, 67,
		73, 253, 128, 165,
		248, 119, 114, 100,
		237, 112, 189, 245,
		180, 75,
	}
	resultLen := 320
	outData := make([]int16, int(resultLen))
	n, err := decoder.Decode(inData, outData)

	if n != resultLen {
		t.Errorf("Decode result error %v", n)
		return
	}
	if err != nil {
		t.Errorf("Decode error %v", err)
		return
	}
	fmt.Println(strings.ReplaceAll(fmt.Sprintf("%v", outData), " ", ","))
	validOut := []int16{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0,
		-1, -1, 0, 1, -1, 0, 1, 0, -1, -1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, -1, -1, 0, 0,
		0, -1, 0, 1, 1, 0, 0, 0, 0, -1, -1, 0, 1, 0, 0, 0, 0, 0, -1, -1, 0, 0, 0, 0, 1, 0,
		0, 0, -1, -1, -1, 0, 0, 0, 1, 1, 0, 0, 1, 1, -1, 0, 1, 1, 0, 0, 2, 0, 0, 0, 0, -6,
		-5, 4, 2, -2, -1, 2, 4, 5, 1, 0, 2, 1, 7, 4, 1, 2, 4, -3, -5, -2, -2, -1, -2, -1,
		1, -5, -5, -1, -1, -8, -13, -4, 7, 8, 1, 1, 11, 9, 2, 0, 2, -3, -3, 0, 2, 1, -1, 2,
		3, 0, -2, -6, -5, -9, -7, -1, -1, 0, 1, 2, 2, 6, 3, 1, 1, 1, 0, 1, 1, 1, 2, -5, -9, 3, 4,
		-2, -2, -6, -2, 0, -2, -3, -1, -4, -2, 1, 0, -7, -4, 2, 1, 4, 1, -4, 1, 1, 0, -6, -4, 1, 1,
		0, -1, 0, 3, 1, -1, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 1, -5, -4, -1, 0, 0, -2, -1, 0, 5, 3, -6, -3,
		1, 1, -2, -3, 2, 8, 4, 0, 2, 2, 0, -1, 0, -1, -1, -5, -2, 2, -1, -1, -1, 0, -1,
		-1, 0, -5, -3, 1, 7, 9, 2, 1, 10, 7, -5, 0, 4, 2, -2, -2, 8, 6, 0, 6, 6, 1, -2, -2, 0, -8, -12, -5,
		-5, -6, -9, -5, 1, 0, -2, -2, 0, -1, -2, -2, 0, 2, 1, 2, 3, 2, 2, 2, 2, 0, -1, -1, -5, -4, 0, 0, 0, -1, -1, 2, 0,
	}
	if !reflect.DeepEqual(outData, validOut) {
		t.Errorf("Decoded result not valid")
	}
}
