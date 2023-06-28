package aes256

import (
	"errors"

	"github.com/wedkarz02/aes256-go/aes256/consts"
	"github.com/wedkarz02/aes256-go/aes256/key"
	"github.com/wedkarz02/aes256-go/aes256/sbox"
)

type AES256 struct {
	Key         [consts.KEY_SIZE]byte
	ExpandedKey *key.ExpandedKey
}

// NewAES256 initializes new AES cipher
// and calculates round keys.
func NewAES256(k []byte) (*AES256, error) {
	if len(k) != consts.KEY_SIZE {
		return nil, errors.New("invalid key size")
	}

	a := AES256{Key: [32]byte(k)}

	var err error
	a.ExpandedKey, err = a.NewExpKey()

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// NewSBox returns a byte slice representation of
// the AES substitution look up table.
//
// https://en.wikipedia.org/wiki/Rijndael_S-box
func NewSBox() *sbox.SBOX {
	return sbox.InitSBOX()
}

// NewInvSBox returns a byte slice representation of
// the AES inverse substitution look up table.
//
// https://en.wikipedia.org/wiki/Rijndael_S-box
func NewInvSBox(sb *sbox.SBOX) *sbox.SBOX {
	return sbox.InitInvSBOX(sb)
}

// NewExpKey returns a key expanded by a key
// schedule to a slice of unique round keys.
//
// https://en.wikipedia.org/wiki/AES_key_schedule
//
// https://www.samiam.org/key-schedule.html
func (a *AES256) NewExpKey() (*key.ExpandedKey, error) {
	xKey, err := key.ExpandKey(a.Key[:])

	if err != nil {
		return nil, err
	}

	return xKey, nil
}

// SubBytes returns a state with every
// byte replaced with it's corresponding
// byte from the sbox.
//
// https://pl.wikipedia.org/wiki/Advanced_Encryption_Standard
func (a *AES256) SubBytes(state []byte) ([]byte, error) {
	if len(state) != consts.BLOCK_SIZE {
		return nil, errors.New("state size not matching the block size")
	}

	var subState []byte

	sbox := NewSBox()
	for i := range state {
		subState = append(subState, sbox[state[i]])
	}

	return subState, nil
}

// InvSubBytes undoes the SubBytes operation
// allowing decryption.
//
// https://pl.wikipedia.org/wiki/Advanced_Encryption_Standard
func (a *AES256) InvSubBytes(state []byte) ([]byte, error) {
	if len(state) != consts.BLOCK_SIZE {
		return nil, errors.New("state size not matching the block size")
	}

	var invSubState []byte

	invsbox := NewInvSBox(NewSBox())
	for i := range state {
		invSubState = append(invSubState, invsbox[state[i]])
	}

	return invSubState, nil
}

// ShiftRows returns a state where the last three
// rows has been transposed in an AES specific way.
//
// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
func (a *AES256) ShiftRows(state []byte) ([]byte, error) {
	if len(state) != consts.BLOCK_SIZE {
		return nil, errors.New("state size not matching the block size")
	}

	var shiftedState []byte
	copy(shiftedState, state)

	for i := 1; i < 4; i++ {
		j := i

		shiftedState[i+(4*0)] = state[i+4*((j+0)%4)]
		shiftedState[i+(4*1)] = state[i+4*((j+1)%4)]
		shiftedState[i+(4*2)] = state[i+4*((j+2)%4)]
		shiftedState[i+(4*3)] = state[i+4*((j+4)%4)]
	}

	return shiftedState, nil
}

// InvShiftRows undoes the ShiftRows operation
// allowing decryption.
//
// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
func (a *AES256) InvShiftRows(state []byte) ([]byte, error) {
	if len(state) != consts.BLOCK_SIZE {
		return nil, errors.New("state size not matching the block size")
	}

	var invShiftedState []byte
	copy(invShiftedState, state)

	for i := 1; i < 4; i++ {
		j := 4 - i

		invShiftedState[i+(4*0)] = state[i+4*((j+0)%4)]
		invShiftedState[i+(4*1)] = state[i+4*((j+1)%4)]
		invShiftedState[i+(4*2)] = state[i+4*((j+2)%4)]
		invShiftedState[i+(4*3)] = state[i+4*((j+4)%4)]
	}

	return invShiftedState, nil
}
