package obinary

import (
	"gopkg.in/istreamdata/orientgo.v2"
	"gopkg.in/istreamdata/orientgo.v2/obinary/rw"
	"math/big"
)

func (f binaryRecordFormatV0) readDecimal(r *rw.ReadSeeker) interface{} {
	scale := int(r.ReadInt())
	value := big.NewInt(0).SetBytes(r.ReadBytes())
	return orient.Decimal{
		Scale: scale,
		Value: value,
	}
}

func (f binaryRecordFormatV0) writeDecimal(w *rw.Writer, o interface{}) {
	var d orient.Decimal
	switch v := o.(type) {
	case int64:
		d = orient.Decimal{Value: big.NewInt(v)}
	case *big.Int:
		d = orient.Decimal{Value: v}
	case orient.Decimal:
		d = v
	default:
		panic(orient.ErrTypeSerialization{Val: o, Serializer: f})
	}
	w.WriteInt(int32(d.Scale))    // scale value, 0 for ints
	w.WriteBytes(d.Value.Bytes()) // unscaled value
}
