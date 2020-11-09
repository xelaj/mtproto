package tl

const (
	WordLen   = 4           // размер слова в TL (32 бита)
	LongLen   = WordLen * 2 // int64 8 байт занимает
	DoubleLen = WordLen * 2 // float64 8 байт занимает

	// Блядские магические числа
	FuckingMagicNumber = 254  // 253 элемента максимум можно закодировать в массиве элементов
	ByteLenMagicNumber = 0xfe // ???

	// https://core.telegram.org/schema/mtproto
	CrcVector = 0x1cb5c415
	CrcFalse  = 0xbc799737
	CrcTrue   = 0x997275b5
	// CrcNull   = 0x56730bcc
)
