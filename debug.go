package gopdf

import "log"

func DebugSubType(b []byte) {
	//b := buff.Bytes()
	var max = len(ch)
	var i = 0
	for i < max {
		if b[i] != ch[i] {
			log.Fatalf("line: %d  real = %d  my = %d\n", i, ch[i], b[i])
		}
		i++
	}
}
