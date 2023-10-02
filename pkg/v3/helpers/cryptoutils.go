package helpers

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/lamassuiot/lamassuiot/pkg/v3/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func GenerateCertificateRequest(subject models.Subject, key *rsa.PrivateKey) (*x509.CertificateRequest, error) {

	template := x509.CertificateRequest{
		Subject: SubjectToPkixName(subject),
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, err
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, err
	}

	return csr, err
}

func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return privkey, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func LoadSytemCACertPool() *x509.CertPool {
	certPool := x509.NewCertPool()
	systemCertPool, err := x509.SystemCertPool()
	if err == nil {
		certPool = systemCertPool
	} else {
		log.Warnf("could not get system cert pool (trusted CAs). Using empty pool: %s", err)
	}

	return certPool
}

func LoadSystemCACertPoolWithExtraCAsFromFiles(casToAdd []string) *x509.CertPool {
	certPool := x509.NewCertPool()
	systemCertPool, err := x509.SystemCertPool()
	if err == nil {
		certPool = systemCertPool
	} else {
		log.Warnf("could not get system cert pool (trusted CAs). Using empty pool: %s", err)
	}

	for _, ca := range casToAdd {
		if ca == "" {
			continue
		}

		caCert, err := ReadCertificateFromFile(ca)
		if err != nil {
			log.Warnf("could not load CA certificate in %s. Skipping CA: %s", ca, err)
			continue
		}

		certPool.AddCert(caCert)
	}

	return certPool
}

func ValidateCertAndPrivKey(cert *x509.Certificate, rsaKey *rsa.PrivateKey, ecKey *ecdsa.PrivateKey) (bool, error) {
	errs := []string{
		"tls: private key type does not match public key type",
		"tls: private key does not match public key",
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if rsaKey != nil {
		keyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)
		_, err := tls.X509KeyPair(pemCert, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}))
		if err == nil {
			return true, nil
		}

		contains := slices.Contains(errs, err.Error())
		if contains {
			return false, nil
		}

		return false, err
	}

	if ecKey != nil {
		keyBytes, err := x509.MarshalECPrivateKey(ecKey)
		if err != nil {
			return false, err
		}

		_, err = tls.X509KeyPair(pemCert, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}))
		if err == nil {
			return true, nil
		}

		contains := slices.Contains(errs, err.Error())
		if contains {
			return false, nil
		}

		return false, err
	}

	return false, fmt.Errorf("both keys are nil")
}

func CalculateRSAKeySizes(keyMin int, KeyMax int) []int {
	var keySizes []int
	key := keyMin
	for {
		if key%128 == 0 {
			keySizes = append(keySizes, key)
			key = key + 128
		}
		if key%1024 == 0 {
			break
		}
	}
	for {
		if key%1024 == 0 {
			keySizes = append(keySizes, key)
			if key == KeyMax {
				break
			}
			key = key + 1024
		} else {
			break
		}
	}
	return keySizes
}

func CalculateECDSAKeySizes(keyMin int, KeyMax int) []int {
	var keySizes []int
	keySizes = append(keySizes, keyMin)
	if keyMin < 224 && KeyMax > 224 {
		keySizes = append(keySizes, 224)
	}
	if keyMin < 256 && KeyMax > 256 {
		keySizes = append(keySizes, 256)
	}
	if keyMin < 384 && KeyMax > 384 {
		keySizes = append(keySizes, 384)
	}
	if keyMin < 512 && KeyMax >= 512 {
		keySizes = append(keySizes, 512)
	}
	return keySizes
}
