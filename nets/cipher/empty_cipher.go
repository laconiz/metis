package cipher

type EmptyCipher struct {
}

func (cipher *EmptyCipher) Encode(raw []byte) ([]byte, error) {
	return raw, nil
}

func (cipher *EmptyCipher) Decode(stream []byte) ([]byte, error) {
	return stream, nil
}

type EmptyMaker struct {
	cipher EmptyCipher
}

func (maker EmptyMaker) New() Cipher {
	return &maker.cipher
}
