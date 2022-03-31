package models

import (
	bytes "bytes"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	types "github.com/gogo/protobuf/types"
	models "gitlab.com/alphaticks/xchanger/models"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ProtocolAsset struct {
	Address   []byte            `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	MaxSupply *types.BytesValue `protobuf:"bytes,2,opt,name=max_supply,json=maxSupply,proto3" json:"max_supply,omitempty"`
	Protocol  *models.Protocol  `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Asset     *models.Asset     `protobuf:"bytes,4,opt,name=asset,proto3" json:"asset,omitempty"`
	Chain     *models.Chain     `protobuf:"bytes,5,opt,name=chain,proto3" json:"chain,omitempty"`
	Meta      map[string]string `protobuf:"bytes,6,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *ProtocolAsset) Reset()      { *m = ProtocolAsset{} }
func (*ProtocolAsset) ProtoMessage() {}
func (*ProtocolAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7bd1bc6b244e966, []int{0}
}
func (m *ProtocolAsset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAsset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAsset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAsset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAsset.Merge(m, src)
}
func (m *ProtocolAsset) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAsset) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAsset.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAsset proto.InternalMessageInfo

func (m *ProtocolAsset) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *ProtocolAsset) GetMaxSupply() *types.BytesValue {
	if m != nil {
		return m.MaxSupply
	}
	return nil
}

func (m *ProtocolAsset) GetProtocol() *models.Protocol {
	if m != nil {
		return m.Protocol
	}
	return nil
}

func (m *ProtocolAsset) GetAsset() *models.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *ProtocolAsset) GetChain() *models.Chain {
	if m != nil {
		return m.Chain
	}
	return nil
}

func (m *ProtocolAsset) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtocolAsset)(nil), "models.ProtocolAsset")
	proto.RegisterMapType((map[string]string)(nil), "models.ProtocolAsset.MetaEntry")
}

func init() { proto.RegisterFile("asset_data.proto", fileDescriptor_f7bd1bc6b244e966) }

var fileDescriptor_f7bd1bc6b244e966 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xbd, 0x6a, 0xe3, 0x40,
	0x10, 0x80, 0xb5, 0xfe, 0xbb, 0xd3, 0xfa, 0x0c, 0x66, 0xb9, 0x42, 0xf8, 0x60, 0xcf, 0xdc, 0x35,
	0x2e, 0xc2, 0x0a, 0xec, 0x40, 0x82, 0xbb, 0x38, 0xa4, 0x0c, 0x84, 0x0d, 0xa4, 0x35, 0x63, 0x69,
	0x23, 0x1b, 0x4b, 0x96, 0xd0, 0xae, 0x13, 0xab, 0xcb, 0x23, 0xe4, 0x31, 0xf2, 0x28, 0x81, 0x34,
	0x2e, 0x5d, 0xc6, 0x72, 0x93, 0xd2, 0x8f, 0x10, 0xb4, 0x6b, 0x19, 0x92, 0x6e, 0x67, 0xe6, 0x9b,
	0xfd, 0x66, 0x06, 0xb7, 0x41, 0x4a, 0xa1, 0xc6, 0x3e, 0x28, 0x60, 0x49, 0x1a, 0xab, 0x98, 0x34,
	0xa2, 0xd8, 0x17, 0xa1, 0xec, 0xf4, 0x83, 0x99, 0x0a, 0x61, 0xc2, 0xbc, 0x38, 0x72, 0x21, 0x4c,
	0xa6, 0xa0, 0x66, 0xde, 0x5c, 0xba, 0x2b, 0x6f, 0x0a, 0x8b, 0x40, 0xa4, 0xae, 0xc1, 0x5c, 0xdd,
	0x24, 0x4d, 0x6f, 0x87, 0x06, 0x71, 0x1c, 0x84, 0xc2, 0x24, 0x27, 0xcb, 0x7b, 0xf7, 0x31, 0x85,
	0x24, 0x11, 0xe9, 0xa1, 0xfe, 0xef, 0xad, 0x82, 0x5b, 0x37, 0xc5, 0xcb, 0x8b, 0xc3, 0x8b, 0x42,
	0x4c, 0x1c, 0xfc, 0x03, 0x7c, 0x3f, 0x15, 0x52, 0x3a, 0xa8, 0x8b, 0x7a, 0xbf, 0x78, 0x19, 0x92,
	0x21, 0xc6, 0x11, 0xac, 0xc6, 0x72, 0x99, 0x24, 0x61, 0xe6, 0x54, 0xba, 0xa8, 0xd7, 0xec, 0xff,
	0x61, 0x46, 0xc0, 0x4a, 0x01, 0x1b, 0x65, 0x4a, 0xc8, 0x3b, 0x08, 0x97, 0x82, 0xdb, 0x11, 0xac,
	0x6e, 0x35, 0x4d, 0x4e, 0xf0, 0xcf, 0xe4, 0xa0, 0x71, 0xaa, 0xba, 0xb3, 0xcd, 0xcc, 0xbc, 0xac,
	0xd4, 0xf3, 0x23, 0x41, 0xfe, 0xe3, 0xba, 0xbe, 0x82, 0x53, 0xd3, 0x68, 0xab, 0x44, 0xf5, 0x84,
	0xdc, 0xd4, 0x0a, 0xc8, 0x9b, 0xc2, 0x6c, 0xe1, 0xd4, 0xbf, 0x42, 0x97, 0x45, 0x92, 0x9b, 0x1a,
	0x19, 0xe0, 0x5a, 0x24, 0x14, 0x38, 0x8d, 0x6e, 0xb5, 0xd7, 0xec, 0xff, 0xfd, 0xee, 0xd4, 0x1f,
	0xb2, 0x6b, 0xa1, 0xe0, 0x6a, 0xa1, 0xd2, 0x8c, 0x6b, 0xb8, 0x73, 0x86, 0xed, 0x63, 0x8a, 0xb4,
	0x71, 0x75, 0x2e, 0x32, 0x7d, 0x0b, 0x9b, 0x17, 0x4f, 0xf2, 0x1b, 0xd7, 0x1f, 0x8a, 0xfd, 0xf4,
	0x09, 0x6c, 0x6e, 0x82, 0x61, 0xe5, 0x1c, 0x8d, 0x4e, 0xd7, 0x5b, 0x6a, 0x6d, 0xb6, 0xd4, 0xda,
	0x6f, 0x29, 0x7a, 0xca, 0x29, 0x7a, 0xc9, 0x29, 0x7a, 0xcd, 0x29, 0x5a, 0xe7, 0x14, 0xbd, 0xe7,
	0x14, 0x7d, 0xe4, 0xd4, 0xda, 0xe7, 0x14, 0x3d, 0xef, 0xa8, 0xb5, 0xde, 0x51, 0x6b, 0xb3, 0xa3,
	0xd6, 0xa4, 0xa1, 0xf7, 0x1e, 0x7c, 0x06, 0x00, 0x00, 0xff, 0xff, 0x78, 0xf4, 0xc6, 0x59, 0xfa,
	0x01, 0x00, 0x00,
}

func (this *ProtocolAsset) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAsset)
	if !ok {
		that2, ok := that.(ProtocolAsset)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Address, that1.Address) {
		return false
	}
	if !this.MaxSupply.Equal(that1.MaxSupply) {
		return false
	}
	if !this.Protocol.Equal(that1.Protocol) {
		return false
	}
	if !this.Asset.Equal(that1.Asset) {
		return false
	}
	if !this.Chain.Equal(that1.Chain) {
		return false
	}
	if len(this.Meta) != len(that1.Meta) {
		return false
	}
	for i := range this.Meta {
		if this.Meta[i] != that1.Meta[i] {
			return false
		}
	}
	return true
}
func (this *ProtocolAsset) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&models.ProtocolAsset{")
	s = append(s, "Address: "+fmt.Sprintf("%#v", this.Address)+",\n")
	if this.MaxSupply != nil {
		s = append(s, "MaxSupply: "+fmt.Sprintf("%#v", this.MaxSupply)+",\n")
	}
	if this.Protocol != nil {
		s = append(s, "Protocol: "+fmt.Sprintf("%#v", this.Protocol)+",\n")
	}
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	if this.Chain != nil {
		s = append(s, "Chain: "+fmt.Sprintf("%#v", this.Chain)+",\n")
	}
	keysForMeta := make([]string, 0, len(this.Meta))
	for k, _ := range this.Meta {
		keysForMeta = append(keysForMeta, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForMeta)
	mapStringForMeta := "map[string]string{"
	for _, k := range keysForMeta {
		mapStringForMeta += fmt.Sprintf("%#v: %#v,", k, this.Meta[k])
	}
	mapStringForMeta += "}"
	if this.Meta != nil {
		s = append(s, "Meta: "+mapStringForMeta+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringAssetData(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *ProtocolAsset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAsset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAsset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Meta) > 0 {
		for k := range m.Meta {
			v := m.Meta[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintAssetData(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintAssetData(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintAssetData(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x32
		}
	}
	if m.Chain != nil {
		{
			size, err := m.Chain.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAssetData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Asset != nil {
		{
			size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAssetData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Protocol != nil {
		{
			size, err := m.Protocol.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAssetData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.MaxSupply != nil {
		{
			size, err := m.MaxSupply.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAssetData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintAssetData(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAssetData(dAtA []byte, offset int, v uint64) int {
	offset -= sovAssetData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProtocolAsset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovAssetData(uint64(l))
	}
	if m.MaxSupply != nil {
		l = m.MaxSupply.Size()
		n += 1 + l + sovAssetData(uint64(l))
	}
	if m.Protocol != nil {
		l = m.Protocol.Size()
		n += 1 + l + sovAssetData(uint64(l))
	}
	if m.Asset != nil {
		l = m.Asset.Size()
		n += 1 + l + sovAssetData(uint64(l))
	}
	if m.Chain != nil {
		l = m.Chain.Size()
		n += 1 + l + sovAssetData(uint64(l))
	}
	if len(m.Meta) > 0 {
		for k, v := range m.Meta {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovAssetData(uint64(len(k))) + 1 + len(v) + sovAssetData(uint64(len(v)))
			n += mapEntrySize + 1 + sovAssetData(uint64(mapEntrySize))
		}
	}
	return n
}

func sovAssetData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAssetData(x uint64) (n int) {
	return sovAssetData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ProtocolAsset) String() string {
	if this == nil {
		return "nil"
	}
	keysForMeta := make([]string, 0, len(this.Meta))
	for k, _ := range this.Meta {
		keysForMeta = append(keysForMeta, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForMeta)
	mapStringForMeta := "map[string]string{"
	for _, k := range keysForMeta {
		mapStringForMeta += fmt.Sprintf("%v: %v,", k, this.Meta[k])
	}
	mapStringForMeta += "}"
	s := strings.Join([]string{`&ProtocolAsset{`,
		`Address:` + fmt.Sprintf("%v", this.Address) + `,`,
		`MaxSupply:` + strings.Replace(fmt.Sprintf("%v", this.MaxSupply), "BytesValue", "types.BytesValue", 1) + `,`,
		`Protocol:` + strings.Replace(fmt.Sprintf("%v", this.Protocol), "Protocol", "models.Protocol", 1) + `,`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models.Asset", 1) + `,`,
		`Chain:` + strings.Replace(fmt.Sprintf("%v", this.Chain), "Chain", "models.Chain", 1) + `,`,
		`Meta:` + mapStringForMeta + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAssetData(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ProtocolAsset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAssetData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProtocolAsset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAsset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxSupply", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaxSupply == nil {
				m.MaxSupply = &types.BytesValue{}
			}
			if err := m.MaxSupply.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Protocol", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Protocol == nil {
				m.Protocol = &models.Protocol{}
			}
			if err := m.Protocol.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Asset == nil {
				m.Asset = &models.Asset{}
			}
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Chain == nil {
				m.Chain = &models.Chain{}
			}
			if err := m.Chain.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAssetData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAssetData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAssetData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAssetData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthAssetData
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthAssetData
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAssetData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthAssetData
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthAssetData
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipAssetData(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthAssetData
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Meta[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAssetData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAssetData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAssetData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAssetData
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAssetData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAssetData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAssetData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAssetData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAssetData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAssetData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAssetData = fmt.Errorf("proto: unexpected end of group")
)
