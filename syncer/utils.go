package syncer

import "strconv"

func MustUint64(value string) int64 {
	r, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return r
}
