package storageunit

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Bits uint64

const (
	Bit     Bits = 1
	Kilobit Bits = 1000
	Megabit Bits = 1000 * Kilobit
	Gigabit Bits = 1000 * Megabit
	Terabit Bits = 1000 * Gigabit
	Petabit Bits = 1000 * Terabit
	Exabit  Bits = 1000 * Petabit

	MaxBitUnit Bits = 1<<64 - 1
	MinBitUnit Bits = 0
)

func ParseBits(s string) (Bits, error) {
	if len(s) == 0 {
		return 0, nil
	}
	if len(s) < 2 || s[len(s)-1] != 'b' {
		return 0, errors.New("Invalid Bit")
	}
	u := Bit
	s = s[:len(s)-1]
	switch s[len(s)-1] {
	case 'K':
		u = Kilobit
		s = s[:len(s)-1]
	case 'M':
		u = Megabit
		s = s[:len(s)-1]
	case 'G':
		u = Gigabit
		s = s[:len(s)-1]
	case 'T':
		u = Terabit
		s = s[:len(s)-1]
	case 'P':
		u = Petabit
		s = s[:len(s)-1]
	case 'E':
		u = Exabit
		s = s[:len(s)-1]
	default:
		if s[len(s)-1] < '0' || s[len(s)-1] > '9' {
			return 0, errors.New("Invalid Bit")
		}
	}
	n, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return 0, err
	}
	return Bits(n * float64(u)), nil
}
func (b Bits) String() string {
	// Largest unit is Exabit
	if b >= Exabit {
		return strconv.FormatFloat(float64(b/Exabit), 'f', -1, 64) + "Eb"
	} else if b >= Petabit {
		return strconv.FormatFloat(float64(b/Petabit), 'f', -1, 64) + "Pb"
	} else if b >= Terabit {
		return strconv.FormatFloat(float64(b/Terabit), 'f', -1, 64) + "Tb"
	} else if b >= Gigabit {
		return strconv.FormatFloat(float64(b/Gigabit), 'f', -1, 64) + "Gb"
	} else if b >= Megabit {
		return strconv.FormatFloat(float64(b/Megabit), 'f', -1, 64) + "Mb"
	} else if b >= Kilobit {
		return strconv.FormatFloat(float64(b/Kilobit), 'f', -1, 64) + "Kb"
	} else {
		return strconv.FormatInt(int64(b), 10) + "b"
	}
}
func (b Bits) Byte() Bytes {
	return Bytes(b)
}
func (b Bits) Kilobits() float64 {
	return float64(b) / float64(Kilobit)
}
func (b Bits) Megabits() float64 {
	return float64(b) / float64(Megabit)
}
func (b Bits) Gigabits() float64 {
	return float64(b) / float64(Gigabit)
}
func (b Bits) Terabits() float64 {
	return float64(b) / float64(Terabit)
}
func (b Bits) Petabits() float64 {
	return float64(b) / float64(Petabit)
}
func (b Bits) Exabits() float64 {
	return float64(b) / float64(Exabit)
}

func (b *Bits) UnmarshalText(bs []byte) error {
	n, err := ParseBits(string(bs))
	if err != nil {
		return err
	}
	*b = n
	return nil
}
func (b Bits) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}
func (b Bits) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}
func (b *Bits) UnmarshalJSON(bs []byte) error {
	return b.UnmarshalText(bs[1 : len(bs)-1])
}

func StringToBitsHookFunc() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflectTypeOfBit {
			return data, nil
		}
		return ParseBits(data.(string))
	}
}

var (
	reflectTypeOfBit = reflect.TypeOf(Kilobit)
)
