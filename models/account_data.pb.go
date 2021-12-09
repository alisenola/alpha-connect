package models

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type Account struct {
	Name             string                   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Exchange         *models.Exchange         `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	ApiCredentials   *models.APICredentials   `protobuf:"bytes,3,opt,name=api_credentials,json=apiCredentials,proto3" json:"api_credentials,omitempty"`
	StarkCredentials *models.STARKCredentials `protobuf:"bytes,4,opt,name=stark_credentials,json=starkCredentials,proto3" json:"stark_credentials,omitempty"`
	EcdsaCredentials *models.ECDSACredentials `protobuf:"bytes,5,opt,name=ecdsa_credentials,json=ecdsaCredentials,proto3" json:"ecdsa_credentials,omitempty"`
}

func (m *Account) Reset()      { *m = Account{} }
func (*Account) ProtoMessage() {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e413886430e50, []int{0}
}
func (m *Account) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Account.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return m.Size()
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Account) GetExchange() *models.Exchange {
	if m != nil {
		return m.Exchange
	}
	return nil
}

func (m *Account) GetApiCredentials() *models.APICredentials {
	if m != nil {
		return m.ApiCredentials
	}
	return nil
}

func (m *Account) GetStarkCredentials() *models.STARKCredentials {
	if m != nil {
		return m.StarkCredentials
	}
	return nil
}

func (m *Account) GetEcdsaCredentials() *models.ECDSACredentials {
	if m != nil {
		return m.EcdsaCredentials
	}
	return nil
}

func init() {
	proto.RegisterType((*Account)(nil), "models.Account")
}

func init() { proto.RegisterFile("account_data.proto", fileDescriptor_941e413886430e50) }

var fileDescriptor_941e413886430e50 = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0xd0, 0x31, 0x4e, 0xf3, 0x30,
	0x14, 0xc0, 0x71, 0xbb, 0x5f, 0xbf, 0x02, 0x46, 0x82, 0xe2, 0x01, 0x45, 0x0c, 0x4f, 0x15, 0x53,
	0x07, 0x94, 0x48, 0x85, 0x1d, 0x85, 0xd2, 0x01, 0xb1, 0xa0, 0x94, 0xbd, 0x7a, 0x75, 0xac, 0x36,
	0x6a, 0x9a, 0x44, 0xb1, 0x91, 0x18, 0x39, 0x02, 0x33, 0x27, 0xe0, 0x28, 0x8c, 0x1d, 0x3b, 0x12,
	0x67, 0x61, 0xec, 0x11, 0x90, 0x9c, 0x16, 0x12, 0x36, 0xfb, 0xf9, 0xfd, 0x7f, 0x83, 0x19, 0x47,
	0x21, 0xd2, 0xa7, 0x44, 0x4f, 0x42, 0xd4, 0xe8, 0x66, 0x79, 0xaa, 0x53, 0xde, 0x59, 0xa6, 0xa1,
	0x8c, 0xd5, 0xd9, 0x60, 0x16, 0xe9, 0x18, 0xa7, 0xae, 0x48, 0x97, 0x1e, 0xc6, 0xd9, 0x1c, 0x75,
	0x24, 0x16, 0xca, 0x7b, 0x16, 0x73, 0x4c, 0x66, 0x32, 0xf7, 0xaa, 0x35, 0xcf, 0x46, 0xaa, 0x6a,
	0xcf, 0xdf, 0x5a, 0x6c, 0xcf, 0xaf, 0x48, 0xce, 0x59, 0x3b, 0xc1, 0xa5, 0x74, 0x68, 0x8f, 0xf6,
	0x0f, 0x02, 0x7b, 0xe6, 0x17, 0x6c, 0x5f, 0x6e, 0x01, 0xa7, 0xd5, 0xa3, 0xfd, 0xc3, 0x41, 0xd7,
	0xad, 0x1c, 0x77, 0xb4, 0x9d, 0x07, 0x3f, 0x1b, 0xfc, 0x9a, 0x1d, 0x63, 0x16, 0x4d, 0x44, 0x2e,
	0x43, 0x99, 0xe8, 0x08, 0x63, 0xe5, 0xfc, 0xb3, 0xd1, 0xe9, 0x2e, 0xf2, 0x1f, 0xee, 0x86, 0xbf,
	0xaf, 0xc1, 0x11, 0x66, 0x51, 0xed, 0xce, 0x47, 0xec, 0x44, 0x69, 0xcc, 0x17, 0x0d, 0xa2, 0x6d,
	0x09, 0x67, 0x47, 0x8c, 0x1f, 0xfd, 0xe0, 0xbe, 0x8e, 0x74, 0x6d, 0xf2, 0x87, 0x91, 0x22, 0x54,
	0xd8, 0x60, 0xfe, 0x37, 0x99, 0xd1, 0xf0, 0x76, 0xec, 0x37, 0x18, 0x9b, 0xd4, 0x26, 0x37, 0x57,
	0xab, 0x02, 0xc8, 0xba, 0x00, 0xb2, 0x29, 0x80, 0xbe, 0x18, 0xa0, 0xef, 0x06, 0xe8, 0x87, 0x01,
	0xba, 0x32, 0x40, 0x3f, 0x0d, 0xd0, 0x2f, 0x03, 0x64, 0x63, 0x80, 0xbe, 0x96, 0x40, 0x56, 0x25,
	0x90, 0x75, 0x09, 0x64, 0xda, 0xb1, 0x3f, 0x7b, 0xf9, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x66,
	0x0f, 0xb9, 0xab, 0x01, 0x00, 0x00,
}

func (this *Account) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Account)
	if !ok {
		that2, ok := that.(Account)
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
	if this.Name != that1.Name {
		return false
	}
	if !this.Exchange.Equal(that1.Exchange) {
		return false
	}
	if !this.ApiCredentials.Equal(that1.ApiCredentials) {
		return false
	}
	if !this.StarkCredentials.Equal(that1.StarkCredentials) {
		return false
	}
	if !this.EcdsaCredentials.Equal(that1.EcdsaCredentials) {
		return false
	}
	return true
}
func (this *Account) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.Account{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	if this.Exchange != nil {
		s = append(s, "Exchange: "+fmt.Sprintf("%#v", this.Exchange)+",\n")
	}
	if this.ApiCredentials != nil {
		s = append(s, "ApiCredentials: "+fmt.Sprintf("%#v", this.ApiCredentials)+",\n")
	}
	if this.StarkCredentials != nil {
		s = append(s, "StarkCredentials: "+fmt.Sprintf("%#v", this.StarkCredentials)+",\n")
	}
	if this.EcdsaCredentials != nil {
		s = append(s, "EcdsaCredentials: "+fmt.Sprintf("%#v", this.EcdsaCredentials)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringAccountData(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Account) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Account) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Account) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EcdsaCredentials != nil {
		{
			size, err := m.EcdsaCredentials.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccountData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.StarkCredentials != nil {
		{
			size, err := m.StarkCredentials.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccountData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.ApiCredentials != nil {
		{
			size, err := m.ApiCredentials.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccountData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Exchange != nil {
		{
			size, err := m.Exchange.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccountData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintAccountData(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAccountData(dAtA []byte, offset int, v uint64) int {
	offset -= sovAccountData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Account) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovAccountData(uint64(l))
	}
	if m.Exchange != nil {
		l = m.Exchange.Size()
		n += 1 + l + sovAccountData(uint64(l))
	}
	if m.ApiCredentials != nil {
		l = m.ApiCredentials.Size()
		n += 1 + l + sovAccountData(uint64(l))
	}
	if m.StarkCredentials != nil {
		l = m.StarkCredentials.Size()
		n += 1 + l + sovAccountData(uint64(l))
	}
	if m.EcdsaCredentials != nil {
		l = m.EcdsaCredentials.Size()
		n += 1 + l + sovAccountData(uint64(l))
	}
	return n
}

func sovAccountData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAccountData(x uint64) (n int) {
	return sovAccountData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Account) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Account{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Exchange:` + strings.Replace(fmt.Sprintf("%v", this.Exchange), "Exchange", "models.Exchange", 1) + `,`,
		`ApiCredentials:` + strings.Replace(fmt.Sprintf("%v", this.ApiCredentials), "APICredentials", "models.APICredentials", 1) + `,`,
		`StarkCredentials:` + strings.Replace(fmt.Sprintf("%v", this.StarkCredentials), "STARKCredentials", "models.STARKCredentials", 1) + `,`,
		`EcdsaCredentials:` + strings.Replace(fmt.Sprintf("%v", this.EcdsaCredentials), "ECDSACredentials", "models.ECDSACredentials", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAccountData(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Account) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountData
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
			return fmt.Errorf("proto: Account: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Account: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAccountData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccountData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exchange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountData
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
				return ErrInvalidLengthAccountData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccountData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Exchange == nil {
				m.Exchange = &models.Exchange{}
			}
			if err := m.Exchange.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApiCredentials", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountData
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
				return ErrInvalidLengthAccountData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccountData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ApiCredentials == nil {
				m.ApiCredentials = &models.APICredentials{}
			}
			if err := m.ApiCredentials.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StarkCredentials", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountData
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
				return ErrInvalidLengthAccountData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccountData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StarkCredentials == nil {
				m.StarkCredentials = &models.STARKCredentials{}
			}
			if err := m.StarkCredentials.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EcdsaCredentials", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountData
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
				return ErrInvalidLengthAccountData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccountData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EcdsaCredentials == nil {
				m.EcdsaCredentials = &models.ECDSACredentials{}
			}
			if err := m.EcdsaCredentials.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccountData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccountData
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
func skipAccountData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAccountData
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
					return 0, ErrIntOverflowAccountData
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
					return 0, ErrIntOverflowAccountData
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
				return 0, ErrInvalidLengthAccountData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAccountData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAccountData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAccountData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAccountData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAccountData = fmt.Errorf("proto: unexpected end of group")
)
