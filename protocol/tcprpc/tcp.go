//解TCP 包

package tcprpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/tinyhole/ap/protocol"
	pbAP "github.com/tinyhole/ap/protocol/pack"
)

var (
	ErrTotalLengthNotEnough = errors.New("message total length not enough")
	ErrHeadIncorrect        = errors.New("message head incorrect")
	ErrPackTypeIncorrect    = errors.New("package type incorrect")
)

const (
	TotalLenBytesSize = 4
	HeadLenBytesSize  = 4
)

type tcpCodec struct{}

func NewTcpCodec() protocol.Codec {
	return &tcpCodec{}
}

func (t tcpCodec) Marshal(v interface{}) ([]byte, error) {
	var (
		byteArray []byte
		headBuf   []byte

		totalLen int
		headLen  int

		headLenBuf  = make([]byte, 0, HeadLenBytesSize)
		totalLenBuf = make([]byte, 0, TotalLenBytesSize)

		err  error
		pack *pbAP.ApPackage
		ok   bool
	)
	pack, ok = v.(*pbAP.ApPackage)
	if !ok {
		return nil, ErrPackTypeIncorrect
	}

	//编码头
	headBuf, err = proto.Marshal(pack.Header)
	if err != nil {
		return nil, err
	}

	//计算长度
	headLen = len(headBuf)
	totalLen = len(pack.Body) + headLen
	byteArray = make([]byte, 0, totalLen+TotalLenBytesSize+HeadLenBytesSize)

	//放入数据
	binary.BigEndian.PutUint32(headLenBuf, uint32(headLen))
	binary.BigEndian.PutUint32(totalLenBuf, uint32(totalLen))
	byteArray = append(byteArray[:0], totalLenBuf...)
	byteArray = append(byteArray[:TotalLenBytesSize], headLenBuf...)
	byteArray = append(byteArray[:totalLen], pack.Body...)

	return byteArray, nil
}

//Unmarshal 反解TCP 协议消息用
func (t *tcpCodec) Unmarshal(data []byte, v interface{}) error {
	var (
		headLen  uint32
		totalLen uint32
		pack     *pbAP.ApPackage
		ok       bool
	)

	pack, ok = v.(*pbAP.ApPackage)
	if !ok {
		return ErrPackTypeIncorrect
	}

	//读取数据长度
	buf := bytes.NewBuffer(data)
	TotalLenBuf := buf.Next(TotalLenBytesSize)
	headLenBuf := buf.Next(HeadLenBytesSize)
	totalLen = binary.BigEndian.Uint32(TotalLenBuf)
	if len(data) != int(totalLen) {
		return ErrTotalLengthNotEnough
	}
	headLen = binary.BigEndian.Uint32(headLenBuf)

	//解头部
	headBuf := buf.Next(int(headLen))
	pbHead := &pbAP.Header{}
	proto.Unmarshal(headBuf, pbHead)
	pack.Header = pbHead

	pack.Body = buf.Next(int(totalLen - headLen))
	return nil
}
