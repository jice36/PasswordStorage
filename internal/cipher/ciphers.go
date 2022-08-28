package cipher

import (
	"crypto/aes"
	"crypto/cipher"
)

type Cipher struct {
	cipher cipher.Block
}

func NewCipher(key []byte) (*Cipher, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Cipher{cipher: cipher}, nil
}

func (c *Cipher) EncryptModeECB(src []byte) []byte {
	dst := make([]byte, 0)
	if len(src)%aes.BlockSize == 0 {
		count := numberBlocks(len(src))
		for i := 0; i < count; i++ {
			dstForCipher := make([]byte, aes.BlockSize)
			c.cipher.Encrypt(dstForCipher, src[i*16:(i+1)*16])
			dst = append(dst, dstForCipher...)
		}
	} else if len(src)%aes.BlockSize != 0 {
		count := numberBlocks(len(src))
		for i := 0; i < count; i++ {
			dstForCipher := make([]byte, aes.BlockSize)
			c.cipher.Encrypt(dstForCipher, src[i*16:(i+1)*16])
			dst = append(dst, dstForCipher...)
		}
		dstForCipher := make([]byte, aes.BlockSize)
		c.cipher.Encrypt(dstForCipher, zeroPadding(src[count*16:], aes.BlockSize))
		dst = append(dst, dstForCipher...)
	}
	return dst
}

func (c *Cipher) DecryptModeECB(src []byte) []byte {
	dst := make([]byte, 0)
	count := numberBlocks(len(src))
	for i := 0; i < count-1; i++ {
		dstForCipher := make([]byte, aes.BlockSize)
		c.cipher.Decrypt(dstForCipher, src[i*16:(i+1)*16])
		dst = append(dst, dstForCipher...)
	}
	dstForCipher := make([]byte, aes.BlockSize)
	c.cipher.Decrypt(dstForCipher, src[(count-1)*16:])
	dst = append(dst, dropPadding(dstForCipher)...)
	return dst
}

func numberBlocks(lengthBlock int) int {
	return (lengthBlock) / 16
}

func zeroPadding(block []byte, blockSize int) []byte {
	pad := blockSize - len(block)
	if pad == 0 {
		return block
	}

	blockWithPad := make([]byte, len(block)+pad)
	copy(blockWithPad, block)
	blockWithPad[len(blockWithPad)-1] = byte(pad)
	return blockWithPad
}

func dropPadding(block []byte) []byte {
	for i := len(block) - 2; i >= 0; i-- {
		if block[i] == 0 {
			break
		}
		return block
	}
	lengthPad := int(block[len(block)-1])
	return block[:len(block)-lengthPad]
}
