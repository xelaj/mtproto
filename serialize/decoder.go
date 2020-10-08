package serialize

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/fatih/structtag"
	"github.com/k0kubun/pp"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
)

type Decoder struct {
	buf *bytes.Buffer
}

func NewDecoder(input []byte) *Decoder {
	return &Decoder{
		buf: bytes.NewBuffer(input),
	}
}

func (d *Decoder) PopLong() int64 {
	pp.Println("PopLong")
	val := make([]byte, LongLen)
	d.mustRead(val)
	return int64(binary.LittleEndian.Uint64(val))
}

func (d *Decoder) PopDouble() float64 {
	pp.Println("PopDouble")
	val := make([]byte, DoubleLen)
	d.mustRead(val)
	return math.Float64frombits(binary.LittleEndian.Uint64(val))
}

func (d *Decoder) PopInt() int32 {
	pp.Println("PopInt")
	val := make([]byte, WordLen)
	d.mustRead(val)
	return int32(binary.LittleEndian.Uint32(val))
}

func (d *Decoder) PopUint() uint32 {
	pp.Println("PopUint")
	val := make([]byte, WordLen)
	d.mustRead(val)
	return binary.LittleEndian.Uint32(val)
}

func (d *Decoder) PopInt128() *Int128 {
	pp.Println("PopInt128")
	val := d.PopRawBytes(Int128Len)
	return &Int128{big.NewInt(0).SetBytes(val)}
}

func (d *Decoder) PopInt256() *Int256 {
	pp.Println("PopInt256")
	val := d.PopRawBytes(Int256Len)
	return &Int256{big.NewInt(0).SetBytes(val)}
}

func (d *Decoder) PopRawBytes(size int) []byte {
	pp.Println("PopRawBytes")
	val := make([]byte, size)
	d.mustRead(val)
	return val
}

func (d *Decoder) PopMessage() []byte {
	pp.Println("PopMessage")
	var firstByte byte
	val := []byte{0}

	d.mustRead(val)
	firstByte = val[0]

	realSize := 0
	lenNumberSize := 0 // сколько байт занимаем число обозначающее длину массива
	if firstByte != FuckingMagicNumber {
		realSize = int(firstByte) // это tinyMessage по сути, первый байт является 8битным числом, которое представляет длину сообщения
		lenNumberSize = 1
	} else {
		// иначе это largeMessage с блядским магитческим числом 0xfe
		realSizeBuf := make([]byte, WordLen-1) // WordLen-1 т.к. 1 байт уже прочитали
		d.mustRead(realSizeBuf)
		realSizeBuf = append(realSizeBuf, 0x0) // добиваем до WordLen

		realSize = int(binary.LittleEndian.Uint32(realSizeBuf))
		lenNumberSize = WordLen
	}

	buf := make([]byte, realSize)
	d.mustRead(buf)
	readLen := lenNumberSize + realSize // lenNumberSize это сколько байт ушло на описание длины а realsize это сколько мы по факту прочитали
	if readLen%WordLen != 0 {
		voidBytes := make([]byte, 4-readLen%WordLen)
		d.mustRead(voidBytes) // читаем оставшиеся пустые байты. пустые, потому что длина слова 4 байта, может остаться 1,2 или 3 лишних байта
		for _, b := range voidBytes {
			if b != 0 {
				panic("some of bytes doesn't equal zero: " + fmt.Sprintf("%#v", voidBytes))
			}
		}
	}

	return buf
}

func (d *Decoder) GetRestOfMessage() []byte {
	return d.buf.Bytes()
}

func (d *Decoder) PopString() string {
	pp.Println("PopString")
	return string(d.PopMessage())
}

// TODO: непонятно, схерали int128 int256 это набор байт?
func (d *Decoder) PopBigInt() *big.Int {
	pp.Println("PopBigInt")
	return new(big.Int).SetBytes(d.PopMessage())
}

