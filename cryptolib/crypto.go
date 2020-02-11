package cryptolib

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/binary"
    "errors"
    "github.com/codahale/sss"
    "io"
    "fmt"
    "strconv"
)

func Encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    ciphertext := make([]byte, aes.BlockSize+len(text))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
    return ciphertext, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(text) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    return text, nil
}

func convertElement(elem []byte, idx byte, lengthOfIndex int) []byte {
    // Base64 convert
    elemConverted := []byte(base64.StdEncoding.EncodeToString(elem))

    // cut "=" at the end
    for {
        if elemConverted[len(elemConverted)-1] == []byte("=")[0] {
            elemConverted = elemConverted[:len(elemConverted)-1]
        } else {
            break
        }
    }

    // Append index to string
    indexFormat := "%0" + fmt.Sprintf("%d",lengthOfIndex) + "d"
    elemConverted = append(elemConverted[:], []byte(fmt.Sprintf(indexFormat,idx))[:]...)
    return elemConverted
}


func ShareSecret(secret string, numberOfShares byte, minimumKeys byte, lengthOfIndex int) ([]string, error) {
    // Check length
    if len([]byte(fmt.Sprintf("%d",int(numberOfShares)))) > len([]byte(fmt.Sprintf("%"+fmt.Sprintf("%d",lengthOfIndex)+"d",1))) {
        return []string{},errors.New("TooShortLengthOfIndexForSuchNumberOfShares")
    }
    //Check Secret
    if len(secret) < 1 {
        return []string{},errors.New("TooShortStringToShare")
    }
    var arrayCodes []string
    // Split Secret
    sharesUnparsed, err := sss.Split(numberOfShares, minimumKeys, []byte(secret))
    // Convert secret to our format
    for index,element := range sharesUnparsed {
        arrayCodes = append(arrayCodes[:], string(convertElement(element, index, lengthOfIndex)))
    }
    return arrayCodes, err
}


func RecoverSecret(strings []string, minimumShares byte, lengthOfIndex int) (string, error) {
    shares := make(map[byte][]byte, minimumShares)
    if len(strings) < int(minimumShares) {
        return "",errors.New("MaloElementovDlyaDecodirovaniya")
    }

    // CUT main array to minimum
    for {
        if len(strings) > int(minimumShares) {
            strings = strings[1:]
        } else {
            break
        }
    }

    // Create array of shares
    for _, element := range strings {
        elementByte := []byte(element)

        //Get Index of share
        indexByte:= make([]byte, 100)
        indexInt,err := strconv.Atoi(string(elementByte[len(elementByte)-lengthOfIndex:]))
        if err != nil {
            return "", errors.New("CannotParseStringToInteger")
        }
        binary.LittleEndian.PutUint32(indexByte, uint32(indexInt))

        //Get Share
        elementByte = elementByte[:len(elementByte)-lengthOfIndex]

        //Fill share with "=" to multiple of 4
        switch {
        case (len(elementByte) % 4) == 1:
            elementByte = append(elementByte[:], []byte("===")[:]...)
        case (len(elementByte) % 4) == 2 :
            elementByte = append(elementByte[:], []byte("==")[:]...)
        case (len(elementByte) % 4) == 3 :
            elementByte = append(elementByte[:], []byte("=")[:]...)
        }

        // Base64 unconvert
        elementString,err := base64.StdEncoding.DecodeString(string(elementByte))
        if err != nil {
            return "", errors.New("CannotDecodeBase64string")
        }

        // Fill array of shares
        shares[indexByte[0]] = elementString
    }

    // DECODE
    recovered := string(sss.Combine(shares))
    return recovered, nil
}

