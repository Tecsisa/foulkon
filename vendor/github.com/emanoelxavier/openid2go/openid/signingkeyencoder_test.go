package openid

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net/http"
	"testing"
)

func Test_pemEncodePublicKey_UsingNilKey_ReturnsMarshallingKeyError(t *testing.T) {
	_, err := pemEncodePublicKey(nil)

	if err == nil {
		t.Fatal("An error was expected but not returned.")
	}

	expectValidationError(t, err, ValidationErrorMarshallingKey, http.StatusInternalServerError, nil)

}

func Test_pemEncodePublicKey_UsingRSAPublicKey(t *testing.T) {
	rsaKey := &rsa.PublicKey{N: big.NewInt(9871234), E: 15}

	ek, err := pemEncodePublicKey(rsaKey)

	if err != nil {
		t.Error("An error was not expected but returned.")
	}

	if ek == nil {
		t.Error("The encoded key should not be nil.")
	}

	pBlock, rest := pem.Decode(ek)

	if pBlock == nil {
		t.Fatal("A pem block was not found in the encoded key.")
	}

	if len(rest) != 0 {
		t.Errorf("The encoded key was not fully pem decoded. Remaining buffer len %v.", len(rest))
	}

	pub, err := x509.ParsePKIXPublicKey(pBlock.Bytes)

	if err != nil {
		t.Errorf("Parsing the key as DER public key returned the error %v.", err)
	}

	if pub == nil {
		t.Fatal("The key could not be parsed as a DER public key.")
	}

	if rpk, ok := pub.(*rsa.PublicKey); ok {
		rn := rpk.N.Int64()
		en := rsaKey.N.Int64()
		if en != rn {
			t.Error("Expected N", en, "but got", rn)
		}
		if rpk.E != rsaKey.E {
			t.Error("Expected E", rsaKey.E, "but got", rpk.E)
		}
	} else {
		t.Errorf("Expected public key type '*rsa.PublicKey' but got %T.", pub)
	}
}
