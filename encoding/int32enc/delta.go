package int32enc

func EncodeDeltaInplace(values []uint32) {
	var currentValue uint32
	for i := range values {
		if i == 0 {
			currentValue = values[i]
			continue
		}
		values[i], currentValue = values[i]-currentValue, values[i]
	}
}

func DecodeDeltaInplace(values []uint32) {
	var currentValue uint32
	for i := range values {
		if i == 0 {
			currentValue = values[i]
			continue
		}
		values[i] = values[i] + currentValue
		currentValue = values[i]
	}
}
