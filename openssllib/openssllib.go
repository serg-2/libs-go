package openssllib

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/serg-2/certinfo"
)

func keyToKeyPEM(key *rsa.PrivateKey, pwd string) []byte {
	// Marshalling
	keyBytes := x509.MarshalPKCS1PrivateKey(key)

	// convert it to PEM Block
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	// Encrypt Block if needed
	if pwd != "" {
		pemBlock, _ = x509.EncryptPEMBlock(rand.Reader, pemBlock.Type, pemBlock.Bytes, []byte(pwd), x509.PEMCipherAES256)
	}

	// PEM encoding of private key to []byte
	keyPEM := pem.EncodeToMemory(pemBlock)

	return keyPEM
}

func generateRSAKeyAndCertificate() (string, string, error) {
	// Key generation
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}
	// printing keyPEM
	fmt.Println(keyToKeyPEM(key, ""))

	// Time of validity
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * 10 * time.Hour)

	//Create certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: "localhost"},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	//Create certificate using template
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return "", "", err

	}
	// pem encoding of certificate
	certPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		},
	))
	fmt.Println(certPem)
	return string(keyToKeyPEM(key, "")), certPem, nil
}

func generateCertificateRequest(commonName string, clientPassword string) ([]byte, []byte) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 4096)

	encryptedKeyBytes := keyToKeyPEM(keyBytes, clientPassword)

	//var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
	//emailAddress := "test@example.com"
	subj := pkix.Name{
		CommonName:         commonName,
		Country:            []string{"AU"},
		Province:           []string{"Some-State"},
		Locality:           []string{"MyCity"},
		Organization:       []string{"Company Ltd"},
		OrganizationalUnit: []string{"IT"},
		/*
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type: oidEmailAddress,
					Value: asn1.RawValue{
						Tag:   asn1.TagIA5String,
						Bytes: []byte(emailAddress),
					},
				},
			},

		*/
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	blockBytes := &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}

	// creating buffer
	var buffer bytes.Buffer

	err := pem.Encode(&buffer, blockBytes)
	if err != nil {
		panic(err)
	}

	return encryptedKeyBytes, buffer.Bytes()
}

func dumpToFile(filename string, bytes []byte) {
	err := os.WriteFile(filename, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func crsToCrtExample(caPassword string, requestBytes []byte, caCertBytes []byte, caKeyBytes []byte, numberOfYears int, serial *big.Int) []byte {
	// Load CA public key
	pemBlock, _ := pem.Decode(caCertBytes)
	if pemBlock == nil {
		panic("Decode CA Public key failed")
	}
	caCRT, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	// Load CA private key
	pemBlock, _ = pem.Decode(caKeyBytes)
	if pemBlock == nil {
		panic("Decode CA Private key failed")
	}
	der, err := x509.DecryptPEMBlock(pemBlock, []byte(caPassword))
	if err != nil {
		fmt.Println("Error of DecryptPemBlock: ", err)
		panic("Decode CA Private key PASSWORD failed")
	}
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		panic("Parsing CA DECODED Private key PASSWORD failed")
	}

	// load client certificate request
	pemBlock, _ = pem.Decode(requestBytes)
	if pemBlock == nil {
		panic("Decode client certificate request failed")
	}
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		panic("Parsing client certificate request failed")
	}
	if err = clientCSR.CheckSignature(); err != nil {
		panic("Check signature client certificate request failed")
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,

		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,

		SerialNumber: serial,
		Issuer:       caCRT.Subject,
		Subject:      clientCSR.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Duration(numberOfYears) * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)
	if err != nil {
		panic("Create client certificate request failed.")
	}

	// creating buffer
	var buffer bytes.Buffer

	// Add String about certificate information to buffer
	buffer.WriteString(certInfo(clientCRTRaw))

	pem.Encode(&buffer, &pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})

	return buffer.Bytes()
}

