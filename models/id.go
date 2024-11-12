package models

var alphanum = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (m Model) NewPublicID(lens []int) string {
	sum := 0
	for _, l := range lens {
		sum += l
	}

	b := make([]rune, sum)
	for i := range b {
		b[i] = alphanum[m.rand.Intn(len(alphanum))]
	}

	pid := ""
	for i, l := range lens {
		pid += string(b[:l])
		if i < len(lens)-1 {
			pid += "-"
		}
		b = b[l:]
	}

	return pid
}
