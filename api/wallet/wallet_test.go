package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"reflect"
	"strings"
	"testing"
)

func TestGenerateMnemonic(t *testing.T) {
	m := GenerateMnemonic()
	// 获取的助记词结果类似下面字符
	// use fan yard throw cage flush drop hint empty toe oyster merge cube aisle casual maximum design adapt
	if len(m) < 1 {
		t.Error("GenerateMnemonic error")
	}
	strs := strings.Split(string(m), " ")
	if len(strs) != 18 {
		t.Error("length of GenerateMnemonic is not 18")
	}
	if !bip39.IsMnemonicValid(string(m)){
		t.Fatalf("test: is mnemonic valid : false")
	}
	fmt.Println(strs, len(strs))
}

func TestEncodeMnemonic(t *testing.T) {
	mnemonic := "rose rocket invest real refuse margin festival danger anger border idle brown"
	key := []byte("")

	plainText := EncodeMnemonic([]byte(mnemonic), key)
	if plainText == nil {
		lastError := GetLastError()
		fmt.Println("last error code :", lastError)
	}

	res := map[string]interface{}{}
	json.Unmarshal(plainText, &res)

	fmt.Println("result:", res)

}

func TestFromEncodeMnemonic(t *testing.T) {

	keyMap := make(map[string]string, 4)
	keyMap["iv"] = "afb4984ec1da4e5f27a04b99147b426a"
	keyMap["mac"] = "1b805e79f35d29d32dc8562be40ba5d310593bb16db7dc61fedb53e732d5c896"
	keyMap["salt"] = "6882600426051f40cfe2d21eafd8da6a32dd217f7e19b77c72f45190ce997d3b"
	keyMap["cipherText"] = "4fc1f9766f3278d763866f1f7f8399cce0e907748d1830d1dbd63ada414ccd17c7b5065ecc042de80964873a3dadaec3a167a60ba247f160ac90bdf227167a1ccf89eac7a4e7016a7c9ca4bbf6"

	jsonKey, _ := json.Marshal(keyMap)

	res := FromEncodeMnemonic(jsonKey,[]byte(""))
	if res == nil {
		lastError := GetLastError()
		fmt.Println("last error code:", lastError)
	}

	// private key : 68e043089ab98b4a9238d51a362deccf49993939c75c2444dce34ac64f207d0d
	// public key : e7a9526208844cbfa52cd1005f53590e0662f420ee87fd8922a3cf0b2e5c09fd8633c70ccf1a1c4da5e98dc5614a90a743dfab48e6e0ffb139aa4b0fff2b7448
	// address : 0x685ce4CbDd5c19b64CA008cB85b83947e5318EFA
	fmt.Println("AES解密完成：", string(res))
}

func TestDecodeMnemonic(t *testing.T) {
	keyMap := make(map[string]string, 4)
	keyMap["iv"] = "afb4984ec1da4e5f27a04b99147b426a"
	keyMap["mac"] = "1b805e79f35d29d32dc8562be40ba5d310593bb16db7dc61fedb53e732d5c896"
	keyMap["salt"] = "6882600426051f40cfe2d21eafd8da6a32dd217f7e19b77c72f45190ce997d3b"
	keyMap["cipherText"] = "4fc1f9766f3278d763866f1f7f8399cce0e907748d1830d1dbd63ada414ccd17c7b5065ecc042de80964873a3dadaec3a167a60ba247f160ac90bdf227167a1ccf89eac7a4e7016a7c9ca4bbf6"

	jsonKey, _ := json.Marshal(keyMap)

	res := DecodeMnemonic(jsonKey,[]byte(""))
	if res == nil {
		lastError := GetLastError()
		fmt.Println("last error code :", lastError)
	}

	// rose rocket invest real refuse margin festival danger anger border idle brown
	fmt.Println("AES解密完成：", string(res))
}

func TestFromMnemonic(t *testing.T) {
	tests := []struct{
		input string
		output struct{
			privateKey string
			publicKey string // 取前64位
			address string
		}
	}{
		{"", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "", publicKey: "", address: ""}},
		{"ffff ffddd asdf dfdf", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "", publicKey: "", address: ""}},
		{"gown distance below", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "", publicKey: "", address: ""}},
		{"rose rocket invest real refuse margin festival danger anger border idle brown", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "68e043089ab98b4a9238d51a362deccf49993939c75c2444dce34ac64f207d0d", publicKey: "e7a9526208844cbfa52cd1005f53590e0662f420ee87fd8922a3cf0b2e5c09fd", address: "0x685ce4CbDd5c19b64CA008cB85b83947e5318EFA"}},
		{"raw fashion exhibit lend soon actual science magnet miracle captain vanish vague pepper derive sure", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "bf8ef751e5da6b99dcb925277c44173dbcf4ca334b4ea0dbdac21dd4da69dd55", publicKey: "7a987e1033058e5779eaf8b6efcbd3155da5742f1990516b0555154b934f4198", address: "0xAc6aB134eEa5C938c19CB72Ec5E57413CD62e430"}},
		{"boost foam cannon message adult faith column tennis boost target aspect juice decrease fabric ignore brush collect snow", struct {
			privateKey string
			publicKey  string
			address    string
		}{privateKey: "4ce3e6606b3d36d94357fb97880cfa6d81d93a1988507113cd9f35f2481a22ae", publicKey: "f2b52c033387b1c2fe29e950b62de65f33dd8bd431932a25f0593176a32a3282", address: "0xFCcc8676C1991d3d1Ae841670aD143D8eeDD55E2"}},
	}

	for i, test := range tests {
		result := FromMnemonic([]byte(test.input))
		fmt.Println(i, "result: " + string(result))

		if result == nil {
			if test.output.privateKey != "" || test.output.publicKey != "" || test.output.address != "" {
				t.Fatalf("test: %d Expected and returned results are not equal", i)
			}
			continue
		}
		if sliceResult := strings.Split(string(result), ":"); !reflect.DeepEqual(struct {
			privateKey string
			publicKey string
			address string
		}{sliceResult[0], (sliceResult[1])[:64],sliceResult[2]}, test.output) {
			t.Fatalf("test: %d Expected and returned results are not equal", i)
		}
	}
}


func TestAddressFromPublickey(t *testing.T) {

	tests := []struct{
		input string
		output []byte
	}{
		{"", nil},
		{"fffffffffffffffffffffffffffffffffffff", nil},
		{"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", nil},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", nil},
		{"35c98a5f7f1aaa1220ea0dec2e62a5045a622e2d060735c0b1fe2380e1c2ef26d97f22419c115416325b1d6d46a9d002ff6ee40b795264a71b9388a96729eb28",
			[]byte("0xa584cABeEF8C69e3579B03F1699B96BeA63978f7")},
		{"e7a9526208844cbfa52cd1005f53590e0662f420ee87fd8922a3cf0b2e5c09fd8633c70ccf1a1c4da5e98dc5614a90a743dfab48e6e0ffb139aa4b0fff2b7448",
			[]byte("0x685ce4CbDd5c19b64CA008cB85b83947e5318EFA")},
	}

	for i, test := range tests {
		if result := AddressFromPublickey([]byte(test.input)); !bytes.Equal(result, test.output) {
			t.Fatalf("test %d: address from public : result[%s], want[%s]", i, string(result), string(test.output))
		}
	}

}
