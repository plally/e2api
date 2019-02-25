package main



type apiKey struct {
	keyHash string
	expiration int32
	parentKey int32
}


func isValidKey(key string) bool {
	return true
}