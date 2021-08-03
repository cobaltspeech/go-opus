package go_opus

import (
	"fmt"
	"unsafe"
)

/*

#cgo linux,!android,!musl,arm LDFLAGS: ${SRCDIR}/lib/libopus.gnu.arm.a -lpthread
#cgo linux,!android,!musl,arm64 LDFLAGS: ${SRCDIR}/lib/libopus.gnu.aarch64.a -lpthread
#cgo linux,!android,!musl,amd64 LDFLAGS: ${SRCDIR}/lib/libopus.gnu.x86_64.a -lpthread
#cgo linux,!android,musl,arm LDFLAGS: ${SRCDIR}/lib/libopus.musl.arm.a
#cgo linux,!android,musl,arm64 LDFLAGS: ${SRCDIR}/lib/libopus.musl.aarch64.a
#cgo linux,!android,musl,amd64 LDFLAGS: ${SRCDIR}/lib/libopus.musl.x86_64.a
#cgo android,arm LDFLAGS: ${SRCDIR}/lib/libopus.android.a
#cgo android,arm64 LDFLAGS: ${SRCDIR}/lib/libopus.android.a
#cgo android,amd64 LDFLAGS: ${SRCDIR}/lib/libopus.android.a
#cgo ios,arm64,mobile LDFLAGS: ${SRCDIR}/lib/libopus.ios.arm64.a -framework Accelerate
#cgo ios,amd64,mobile LDFLAGS: ${SRCDIR}/lib/libopus.ios.simulator.x86_64.a -framework Accelerate
#cgo darwin,arm64,!mobile LDFLAGS: ${SRCDIR}/lib/libopus.macos.arm64.a -framework Accelerate
#cgo darwin,amd64,!mobile LDFLAGS: ${SRCDIR}/lib/libopus.macosx.x86_64.a -framework Accelerate
#cgo darwin LDFLAGS: ${SRCDIR}/lib/libopus.macosx.x86_64.a  -lm

#include "opus.h"

int
bridge_decoder_get_last_packet_duration(OpusDecoder *st, opus_int32 *samples)
{
	return opus_decoder_ctl(st, OPUS_GET_LAST_PACKET_DURATION(samples));
}
*/
import "C"

var errDecUninitialized = fmt.Errorf("opus decoder uninitialized")

type Decoder interface {
	Decode(data []byte, pcm []int16) (int, error)
}

type opusDecoder struct {
	p *C.struct_OpusDecoder
	// Same purpose as encoder struct
	mem        []byte
	sampleRate int
	channels   int
}

// NewDecoder allocates a new Opus decoder and initializes it with the
// appropriate parameters. All related memory is managed by the Go GC.
func NewDecoder(sampleRate, channels int) (Decoder, error) {
	var dec opusDecoder
	err := dec.Init(sampleRate, channels)

	if err != nil {
		return nil, err
	}

	return &dec, nil
}

func (dec *opusDecoder) Init(sampleRate, channels int) error {
	if dec.p != nil {
		return fmt.Errorf("opus decoder already initialized")
	}

	if channels != 1 && channels != 2 {
		return fmt.Errorf("number of channels must be 1 or 2: %d", channels)
	}

	size := C.opus_decoder_get_size(C.int(channels))
	dec.sampleRate = sampleRate
	dec.channels = channels
	dec.mem = make([]byte, size)
	dec.p = (*C.OpusDecoder)(unsafe.Pointer(&dec.mem[0]))
	errno := C.opus_decoder_init(dec.p, C.opus_int32(sampleRate), C.int(channels))

	if errno != 0 {
		return fmt.Errorf("%v", errno)
	}

	return nil
}

// Decode encoded Opus data into the supplied buffer. On success, returns the
// number of samples correctly written to the target buffer.
func (dec *opusDecoder) Decode(data []byte, pcm []int16) (int, error) {
	if dec.p == nil {
		return 0, errDecUninitialized
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}

	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return 0, fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))

	if n < 0 {
		return 0, fmt.Errorf("%v", n)
	}

	return n, nil
}

