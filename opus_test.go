// Copyright (2021) Cobalt Speech and Language Inc.

package go_opus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

//play -t raw -r 16000 -b 16 -e signed-integer ./output.pcm
func TestOpusFile(t *testing.T) {
	files, err := ioutil.ReadDir("./samples")
	if err != nil {
		t.Error(err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	channels := 1
	sampleRate := 16000
	frame := 320
	decoder, err := NewDecoder(sampleRate, channels)
	if err != nil {
		t.Errorf("Create decoder error %v", err)
	}
	outBuffer := make([]int16, frame)
	if err != nil {
		t.Error(err)
	}
	outF, err := os.Create("output.pcm")
	if err != nil {
		t.Error(err)
	}
	w := bufio.NewWriter(outF)
	defer outF.Close()
	outBytesCount := 0
	for _, f := range files {

		fileName := f.Name()
		bytes, err := ioutil.ReadFile("./samples/" + fileName)
		if err != nil {
			t.Error(err)
		}
		count, err := decoder.Decode(bytes, outBuffer)
		outBytesCount += count
		if err != nil {
			t.Errorf("Decode error %v", err)
		}
		if count != frame {
			t.Errorf("Frame is not valid: %v!=%v", count, frame)
		}
		err = binary.Write(w, binary.LittleEndian, outBuffer)
		if err != nil {
			t.Errorf("Write to file error %v", err)
		}

	}
	t.Logf("out size: %d", outBytesCount)
	origFile, err := ioutil.ReadFile("output_orig.pcm")
	if err != nil {
		t.Error(err)
	}
	outputFile, err := ioutil.ReadFile("output.pcm")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(origFile, outputFile) {
		t.Error("Files not equal")
	}
}

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
