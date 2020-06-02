package packet

type Packet struct {
	Meta   *Meta
	Msg    interface{}
	Raw    []byte
	Stream []byte
}
