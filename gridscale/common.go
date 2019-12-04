package gridscale

//convSOStrings converts slice of interfaces to slice of strings
func convSOStrings(interfaceList []interface{}) []string {
	var labels []string
	for _, labelInterface := range interfaceList {
		labels = append(labels, labelInterface.(string))
	}
	return labels
}
