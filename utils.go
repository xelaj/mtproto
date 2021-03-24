// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"context"
	"io"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
)

type any = interface{}
type null = struct{}

// это неофициальная информация, но есть подозрение, что список датацентров АБСОЛЮТНО идентичный для всех
// приложений. Несмотря на это, любой клиент ОБЯЗАН явно указывать список датацентров, ради надежности.
// данный список лишь эксперементальный и не является частью протокола.
func defaultDCList() map[int]string {
	return map[int]string{
		1: "149.154.175.58:443",
		2: "149.154.167.50:443",
		3: "149.154.175.100:443",
		4: "149.154.167.91:443",
		5: "91.108.56.151:443",
	}
}

func MessageRequireToAck(msg tl.Object) bool {
	switch msg.(type) {
	case /**objects.Ping,*/ *objects.MsgsAck:
		return false
	default:
		return true
	}
}

func CloseOnCancel(ctx context.Context, c io.Closer) {
	go func() {
		<-ctx.Done()
		c.Close()
	}()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
