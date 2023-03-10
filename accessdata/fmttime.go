package accessdata

import "time"

func FmtTimestamp(t time.Time) string {
	buf := []byte{}
	t = t.UTC()
	year, month, day := t.Date()
	itoa(&buf, year, 4)
	buf = append(buf, '-')
	itoa(&buf, int(month), 2)
	buf = append(buf, '-')
	itoa(&buf, day, 2)
	buf = append(buf, ' ')

	hour, min, sec := t.Clock()
	itoa(&buf, hour, 2)
	buf = append(buf, ':')
	itoa(&buf, min, 2)
	buf = append(buf, ':')
	itoa(&buf, sec, 2)
	//if l.flag&Lmicroseconds != 0 {
	buf = append(buf, '.')
	itoa(&buf, t.Nanosecond()/1e3, 6)
	//}
	//buf = append(buf, ' ')
	return string(buf)
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
