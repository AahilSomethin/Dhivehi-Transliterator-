package transliterator

import "unsafe"

func Transliterate(input string) string {
	n := len(input)
	buf := make([]byte, n*2)
	w := 0
	prevIdx := -1

	i := 0
	for i < n {
		if input[i] != 0xDE || i+1 >= n {
			goto nonThaana
		}

		{
			idx := uint(input[i+1]) - 0x80

			if akuruMask>>idx&1 == 0 {
				buf[w] = input[i]
				buf[w+1] = input[i+1]
				w += 2
				prevIdx = int(idx)
				i += 2
				continue
			}

			if i+3 < n && input[i+2] == 0xDE {
				nextIdx := uint(input[i+3]) - 0x80

				if filiMask>>nextIdx&1 != 0 {
					if idx == ainuIdx {
						fili := filiValues[nextIdx]
						buf[w] = fili[0]
						buf[w+1] = '\''
						w += 2
						w += copy(buf[w:], fili[1:])
					} else {
						w += copy(buf[w:], akuruValues[idx])
						w += copy(buf[w:], filiValues[nextIdx])
					}
					prevIdx = int(nextIdx)
					i += 4
					continue
				}

				if nextIdx == sukunIdx {
					if idx == alifuIdx || idx == shaviyaniIdx {
						if i+5 < n && input[i+4] == 0xDE {
							afterIdx := uint(input[i+5]) - 0x80
							if akuruMask>>afterIdx&1 != 0 && len(akuruValues[afterIdx]) > 0 {
								buf[w] = akuruValues[afterIdx][0]
								w++
								prevIdx = int(sukunIdx)
								i += 4
								continue
							}
						}
						buf[w] = 'h'
						w++
						prevIdx = int(sukunIdx)
						i += 4
						continue
					}

					if idx == noonuIdx {
						if i+5 < n && input[i+4] == 0xDE {
							afterIdx := uint(input[i+5]) - 0x80
							if afterIdx == meymuIdx || afterIdx == baaIdx || afterIdx == paviyaniIdx {
								buf[w] = akuruValues[afterIdx][0]
								w++
								prevIdx = int(sukunIdx)
								i += 4
								continue
							}
						}
						w += copy(buf[w:], akuruValues[idx])
						prevIdx = int(sukunIdx)
						i += 4
						continue
					}

					if sukunOvrdMask>>idx&1 != 0 {
						w += copy(buf[w:], sukunOvrdValues[idx])
					} else {
						w += copy(buf[w:], akuruValues[idx])
					}
					prevIdx = int(sukunIdx)
					i += 4
					continue
				}
			}

			// Bare akuru

			if idx == noonuIdx && filiMask>>uint(prevIdx)&1 != 0 {
				if i+3 < n && input[i+2] == 0xDE {
					peekIdx := uint(input[i+3]) - 0x80
					if akuruMask>>peekIdx&1 != 0 {
						buf[w] = 'n'
						buf[w+1] = '\''
						w += 2
						prevIdx = int(idx)
						i += 2
						continue
					}
				}
			}

			if idx == raaIdx {
				afterFiliOrAkuru := filiOrAkuruMask>>uint(prevIdx)&1 != 0
				beforeAkuru := false
				if i+3 < n && input[i+2] == 0xDE {
					peekIdx := uint(input[i+3]) - 0x80
					beforeAkuru = akuruMask>>peekIdx&1 != 0
				}
				if afterFiliOrAkuru || beforeAkuru {
					buf[w] = 'r'
					w++
					prevIdx = int(idx)
					i += 2
					continue
				}
			}

			if akuruNameMask>>idx&1 != 0 {
				w += copy(buf[w:], akuruNameValues[idx])
			}
			prevIdx = int(idx)
			i += 2
			continue
		}

	nonThaana:
		b := input[i]
		if b == 0xD8 && i+1 < n {
			switch input[i+1] {
			case 0x8C:
				buf[w] = ','
				w++
				prevIdx = -1
				i += 2
				continue
			case 0x9B:
				buf[w] = ';'
				w++
				prevIdx = -1
				i += 2
				continue
			case 0x9F:
				buf[w] = '?'
				w++
				prevIdx = -1
				i += 2
				continue
			}
		}

		if b < 0x80 {
			buf[w] = b
			w++
			prevIdx = -1
			i++
		} else if b < 0xC0 {
			buf[w] = b
			w++
			prevIdx = -1
			i++
		} else if b < 0xE0 {
			buf[w] = input[i]
			buf[w+1] = input[i+1]
			w += 2
			prevIdx = -1
			i += 2
		} else if b < 0xF0 {
			copy(buf[w:], input[i:i+3])
			w += 3
			prevIdx = -1
			i += 3
		} else {
			copy(buf[w:], input[i:i+4])
			w += 4
			prevIdx = -1
			i += 4
		}
	}

	return unsafe.String(unsafe.SliceData(buf[:w]), w)
}
