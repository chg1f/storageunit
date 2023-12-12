package storageunit

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Bytes uint64

const (
	Byte     Bytes = 8
	Kilobyte Bytes = 1024 * Byte
	Megabyte Bytes = 1024 * Kilobyte
	Gigabyte Bytes = 1024 * Megabyte
	Terabyte Bytes = 1024 * Gigabyte
	Petabyte Bytes = 1024 * Terabyte
	Exabyte  Bytes = 1024 * Petabyte

	MaxByteUnit Bytes = 1<<64 - 1
	MinByteUnit Bytes = 0
)

func ParseBytes(s string) (Bytes, error) {
	if len(s) == 0 {
		return 0, nil
	}
	if len(s) < 2 || s[len(s)-1] != 'B' {
		return 0, errors.New("Invalid Byte")
	}
	u := Byte
	s = s[:len(s)-1]
	switch s[len(s)-1] {
	case 'K':
		u = Kilobyte
		s = s[:len(s)-1]
	case 'M':
		u = Megabyte
		s = s[:len(s)-1]
	case 'G':
		u = Gigabyte
		s = s[:len(s)-1]
	case 'T':
		u = Terabyte
		s = s[:len(s)-1]
	case 'P':
		u = Petabyte
		s = s[:len(s)-1]
	case 'E':
		u = Exabyte
		s = s[:len(s)-1]
	default:
		if s[len(s)-1] < '0' || s[len(s)-1] > '9' {
			return 0, errors.New("Invalid Byte")
		}
	}
	n, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return 0, err
	}
	return Bytes(n * float64(u)), nil
}
func (b Bytes) String() string {
	// Largest unit is Exabyte
	if b >= Exabyte {
		return strconv.FormatFloat(float64(b/Exabyte), 'f', -1, 64) + "EB"
	} else if b >= Petabyte {
		return strconv.FormatFloat(float64(b/Petabyte), 'f', -1, 64) + "PB"
	} else if b >= Terabyte {
		return strconv.FormatFloat(float64(b/Terabyte), 'f', -1, 64) + "TB"
	} else if b >= Gigabyte {
		return strconv.FormatFloat(float64(b/Gigabyte), 'f', -1, 64) + "GB"
	} else if b >= Megabyte {
		return strconv.FormatFloat(float64(b/Megabyte), 'f', -1, 64) + "MB"
	} else if b >= Kilobyte {
		return strconv.FormatFloat(float64(b/Kilobyte), 'f', -1, 64) + "KB"
	} else {
		return strconv.FormatUint(uint64(b), 10) + "B"
	}
}
func (b Bytes) Bit() Bits {
	return Bits(b)
}
func (b Bytes) Kilobytes() float64 {
	return float64(b) / float64(Kilobyte)
}
func (b Bytes) Megabytes() float64 {
	return float64(b) / float64(Megabyte)
}
func (b Bytes) Gigabytes() float64 {
	return float64(b) / float64(Gigabyte)
}
func (b Bytes) Terabytes() float64 {
	return float64(b) / float64(Terabyte)
}
func (b Bytes) Petabytes() float64 {
	return float64(b) / float64(Petabyte)
}
func (b Bytes) Exabytes() float64 {
	return float64(b) / float64(Exabyte)
}

func (b *Bytes) UnmarshalText(bs []byte) error {
	n, err := ParseBytes(string(bs))
	if err != nil {
		return err
	}
	*b = n
	return nil
}
func (b Bytes) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}
func (b *Bytes) UnmarshalJSON(bs []byte) error {
	return b.UnmarshalText(bs[1 : len(bs)-1])
}
func (b Bytes) MarshalJSON() ([]byte, error) {
	return []byte(`"` + b.String() + `"`), nil
}

func StringToBytesHookFunc() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflectTypeOfByte {
			return data, nil
		}
		return ParseBytes(data.(string))
	}
}

var (
	reflectTypeOfByte = reflect.TypeOf(Kilobyte)
)
