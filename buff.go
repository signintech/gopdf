package gopdf

type Buff struct {
	position int
	datas    []byte
}

func (me *Buff) Write(p []byte) (int, error) {
	for len(me.datas) < me.position+len(p) {
		me.datas = append(me.datas, 0)
	}
	i := 0
	max := len(p)
	for i < max {
		me.datas[i+me.position] = p[i]
		i++
	}
	me.position += i
	return 0, nil
}

func (me *Buff) Len() int {
	return len(me.datas)
}

func (me *Buff) Bytes() []byte {
	return me.datas
}

func (me *Buff) Position() int {
	return me.position
}

func (me *Buff) SetPosition(pos int) {
	me.position = pos
}
