package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"testing"
)

func TestEncrypt(t *testing.T) {
	chatKey := "test_key"
	signatureStr := "7ef68ca6303fb25a83044abdcd7946cb9dca6f83b72163b40741ecaf24a4e7da388dd855ee79d5fdcda95a1e58c88aae6ec60e9dd88748800bb9f84b0c720ff900"
	encrypt := func(src []byte) ([]byte, error) {
		nonce := sha1.Sum([]byte(chatKey))
		signature, err := hex.DecodeString(signatureStr)
		if err != nil {
			return nil, err
		}

		nonce12 := nonce[:chacha20poly1305.NonceSize]
		key32 := signature[:chacha20poly1305.KeySize]
		aead, err := chacha20poly1305.New(key32)
		if err != nil {
			return nil, err
		}
		src = aead.Seal(nil, nonce12, src, nil)
		dst := make([]byte, hex.EncodedLen(len(src)))
		hex.Encode(dst, src)
		return dst, nil
	}
	decrypt := func(src []byte) ([]byte, error) {
		nonce := sha1.Sum([]byte(chatKey))
		signature, err := hex.DecodeString(signatureStr)
		if err != nil {
			return nil, err
		}
		nonce12 := nonce[:chacha20poly1305.NonceSize]
		key32 := signature[:chacha20poly1305.KeySize]
		aead, err := chacha20poly1305.New(key32)
		if err != nil {
			return nil, err
		}

		dst1 := make([]byte, hex.DecodedLen(len(src)))

		hex.Decode(dst1, src)
		return aead.Open(nil, nonce12, dst1, nil)
	}
	ciperText, err := encrypt([]byte("test"))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	plainText, err := decrypt(ciperText)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(string(plainText))
}