// Decode encoded Opus data into the supplied buffer. On success, returns the
// number of samples correctly written to the target buffer.
func (dec *opusDecoder) DecodeFloat32(data []byte, pcm []float32) (int, error) {
	if dec.p == nil {
		return 0, errDecUninitialized
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}

	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return 0, fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode_float(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))

	if n < 0 {
		return 0, fmt.Errorf("%v", n)
	}

	return n, nil
}

// DecodeFEC encoded Opus data into the supplied buffer with forward error
// correction.
//
// It is to be used on the packet directly following the lost one.  The supplied
// buffer needs to be exactly the duration of audio that is missing
//
// When a packet is considered "lost", DecodeFEC can be called on the next
// packet in order to try and recover some of the lost data. The PCM needs to be
// exactly the duration of audio that is missing.  `LastPacketDuration()` can be
// used on the decoder to get the length of the last packet.  Note also that in
// order to use this feature the encoder needs to be configured with
// SetInBandFEC(true) and SetPacketLossPerc(x) options.
//
// Note that DecodeFEC automatically falls back to PLC when no FEC data is
// available in the provided packet.
func (dec *opusDecoder) DecodeFEC(data []byte, pcm []int16) error {
	if dec.p == nil {
		return errDecUninitialized
	}

	if len(data) == 0 {
		return fmt.Errorf("opus: no data supplied")
	}

	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		1))

	if n < 0 {
		return fmt.Errorf("%v", n)
	}

	return nil
}

// DecodeFECFloat32 encoded Opus data into the supplied buffer with forward error
// correction. It is to be used on the packet directly following the lost one.
// The supplied buffer needs to be exactly the duration of audio that is missing
func (dec *opusDecoder) DecodeFECFloat32(data []byte, pcm []float32) error {
	if dec.p == nil {
		return errDecUninitialized
	}

	if len(data) == 0 {
		return fmt.Errorf("opus: no data supplied")
	}

	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode_float(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		1))

	if n < 0 {
		return fmt.Errorf("%v", n)
	}

	return nil
}

// DecodePLC recovers a lost packet using Opus Packet Loss Concealment feature.
//
// The supplied buffer needs to be exactly the duration of audio that is missing.
// When a packet is considered "lost", `DecodePLC` and `DecodePLCFloat32` methods
// can be called in order to obtain something better sounding than just silence.
// The PCM needs to be exactly the duration of audio that is missing.
// `LastPacketDuration()` can be used on the decoder to get the length of the
// last packet.
//
// This option does not require any additional encoder options. Unlike FEC,
// PLC does not introduce additional latency. It is calculated from the previous
// packet, not from the next one.
func (dec *opusDecoder) DecodePLC(pcm []int16) error {
	if dec.p == nil {
		return errDecUninitialized
	}

	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: output buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode(
		dec.p,
		nil,
		0,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))

	if n < 0 {
		return fmt.Errorf("%v", n)
	}

	return nil
}

// DecodePLCFloat32 recovers a lost packet using Opus Packet Loss Concealment feature.
// The supplied buffer needs to be exactly the duration of audio that is missing.
func (dec *opusDecoder) DecodePLCFloat32(pcm []float32) error {
	if dec.p == nil {
		return errDecUninitialized
	}

	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}

	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: output buffer capacity must be multiple of channels")
	}

	n := int(C.opus_decode_float(
		dec.p,
		nil,
		0,
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))

	if n < 0 {
		return fmt.Errorf("%v", n)
	}

	return nil
}

// LastPacketDuration gets the duration (in samples)
// of the last packet successfully decoded or concealed.
func (dec *opusDecoder) LastPacketDuration() (int, error) {
	var samples C.opus_int32

	res := C.bridge_decoder_get_last_packet_duration(dec.p, &samples)

	if res != C.OPUS_OK {
		return 0, fmt.Errorf("%v", res)
	}

	return int(samples), nil
}

func main() {

}
