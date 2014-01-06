package tfwf

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
