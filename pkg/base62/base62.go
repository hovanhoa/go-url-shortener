package base62

import (
	"log"
)

// Base62 represents the base62 const.
const Base62 uint64 = 62

var (
	// DefaultChars character set, [A-Za-z0-9]. https://tools.ietf.org/html/rfc4648
	// We follow the orders similar to base64.
	// Note that there are several different variants in the web, such as
	// [0-9a-zA-Z] or [a-zA-Z0-9] and so on.
	DefaultChars = [Base62]rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	// Base60.
	// Alternative https://github.com/transitive-bullshit/id-shortener/blob/58eea348551eea95d662d664139abd9c13c2f449/index.js#L5
	// 0123456789ABCDEFGHJKLMNPQRSTUVWXYZ_abcdefghijkmnopqrstuvwxyz
	defaultLookup map[rune]int
)

func init() {
	defaultLookup = make(map[rune]int, Base62)
	for k, v := range DefaultChars {
		defaultLookup[v] = k
	}

}

// Factory allows the creation of new base62 encoder/decoder with different
// variant of characters.
type Factory struct {
	chars  [Base62]rune
	lookup map[rune]int
}

// New returns a new Factory with the given variant of characters.
func New(chars [Base62]rune) *Factory {
	if uint64(len(chars)) != Base62 {
		log.Fatal("length must be equal to 62")
	}
	lookup := make(map[rune]int, Base62)
	for i, c := range chars {
		lookup[c] = i
	}
	return &Factory{
		chars:  chars,
		lookup: lookup,
	}
}

// Encode encodes the given integer into a base62 string.
func (f *Factory) Encode(in uint64) string {
	return encode(f.chars, in)
}

// Decode decodes the string back into integer.
func (f *Factory) Decode(s string) uint64 {
	return decode(f.lookup, s)
}

// Encode converts the given integer into a base 62 string.
func Encode(in uint64) string {
	return encode(DefaultChars, in)
}

func encode(chars [Base62]rune, in uint64) string {
	if in < 1 {
		return ""
	}
	i, tmp := 0, in
	for tmp > 0 {
		i++
		// This will add a character if not handled.
		if tmp == Base62 {
			break
		}
		tmp /= Base62
	}

	out := make([]rune, i)
	for in > 0 {
		i--
		// Overflows when modulus 62 % 62. Last character, set it and
		// break.
		// Mod zero, we can't get the -1 index of char, shift it to last position 61.
		if in%Base62 == 0 {
			out[i] = chars[Base62-1]
			// If we divide in by 62, we will get 1, which will repeat another cycle. Terminate it.
			if in == Base62 {
				break
			}
		} else {
			out[i] = chars[in%Base62-1]
		}
		in /= Base62
	}
	return string(out)
}

// Decode attempts to convert a base 62 string back into an integer.
func Decode(s string) uint64 {
	return decode(defaultLookup, s)
}

func decode(lookup map[rune]int, s string) uint64 {
	var sum uint64
	for i, v := range s {
		val := uint64(lookup[v] + 1)
		if i > 0 {
			val = val % Base62
		}
		sum = sum*Base62 + val
	}
	return sum
}
