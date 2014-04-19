package pcap

import (
	"encoding/binary"
	"io"
	"time"
)

const DLT_IEEE802_15_4 = uint32(195)

// see http://wiki.wireshark.org/Development/LibpcapFileFormat
type FileHeader struct {
	//magic_number: used to detect the file format itself and the byte ordering. The writing application writes 0xa1b2c3d4 with it's native byte ordering format into this field.
	MagicNumber uint32

	//version_major, version_minor: the version number of this file format (current version is 2.4)
	VersionMajor uint16
	VersionMinor uint16

	//thiszone: the correction time in seconds between GMT (UTC) and the local timezone of the following packet header timestamps.
	TimeZone int32

	//sigfigs: in theory, the accuracy of time stamps in the capture; in practice, all tools set it to 0
	SigFigs uint32

	//snaplen: the "snapshot length" for the capture (typically 65535 or even more, but might be limited by the user)
	SnapLen uint32

	// See: http://www.tcpdump.org/linktypes.html
	LinkType uint32
}

type Packet struct {
	// porting from 'pcap_pkthdr' struct
	Time   time.Time // packet send/receive time
	Caplen uint32    // bytes stored in the file (caplen <= len)
	Len    uint32    // bytes sent/received

	Data []byte // packet data
}

// create a new file header with sane defaults
func NewFileHeader(linkType uint32) *FileHeader {
	return &FileHeader{
		MagicNumber:  0xa1b2c3d4,
		VersionMajor: 2,
		VersionMinor: 4,
		TimeZone:     0,
		SigFigs:      0,
		SnapLen:      65535,
		LinkType:     linkType,
	}
}

// create a new packet with defaults which reflect realtime capture
func NewPacket(data []byte, len uint32) *Packet {
	return &Packet{
		Time:   time.Now(),
		Caplen: len,
		Len:    len,
		Data:   data,
	}
}

type Writer struct {
	writer io.Writer
	buf    []byte
}

// NewWriter creates a Writer that stores output in an io.Writer.
// The FileHeader is written immediately.
func NewWriter(writer io.Writer, header *FileHeader) (*Writer, error) {
	w := &Writer{
		writer: writer,
		buf:    make([]byte, 24),
	}
	binary.LittleEndian.PutUint32(w.buf, header.MagicNumber)
	binary.LittleEndian.PutUint16(w.buf[4:], header.VersionMajor)
	binary.LittleEndian.PutUint16(w.buf[6:], header.VersionMinor)
	binary.LittleEndian.PutUint32(w.buf[8:], uint32(header.TimeZone))
	binary.LittleEndian.PutUint32(w.buf[12:], header.SigFigs)
	binary.LittleEndian.PutUint32(w.buf[16:], header.SnapLen)
	binary.LittleEndian.PutUint32(w.buf[20:], header.LinkType)
	if _, err := writer.Write(w.buf); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) Write(pkt *Packet) error {
	binary.LittleEndian.PutUint32(w.buf, uint32(pkt.Time.Unix()))
	binary.LittleEndian.PutUint32(w.buf[4:], uint32(pkt.Time.Nanosecond()))
	binary.LittleEndian.PutUint32(w.buf[8:], pkt.Caplen)
	binary.LittleEndian.PutUint32(w.buf[12:], pkt.Len)
	if _, err := w.writer.Write(w.buf[:16]); err != nil {
		return err
	}
	_, err := w.writer.Write(pkt.Data)
	return err
}
