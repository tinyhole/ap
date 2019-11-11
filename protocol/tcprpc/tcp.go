//解TCP 包

package tcprpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/tinyhole/ap/protocol"
	pbAP "github.com/tinyhole/ap/protocol/pack"
	"sync"
)

var (
	ErrFlagLengthNotEnough  = errors.New("message flag length not enough")
	ErrTotalLengthNotEnough = errors.New("message total length not enough")
	ErrPackTypeIncorrect    = errors.New("package type incorrect")
)

var (
	codecPool = sync.Pool{
		New: NewTcpCodec,
	}
)

const (
	TotalLenBytesSize = 4
	HeadLenBytesSize  = 4
)

type TcpCodec struct {
	TotalLen int32
}

func GetCodec() protocol.Codec {
	codec := codecPool.Get()
	if codec == nil {
		return NewTcpCodec().(protocol.Codec)
	}

	return codec.(protocol.Codec)
}

func PutCodec(codec protocol.Codec) {
	codecPool.Put(codec)
}

func NewTcpCodec() interface{} {
	return &TcpCodec{}
}

func (t TcpCodec) Marshal(v interface{}) ([]byte, error) {
	var (
		byteArray []byte
		headBuf   []byte

		totalLen int
		headLen  int

		headLenBuf  = make([]byte, HeadLenBytesSize, HeadLenBytesSize)
		totalLenBuf = make([]byte, TotalLenBytesSize, TotalLenBytesSize)

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

	fmt.Printf("head len [%d]", headLen)
	//放入数据
	binary.BigEndian.PutUint32(headLenBuf, uint32(headLen))
	binary.BigEndian.PutUint32(totalLenBuf, uint32(totalLen))
	byteArray = append(byteArray[:0], totalLenBuf...) //写入总长度
	byteArray = append(byteArray[:TotalLenBytesSize], headLenBuf...) //写入头长度
	byteArray = append(byteArray[:TotalLenBytesSize+HeadLenBytesSize], headBuf...) //写入头
	byteArray = append(byteArray[:TotalLenBytesSize+HeadLenBytesSize+headLen], pack.Body...) //写入body
	return byteArray, nil
}

//Unmarshal 反解TCP 协议消息用
func (t *TcpCodec) Unmarshal(data []byte, v interface{}) error {
	var (
		headLen  uint32
		totalLen uint32
		pack     *pbAP.ApPackage
		ok       bool
		err      error
	)

	pack, ok = v.(*pbAP.ApPackage)
	if !ok {
		return ErrPackTypeIncorrect
	}

	if len(data) < TotalLenBytesSize+HeadLenBytesSize {
		return ErrFlagLengthNotEnough
	}

	//读取数据长度
	buf := bytes.NewBuffer(data)
	TotalLenBuf := buf.Next(TotalLenBytesSize)
	headLenBuf := buf.Next(HeadLenBytesSize)
	totalLen = binary.BigEndian.Uint32(TotalLenBuf)
	t.TotalLen = int32(totalLen)
	headLen = binary.BigEndian.Uint32(headLenBuf)
	if buf.Len() < int(totalLen) {
		return ErrTotalLengthNotEnough
	}
	//解头部
	headBuf := buf.Next(int(headLen))
	pbHead := &pbAP.Header{}
	err = proto.Unmarshal(headBuf, pbHead)
	if err != nil {
		return err
	}
	pack.Header = pbHead
	pack.Body = buf.Next(int(totalLen - headLen))

	return nil
}
