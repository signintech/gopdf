package gopdf

//Buff for pdf content
type Buff struct {
	position int
	datas    []byte
}

//Write : write []byte to buffer
func (b *Buff) Write(p []byte) (int, error) {
	for len(b.datas) < b.position+len(p) {
		b.datas = append(b.datas, 0)
	}
	i := 0
	max := len(p)
	for i < max {
		b.datas[i+b.position] = p[i]
		i++
	}
	b.position += i
	return 0, nil
}

//Len : len of buffer
func (b *Buff) Len() int {
	return len(b.datas)
}

//Bytes : get bytes
func (b *Buff) Bytes() []byte {
	return b.datas
}

//Position : get current postion
func (b *Buff) Position() int {
	return b.position
}

//SetPosition : set current postion
func (b *Buff) SetPosition(pos int) {
	b.position = pos
}
