package vault

func DetectContentType(sample []byte) (string, error) {
	return detectContentType(sample)
}
