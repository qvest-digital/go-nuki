package command

import "fmt"

type Version []byte

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major(), v.Minor(), v.Patch())
}

func (v Version) Major() uint8 {
	if len(v) >= 1 {
		return v[0]
	}
	return 0
}

func (v Version) Minor() uint8 {
	if len(v) >= 2 {
		return v[1]
	}
	return 0
}

func (v Version) Patch() uint8 {
	if len(v) >= 3 {
		return v[2]
	}
	return 0
}
