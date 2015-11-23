package hasher

import "crypto/md5"
import "crypto/sha1"

func Md5(bytes []byte) [16]byte {
	return md5.Sum(bytes)
}

func Sha1(bytes []byte) [20]byte {
	return sha1.Sum(bytes)
}

func Hash(bytes []byte) ([20]byte, [16]byte) {
	return Sha1(bytes), Md5(bytes)
}