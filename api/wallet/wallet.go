package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

const (
	scryptR     = 8
	scryptDKLen = 32

	InternalServerError          = 10001 // 内部服务器错误
	RequestParameterError        = 10002 // 请求参数错误
	CreateMnemonicError          = 10003 // 创建助记词出错
	EncryptMnemonicError         = 10004 // 加密助记词出错
	MarshalEnMnemoicError        = 10005 // 对助记词密文信息序列化时出错
	UnMarshalEnMnemoicError      = 10006 // 解析助记词密文出错
	DecryptEnMnemoicError        = 10007 // 解密助记词错误
	CreateWalletFromMnemoicError = 10008 // 通过助记词创建钱包时出错
	DeriveWalletAccountError     = 10009 // 获取钱包用户账户时出错
	GetWalletPrivateError        = 10010 // 获取钱包用户私钥信息时错误
	GetWalletPublicError         = 10011 // 获取钱包用户公钥信息时错误
	HexPublicStringError         = 10012 // 公钥信息转十六进制字符串错误
	UnMarshalPublicBytesError    = 10013 // 反序列化二进制公钥信息错误
	CreateEntropyError			 = 10014 // 创建熵值出错
	EnMnemonicSignError			 = 10015 // 使用加密助记词签名时错误
	HexStringToBytesError		 = 10016 // 十六进制字符串转bytes错误
)

var (
	lastError = 0 // last error default 0
)

type Bytes = []byte

//字符串的编码为utf-8

//返回最后一次的错误
func GetLastError() int {
	return lastError
}

//16进制的字符串，转换为byte数据
func HexStringToBytes(str string) Bytes {
	lastError = 0
	bs, err := hex.DecodeString(str)
	if err != nil {
		lastError = HexStringToBytesError
		bs = nil
	}
	return bs
}

//byte数组转换为16进制字符串
func BytesToHexString(bytes Bytes) string {
	return hex.EncodeToString(bytes)
}

func Hash256(data Bytes) Bytes {
	hash := sha256.New()
	hash.Write(data)

	return hash.Sum(nil)
}

//返回值为签名后的交易
func SignTX(enMnemonic Bytes, password Bytes, tx string) string {

	lastError = 0

	mnemonic := DecodeMnemonic(enMnemonic, password)
	if mnemonic == nil {
		lastError = DecryptEnMnemoicError
		return ""
	}

	wallet, err := NewFromMnemonic(string(mnemonic))
	if err != nil {
		lastError = CreateWalletFromMnemoicError
		return ""
	}

	path := MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		lastError = DeriveWalletAccountError
		return ""
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		lastError = GetWalletPrivateError
		return ""
	}

	signBytes, err := crypto.Sign(Hash256([]byte(tx)), privateKey)
	if err != nil {
		lastError = EnMnemonicSignError
		return ""
	}

	return string(signBytes)
}

//Sign string
//privatekey
func Sign(privateKey Bytes, hash256 Bytes) (signedHash Bytes) {
	lastError = 0

	signedHash = nil
	key, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return
	}

	sig, err := crypto.Sign(hash256, key)
	if err == nil {
		signedHash = sig
	}
	return
}

//验证签名
//hex string
//pubkey without '0x'
func VerifySignature(publicKey Bytes, digestHash Bytes, signature Bytes) bool {
	return crypto.VerifySignature(publicKey, digestHash, signature[:len(signature)-1])
}

//生成助记词, 返回值为utf-8的字符串的byte数组， 使用byte数组比使用string更安全，以空格间隔
//Mnemonic, sample: picnic excite chat garden favorite exact maximum hire toddler asthma stove ramp vicious crazy employ cactus insane host
func GenerateMnemonic() Bytes {
	lastError = 0

	entropy, err := bip39.NewEntropy(192)
	if err != nil {
		lastError = CreateEntropyError
		return nil
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		lastError = CreateMnemonicError
		return nil
	}

	return Bytes(mnemonic)
}

//加密助记词
func EncodeMnemonic(mnemonic Bytes, password Bytes) Bytes {

	lastError = 0

	cryptoJson, err := keystore.EncryptDataV3(mnemonic, password, keystore.LightScryptN, keystore.StandardScryptP)
	if err != nil {
		lastError = EncryptMnemonicError
		return nil
	}

	cryptoRes := make(map[string]interface{}, 4)

	cryptoRes["iv"] = cryptoJson.CipherParams.IV
	cryptoRes["mac"] = cryptoJson.MAC
	cryptoRes["salt"] = cryptoJson.KDFParams["salt"]
	cryptoRes["cipherText"] = cryptoJson.CipherText

	bytes, err := json.Marshal(cryptoRes)
	if err != nil {
		lastError = MarshalEnMnemoicError
	}

	return bytes
}

//解密助记词
func DecodeMnemonic(enMnemonic Bytes, password Bytes) Bytes {

	lastError = 0

	keyMap := make(map[string]string, 4)
	err := json.Unmarshal(enMnemonic, &keyMap)
	if err != nil {
		lastError = UnMarshalEnMnemoicError
		return nil
	}
	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = keystore.LightScryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = keystore.StandardScryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = keyMap["salt"]

	cryptoStruct := keystore.CryptoJSON{
		Cipher:     "aes-128-ctr",
		CipherText: keyMap["cipherText"],
		KDF:        "scrypt",
		KDFParams:  scryptParamsJSON,
		MAC:        keyMap["mac"],
	}
	cryptoStruct.CipherParams.IV = keyMap["iv"]

	bytes, err := keystore.DecryptDataV3(cryptoStruct, string(password))
	if err != nil {
		lastError = DecryptEnMnemoicError
		return nil
	}

	return bytes
}

//通过助记词返回主私钥、公钥、地址， utf-8编码的byte数组
// 返回值为： 私钥：公钥：地址， 中间使用字符“:”间隔
func FromMnemonic(mnemonic Bytes) Bytes {

	lastError = 0

	wallet, err := NewFromMnemonic(string(mnemonic))
	if err != nil {
		lastError = CreateWalletFromMnemoicError
		return nil
	}

	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	account, err := wallet.Derive(path, false)
	if err != nil {
		lastError = DeriveWalletAccountError
		return nil
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		lastError = GetWalletPrivateError
		return nil
	}

	publicKey, err := wallet.PublicKeyHex(account)
	if err != nil {
		lastError = GetWalletPublicError
		return nil
	}

	return []byte(privateKey + ":" + publicKey + ":" + account.Address.Hex())
}

//通过加密的助记词返回主私钥：公钥：地址
// 返回值为： 私钥：公钥：地址， 中间使用字符“:”间隔
func FromEncodeMnemonic(enMnemonic Bytes, password Bytes) Bytes {

	// 解密助记词
	bytes := DecodeMnemonic(enMnemonic, password)

	return FromMnemonic(bytes)
}

// 通过公钥解析地址
func AddressFromPublickey(pubkey Bytes) Bytes {
	lastError = 0

	publicKeyStr := "04" + string(pubkey)
	bytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		lastError = HexPublicStringError
		return nil
	}

	pulicKey, err := crypto.UnmarshalPubkey(bytes)
	if err != nil {
		lastError = UnMarshalPublicBytesError
		return nil
	}

	return []byte(crypto.PubkeyToAddress(*pulicKey).String())
}
