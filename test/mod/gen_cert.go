package mod

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func GenCert() {
	// 初始化CA证书和私钥
	var err error
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate CA key: %v", err)
	}

	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"MITM Proxy CA"},
			CommonName:   "MITM Proxy CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("Failed to create CA cert: %v", err)
	}

	// 保存CA证书到文件(方便用户安装)
	saveCertToFile("ca-cert.pem", caCertBytes)
}

func saveCertToFile(filename string, certBytes []byte) {
	certOut, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %v", filename, err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		log.Fatalf("Failed to write data to %s: %v", filename, err)
	}
}
