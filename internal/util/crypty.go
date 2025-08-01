package util

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
)

func GetSHA1(cert *x509.Certificate) string {
	// 计算SHA-1指纹
	sha1Hash := sha1.Sum(cert.Raw)
	return hex.EncodeToString(sha1Hash[:])

	// 计算SHA-256指纹
	//sha256Hash := sha256.Sum256(cert.Raw)
	//sha256Fingerprint = hex.EncodeToString(sha256Hash[:])
}