func certInfo(rawCert []byte) string {
	certX509, err := x509.ParseCertificate(rawCert)
	if err != nil {
		log.Fatal(err)
	}
	res, err := certinfo.CertificateText(certX509)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func combineFolder(folderNames ...string) string {
	var a = ""
	for _, folder := range folderNames {
		a += folder
		a += "/"
	}
	a = a[:len(a)-1]

	return a
}

func GenerateCert(commonName string, clientPassword string, caPassword string, config Openssl) {
	configFileName := commonName + ".ovpn"
	userKeyFileName := commonName + ".key"
	requestFileName := commonName + ".csr"
	certFileName := commonName + ".crt"

	var certStruct CertStruct
	var err error

	// Generate Key And CertificateRequest
	keyBytes, requestBytes := generateCertificateRequest(commonName, clientPassword)

	// Dump KEY File, if needed
	dumpToFile(combineFolder(config.EasyRSAFolder, config.KeyFolder, userKeyFileName), keyBytes)
	// Dump request to File, if needed
	dumpToFile(combineFolder(config.EasyRSAFolder, config.RequestFolder, requestFileName), requestBytes)

	// Loading ca certificate
	certStruct.CABytes, err = ioutil.ReadFile(combineFolder(config.EasyRSAFolder, config.CaCertFileName))
	if err != nil {
		log.Fatal(err)
	}

	// Loading ca key
	certStruct.KeyBytes, err = ioutil.ReadFile(combineFolder(config.EasyRSAFolder, config.KeyFolder, config.CaKeyFileName))
	if err != nil {
		log.Fatal(err)
	}

	// Load SERIAL
	serialBytes, err := ioutil.ReadFile(combineFolder(config.EasyRSAFolder, config.SerialFile))
	if err != nil {
		log.Fatal(err)
	}
	// Saving old serial
	dumpToFile(combineFolder(config.EasyRSAFolder, config.SerialOldFile), serialBytes)

	// Creating name fore cert_serial folder
	certSerialFileName := strings.TrimSuffix(string(serialBytes), "\n") + ".pem"

	// Converting to big int
	serial := new(big.Int)
	serial.SetString(string(serialBytes), 16)

	certStruct.CertBytes = crsToCrtExample(caPassword, requestBytes, certStruct.CABytes, certStruct.KeyBytes, config.NumberOfYearsValidity, serial)

	// Incrementing
	serial.Add(serial, big.NewInt(1))

	// Converting to HEX string with end of line
	serialNewBytes := fmt.Sprintf("%X\n", serial)

	// Storing new serial
	dumpToFile(combineFolder(config.EasyRSAFolder, config.SerialFile), []byte(serialNewBytes))

	// Dump user certificate to File, if needed
	dumpToFile(combineFolder(config.EasyRSAFolder, config.CertFolder, certFileName), certStruct.CertBytes)
	dumpToFile(combineFolder(config.EasyRSAFolder, config.CertSerialFolder, certSerialFileName), certStruct.CertBytes)

	// Check Config directory exists
	err = os.MkdirAll(config.ConfigsFolder, os.ModePerm)
	if err != nil {
		log.Fatal("Can't create directory: " + config.ConfigsFolder)
	}

	// Reading TLS static key
	certStruct.TABytes, err = ioutil.ReadFile(config.TaFileName)
	if err != nil {
		log.Fatal(err)
	}

	WriteToFileUsingTemplate(config.TemplateFilename, combineFolder(config.ConfigsFolder, configFileName), certStruct)
}

func WriteToFileUsingTemplate(template string, fullFileName string, cert CertStruct) {
	// load template
	templateBytes, err := ioutil.ReadFile(template)
	if err != nil {
		log.Fatal(err)
	}

	// Create config file
	targetFile, err := os.Create(fullFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer targetFile.Close()

	// Writing template
	targetFile.Write(templateBytes)

	// Writing ca certificate
	targetFile.WriteString("<ca>\n")
	targetFile.Write(cert.CABytes)
	targetFile.WriteString("</ca>\n")

	// Writing certificate
	targetFile.WriteString("<cert>\n")
	targetFile.Write(cert.CertBytes)
	targetFile.WriteString("</cert>\n")

	// Writing key
	targetFile.WriteString("<key>\n")
	targetFile.Write(cert.KeyBytes)
	targetFile.WriteString("</key>\n")

	targetFile.WriteString("<tls-auth>\n")
	targetFile.Write(cert.TABytes)
	targetFile.WriteString("</tls-auth>\n")
}

func XorArrays(a []byte, b []byte) ([]byte, error) {
	if len(a) == 0 || len(b) == 0 {
		return []byte{}, errors.New("Empty array.")
	}

	if len(a) != len(b) {
		return []byte{}, errors.New("Arrays not equal.")
	}

	for i := range a {
		a[i] ^= b[i]
	}

	return a, nil
}

func ArrayFromHex(a string) ([]byte, error) {
	data, err := hex.DecodeString(a)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func ArrayToHex(a []byte) string {
	return fmt.Sprintf("%x", a)
}
