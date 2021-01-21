package log

import (
	"fmt"
	"github.com/kiyoptr/su/log/tagprovider"
	"sort"
	"strings"
)

type tagorder int

const (
	toName tagorder = iota
	toMode
	toCustom
	toMessage tagorder = 1000
)

type taglist struct {
	list        map[tagorder]tagprovider.Provider
	customCount int
}

func newTagList() taglist {
	return taglist{
		list: make(map[tagorder]tagprovider.Provider),
	}
}

func (l taglist) set(order tagorder, provider tagprovider.Provider) {
	if order == toCustom {
		l.list[toCustom+tagorder(l.customCount)] = provider
		l.customCount++
	} else {
		l.list[order] = provider
	}
}

func (l taglist) sortKeys() []int {
	orders := make([]int, len(l.list))
	i := 0
	for o := range l.list {
		orders[i] = int(o)
		i++
	}
	sort.Ints(orders)

	return orders
}

func (l taglist) build() string {
	orders := l.sortKeys()

	sb := &strings.Builder{}
	for i, o := range orders {
		key, value := l.list[tagorder(o)]()
		fmt.Fprintf(sb, "[%s=%s]", key, value)

		if i+1 < len(orders) {
			sb.WriteByte(' ')
		}
	}

	return sb.String()
}

func (l taglist) merge(in taglist) taglist {
	newList := newTagList()

	lastCustomOrder := toCustom
	for k, v := range l.list {
		newList.list[k] = v
		if k >= toCustom && k < toMessage {
			lastCustomOrder = k
		}
	}

	for k, v := range in.list {
		if k >= toCustom && k < toMessage {
			k = lastCustomOrder + (k - toCustom + 1)
		}
		newList.list[k] = v
	}

	orders := newList.sortKeys()
	for _, o := range orders {
		if o >= int(toCustom) && o < int(toMessage) {
			newList.customCount++
		}
	}

	return newList
}
