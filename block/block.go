package block

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io"

	"golang.org/x/crypto/sha3"
)
import "crypto/cipher"
import "crypto/rsa"
import "crypto/aes"
import "crypto/rand"

type BlockData struct {
	encryptedKey     string   // RSA encrypted AES key
	encryptedMessage string   // AES encrypted message
	salt             [8]byte  // Random salt
	parent           [64]byte // Hash of parent block
}

type Block struct {
	// Block Data
	data BlockData

	// Non-hashed data
	ID     [64]byte // Hash of block data
	pepper [8]byte  // Random non-hashed salt
}

// Returns n random bytes
func RandomBytes(n int) []byte {
	out := make([]byte, n)
	rand.Read(out)
	return out
}

// Selects a block parent based on the encrypted message
func selectParentHash(encryptedMessage string) [64]byte {
	// TODO: Connect this to the blockpool
	var out [64]byte
	copy(out[:], sha3.New512().Sum(RandomBytes(32))[:64])
	return out
}

func CreateBlockData(message string, key *rsa.PublicKey) BlockData {
	var out BlockData

	// Block salt
	copy(out.salt[:], RandomBytes(8)[:8])

	// Message encryption
	// Random AES256 key
	AESkey := RandomBytes(32)
	// Block cipher for that key
	AESCipher, e := aes.NewCipher(AESkey)
	if e != nil {
		panic(e)
	}
	// encrypted message bytes
	cipherBytes := make([]byte, aes.BlockSize+len(message))
	// Initialization Vector
	// Delivered with ciphertext as it is necessary for decryption...
	// But it doesn't have to be private to be secure
	iv := cipherBytes[:aes.BlockSize]
	// If error reading IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err) // panic!!
	}
	// Stream cipher
	stream := cipher.NewCTR(AESCipher, iv)
	// Plaintext bytes
	plaintext := []byte(message)
	// Encryption
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], plaintext)
	// Convert to base64 and place in block
	out.encryptedMessage = base64.URLEncoding.EncodeToString(cipherBytes)

	// AES key encryption
	cipheredKey, e := rsa.EncryptOAEP(sha3.New512(), rand.Reader, key, AESkey, nil)
	// Panic on error
	if e != nil {
		panic(e)
	}
	// Convert to base64 and place in block
	out.encryptedKey = base64.URLEncoding.EncodeToString(cipheredKey)

	// Select blockparent using blockpool
	out.parent = selectParentHash(out.encryptedMessage)

	// Done.
	return out
}

// BlockData -> string
func StringifyBlockData(data BlockData) string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(data)
	raw := buf.Bytes()
	return string(raw[:buf.Len()])
}

func DestringifyBlockData(data string) BlockData {
	var out BlockData

	var buf bytes.Buffer
	buf.WriteString(data)
	decoder := gob.NewDecoder(&buf)
	decoder.Decode(out)

	return out
}

func CreateBlock(message string, key *rsa.PublicKey) Block {
	var out Block

	// Block data
	out.data = CreateBlockData(message, key)

	// Block ID
	dataString := StringifyBlockData(out.data)
	copy(out.ID[:], sha3.New512().Sum([]byte(dataString))[:64])

	// Block pepper
	copy(out.pepper[:], RandomBytes(8)[:8])

	return out
}

func StringifyBlock(block Block) string {
	// TODO Implement this

	return "Block serialization not yet ready. This is not your bug."
}

func DestringifyBlock(block string) Block {
	var out Block

	// TODO implement this
	return out
}

// Call on main startup
func Initialize() {
	gob.Register(BlockData{})
	gob.Register(Block{})
}
