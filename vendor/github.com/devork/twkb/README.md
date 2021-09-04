# TWKB Go Library

[![Build Status](https://travis-ci.org/devork/twkb.png?branch=master)](https://travis-ci.org/devork/twkb)

A small GO parser for the [TWKB specification](https://github.com/TWKB/Specification)

# Usage

`go get github.com/devork/twkb`

```go
data, _ := hex.DecodeString("01000204")
geom, err := Decode(bytes.NewReader(data))

if err != nil {
		t.Fatalf("Failed to decode point geometry: err = %s", err)
}

// Do something magical with the geometry returned
```

The library is usable, but could do with some optimisations.
