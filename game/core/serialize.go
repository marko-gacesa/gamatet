// Copyright (c) 2020 by Marko Gaćeša

package core

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"slices"

	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/op"
)

type serializer struct {
	rawBuff  bytes.Buffer
	gzipBuff bytes.Buffer
}

func (s *serializer) Serialize(events *event.List) []byte {
	s.rawBuff.Reset()
	events.Range(func(e event.Event) {
		err := op.Write(&s.rawBuff, e)
		if err != nil {
			log.Printf("failed to serialize event %T: %s\n", e, err)
		}
	})

	if s.rawBuff.Len() == 0 {
		return nil
	}

	raw := s.rawBuff.Bytes()
	if raw[0] != 'Z' && len(raw) < 80 {
		return slices.Clone(raw)
	}

	s.gzipBuff.Reset()
	s.gzipBuff.Write([]byte{'Z'})
	zipper := gzip.NewWriter(&s.gzipBuff)
	_, _ = zipper.Write(raw)
	_ = zipper.Close()

	return slices.Clone(s.gzipBuff.Bytes())
}

func (s *serializer) Deserialize(data []byte, p event.Pusher) error {
	if len(data) == 0 {
		return nil
	}

	readEvents := func(reader io.Reader) error {
		for {
			e, err := op.Read(reader)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}

			p.Push(e)
		}
	}

	if data[0] == 'Z' {
		return _unzip(bytes.NewReader(data[1:]), readEvents)
	} else {
		return readEvents(bytes.NewReader(data))
	}
}

func _unzip(reader io.Reader, fn func(io.Reader) error) error {
	unzipper, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	err = fn(unzipper)
	if err != nil {
		return err
	}

	err = unzipper.Close()
	if err != nil {
		return err
	}

	return nil
}
