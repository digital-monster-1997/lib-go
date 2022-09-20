package dmap

import (
	"fmt"
	"github.com/DigitakMonster1997/lib-go/pkg/utils/dcast"
)

// lookup Flatten a multilevel map
func lookup(prefix, sep string, sourceData, destData map[string]interface{}) {
	for key, value := range sourceData {
		fullIndex := fmt.Sprintf("%s%s%s", prefix, sep, key)
		if prefix == "" {
			fullIndex = fmt.Sprintf("%s", key)
		}
		if dd, err := dcast.ToStringMapE(value); err != nil {
			// 最後一層了,無法分割
			destData[fullIndex] = value
		} else {
			// 有下一層，再繼續下去
			lookup(fullIndex, sep, dd, destData)
		}
	}
}