func (d *Decoder) PopBool() bool {
	pp.Println("PopBool")
	switch crc := d.PopUint(); crc {
	case crc_boolTrue:
		return true
	case crc_boolFalse:
		return false
	default:
		panic("not a bool value, actually: " + fmt.Sprintf("%#v", crc))
	}
}

func (d *Decoder) PopNull() interface{} {
	if d.PopUint() != crc_null {
		panic("not a bool value, actually")
	}
	return nil
}

// PopObj создает структуру исходя из кода объекта, который находится в буффере.
// Следует использовать только вкупе с функциями-генераторами, которые должны
// быть объявлены в CustomDecoders. поиск и создание объекта выполняется в том
// порядке, в котором были объявлены сами функции в CustomDecoders.
func (d *Decoder) PopObj() TL {
	pp.Println("PopObj")
	constructorID := d.PopCRC()

	var obj TL
	var isEnum bool
	var err error
	for _, f := range customDecoders {
		obj, isEnum, err = f(constructorID)
		if errs.IsNotFound(err) {
			continue
		}
		if err != nil {
			panic(err)
		}
		break
	}
	if obj == nil {
		panic(errs.NotFound("constructorID", fmt.Sprintf("%#v", constructorID)))
	}

	if !isEnum {
		d.PopToObjUsingReflection(obj, true)
	}

	return obj
}

func (d *Decoder) PopToObjUsingReflection(item TL, ignoreCRCReading bool) {
	pp.Println("PopToObjUsingReflection")
	if !ignoreCRCReading {
		crcCode := d.PopCRC()
		if crcCode != item.CRC() {
			pp.Println("PANIC", d.GetRestOfMessage())
			panic("invalid crc code: " + fmt.Sprintf("%#v", crcCode) + ", want: " + fmt.Sprintf("%#v", item.CRC()))
		}
	}

	// если есть метод DecodeFrom, то нам незачем париться
	if v, ok := item.(TLDecoder); ok {
		pp.Println("decoding native " + reflect.TypeOf(item).String())
		pp.Println(d.GetRestOfMessage())
		pp.Println("END decoding native " + reflect.TypeOf(item).String())
		v.DecodeFrom(d)
		return
	}

	value := reflect.ValueOf(item)
	if value.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	value = reflect.Indirect(value)
	if value.Kind() != reflect.Struct {
		panic("not recieving on struct: " + value.Type().String() + " -> " + value.Kind().String())
	}

	vtyp := value.Type()

	var optionalBitSet uint32

	for i := 0; i < value.NumField(); i++ {
		ftyp := value.Field(i).Type()

		// если в тегах указан flag значит нужно узнать, есть ли такой то бит, что бы уточнить, может вообще этот кусок пропустить?
		tags, err := structtag.Parse(string(vtyp.Field(i).Tag))
		dry.PanicIfErr(err)
		flagTag, err := tags.Get("flag")
		if err != nil {
			if err.Error() != "tag does not exist" {
				panic(err)
			}
		}
		if flagTag != nil {
			fmt.Println(vtyp.Field(i).Name + " имеет теги flag " + flagTag.String())
			triggerBit, err := strconv.Atoi(flagTag.Name)
			dry.PanicIfErr(err)
			if optionalBitSet&(1<<triggerBit) == 0 {
				fmt.Println("бит не задан так что пропускаем")
				continue
			}

			if dry.StringInSlice("encoded_in_bitflags", flagTag.Options) {
				fmt.Println("значение находится в битфлагах")
				value.Field(i).Set(reflect.ValueOf(true).Convert(ftyp))
				continue
			}

			fmt.Println("у этого филда задан бит тег!")
		}
		fmt.Println("парсим " + vtyp.Field(i).Name)
		switch value.Field(i).Kind() {
		case reflect.Int64:
			println("это int64!")
			value.Field(i).Set(reflect.ValueOf(d.PopLong()).Convert(ftyp))
		case reflect.Uint32: // это применимо так же к енумам
			println("это uint32!")
			value.Field(i).Set(reflect.ValueOf(d.PopUint()).Convert(ftyp))
		case reflect.Int32:
			println("это int32!")
			value.Field(i).Set(reflect.ValueOf(d.PopInt()).Convert(ftyp))
		case reflect.Bool:
			println("это bool!")
			value.Field(i).Set(reflect.ValueOf(d.PopBool()).Convert(ftyp))
		case reflect.String:
			println("это string!")
			value.Field(i).Set(reflect.ValueOf(d.PopString()).Convert(ftyp))
		case reflect.Struct:
			println("это struct!")
			if vtyp.Field(i).Name == "__flagsPosition" {
				optionalBitSet = d.PopUint()
				fmt.Printf("нашли битфлаг! %b\n", optionalBitSet)
				continue
			}
			fieldValue := reflect.New(ftyp).Elem().Interface().(TL)
			d.PopToObjUsingReflection(fieldValue, false)
			value.Field(i).Set(reflect.ValueOf(fieldValue).Convert(ftyp))

		case reflect.Slice:
			println("это slice!")
			if _, ok := value.Field(i).Interface().([]byte); ok {
				value.Field(i).Set(reflect.ValueOf(d.PopMessage()))
			} else {
				value.Field(i).Set(reflect.ValueOf(d.PopVector(ftyp.Elem())).Convert(ftyp))
			}
		case reflect.Ptr:
			println("это ptr!")
			// если поинтер то это структура на что-то
			if _, ok := value.Field(i).Interface().(*Int128); ok {
				value.Field(i).Set(reflect.ValueOf(d.PopInt128()))
			} else if _, ok := value.Field(i).Interface().(*Int256); ok {
				value.Field(i).Set(reflect.ValueOf(d.PopInt256()))
			} else if _, ok := value.Field(i).Interface().(TL); ok {
				value.Field(i).Set(reflect.New(value.Field(i).Type().Elem()))
				d.PopToObjUsingReflection(value.Field(i).Interface().(TL), false)

			} else {
				panic("неизвестная штука: " + value.Field(i).Type().String())
			}

		default:
			panic("неизвестная штука: " + value.Field(i).Type().String())
		}
	}

}

