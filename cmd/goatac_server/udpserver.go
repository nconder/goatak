package main

import (
	"net"

	"github.com/kdudkov/goatak/cot"
	"github.com/kdudkov/goatak/xml"
)

func (app *App) ListenUDP(addr string) error {
	p, err := net.ListenPacket("udp", addr)

	if err != nil {
		app.Logger.Error(err)
		return err
	}

	buf := make([]byte, 65535)

	for app.ctx.Err() == nil {
		n, _, err := p.ReadFrom(buf)
		if err != nil {
			app.Logger.Errorf("read error: %v", err)
			return err
		}

		evt := &cot.Event{}
		if err := xml.Unmarshal(buf[:n], evt); err != nil {
			app.Logger.Errorf("decode error: %v", err)
			continue
		}

		dat := make([]byte, n)
		copy(dat, buf[:n])
		app.ch <- &Msg{event: evt, dat: dat}
	}

	return nil
}
