package int64enc

func EncodeDeltaInplace(values []uint64) {
	var currentValue uint64
	for i := range values {
		if i == 0 {
			currentValue = values[i]
			continue
		}
		values[i], currentValue = values[i]-currentValue, values[i]
	}
}

func DecodeDeltaInplace(values []uint64) {
	var currentValue uint64
	for i := range values {
		if i == 0 {
			currentValue = values[i]
			continue
		}
		values[i] = values[i] + currentValue
		currentValue = values[i]
	}
}
