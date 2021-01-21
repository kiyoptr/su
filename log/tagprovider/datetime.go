package tagprovider

import "time"

func DateTime(format string) Provider {
	return func() (key, value string) {
		now := time.Now()
		key = "date"
		value = now.Format(format)
		return
	}
}