func (d *Decoder) PopCRC() uint32 {
	pp.Println("PopCRC")
	return d.PopUint() // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (d *Decoder) PopVector(as reflect.Type) interface{} {
	pp.Println("PopVector")
	constructorID := d.PopCRC()

	if constructorID != crc_vector {
		pp.Println(d.buf.Bytes())
		panic("not a vector: " + fmt.Sprintf("%#v", constructorID) + " want: 0x1cb5c415")
	}
	size := int(d.PopUint())

	x := reflect.MakeSlice(reflect.SliceOf(as), size, size)

	for i := 0; i < size; i++ {
		var v interface{}

		switch as.Kind() {
		case reflect.Bool:
			v = d.PopBool()
		case reflect.String:
			v = d.PopString()
		case reflect.Int8, reflect.Int16, reflect.Int32:
			v = d.PopInt()
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			v = d.PopUint()
		case reflect.Struct:
			v = d.PopObj()
		case reflect.Int64:
			v = d.PopLong()
		case reflect.Slice:
			if as.Elem().Kind() == reflect.Uint8 { // []byte
				v = d.PopMessage()
			} else {
				v = d.PopVector(as.Elem())
			}
		case reflect.Ptr:
			n := reflect.New(as.Elem()).Interface().(TL)
			d.PopToObjUsingReflection(n, false)
			v = n
		default:
			panic("как обрабатывать? " + as.String())
		}

		x.Index(i).Set(reflect.ValueOf(v))
	}

	return x.Interface()
}

func (d *Decoder) mustRead(into []byte) {
	if len(into) == 0 {
		return
	}
	n, err := d.buf.Read(into)
	if err != nil {
		panic(fmt.Errorf("read %d bytes: %w", n, err))
	}
	dry.PanicIf(n != len(into), fmt.Sprintf("expected to read equal %v bytes, got %v", len(into), n))
}
