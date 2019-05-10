package aws

// Add cfctl stamp to the stack tags
func tagPkgStamp(tags map[string]string) map[string]string {
	if len(tags) == 0 {
		tags = make(map[string]string)
	}

	tags["CreatedBy"] = "cfctl"

	return tags
}
