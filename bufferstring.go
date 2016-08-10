package bufferstring

import (
	"sync"
	"unicode/utf8"

	"github.com/quanzo/gservice/bufferint"
	"github.com/quanzo/gservice/status"
)

type BufferString struct {
	buffer         []rune
	size           int
	addSpace       int
	lock           *sync.Mutex
	modeThreadSafe bool
}

// инициализация
func (this *BufferString) init(buffSize int, addSpace int) {
	this.size = 0
	if buffSize < 0 {
		buffSize = BUFFER_SIZE
	}
	if addSpace < 0 {
		addSpace = 0
	}
	this.addSpace = addSpace
	this.buffer = make([]rune, buffSize, buffSize+addSpace)
	this.lock = new(sync.Mutex)
	this.modeThreadSafe = true
}

// Инициализация из строки.
func (this *BufferString) initFromString(s string, addSpace int) {
	sLen := utf8.RuneCountInString(s)
	this.init(sLen, addSpace)
	if sLen > 0 {
		copy(this.buffer, []rune(s))
		this.size = sLen
	} else {
		this.size = 0
	}
} // end initFromString

//*****************************************************************************

// Увеличить размер буфера.
func (this *BufferString) alloc(newSize int) {
	if len(this.buffer) < newSize {
		this.buffer = append(this.buffer, make([]rune, newSize-len(this.buffer)+this.addSpace)...)
	}
}

// Получить одно значение из буфера с индексом i.
func (this *BufferString) one(i int) (rune, error) {
	if i >= 0 && i < this.size {
		return this.buffer[i], nil
	} else {
		return -1, status.New(0, "Index out of range.", true, false)
	}
}

// Заменить часть строки начиная с start длинной count на символы sRunes.
func (this *BufferString) replace(start int, count int, sRunes *[]rune) int {
	var (
		sLen int
	)
	if start < 0 {
		start = 0
	}
	if start > this.size {
		start = this.size
	}
	if count < 0 {
		count = 0
	}
	if count > (this.size - start) {
		count = this.size - start
	}
	if sRunes != nil {
		sLen = len(*sRunes)
	} else {
		sLen = 0
	}

	newBuffSize := this.size + sLen - count
	if len(this.buffer) < newBuffSize {
		this.alloc(newBuffSize)
	}
	if sLen != count {
		copy(this.buffer[start+sLen:], this.buffer[start+count:])
	}
	this.size = newBuffSize
	if sLen > 0 {
		copy(this.buffer[start:], *sRunes)
	}
	return this.size
} // end replace

// Получить подстроку: начиная с символа start выдать count символов.
func (this *BufferString) substr(start int, count int) []rune {
	if start > this.size || count <= 0 || start < 0 {
		return nil
	} else {
		if start+count > this.size {
			return this.buffer[start:this.size]
		} else {
			return this.buffer[start:(start + count)]
		}
	}
} // end substr

// Поиск цепочки символов в буфере.
func (this *BufferString) find(needle *[]rune, start int, back bool) int {
	sNeedle := len(*needle)
	if start == -1 {
		if back {
			start = this.size
		} else {
			start = 0
		}
	}
	if (!back && (start+sNeedle) > this.size) || (back && (start-sNeedle+1) < 0) {
		return -1
	} else {
		var (
			i, i_start, i_end int
		)
		i = start
		for i >= 0 && i <= this.size {
			if back {
				i_start = i - sNeedle
				i_end = i
			} else {
				i_start = i
				i_end = i + sNeedle
			}
			if i_start < 0 || i_end > this.size {
				return -1
			} else {
				if this.Equal(this.buffer[i_start:i_end], *needle) {
					return i_start
				}
			}
			if back {
				i--
			} else {
				i++
			}
		}
		return -1
	}
} // end find

// Возвращает первую встретившуюся строку из needle. Возвращает позицию строки и его номер в needle.
func (this *BufferString) findFirst(needle *[]string, start int) (int, int) {
	var (
		p, res_p int
		i, res_i int
		val      string
		valRune  []rune
	)
	res_p = -1
	res_i = -1
	for i, val = range *needle {
		valRune = []rune(val)
		p = this.find(&valRune, start, false)
		if p > -1 && ((res_p != -1 && res_p > p) || res_p == -1) {
			res_p = p
			res_i = i
		}
	}
	return res_p, res_i
} // end firstFind

// Найти все вхождения цепочки символов в буфере. Результат будет добавлен в буфер res. Возвращено кол-во найденых цепочек символов.
func (this *BufferString) findAll(needle *[]rune, res *bufferint.BufferInt) int {
	var (
		needle_len int
		start      int = 0
		pos        int
		res_count  int = 0
	)
	needle_len = len(*needle)
	if needle_len > 0 {
		for pos = this.find(needle, start, false); pos != -1; pos = this.find(needle, start, false) {
			res.Append(pos)
			start = pos + needle_len
			res_count++
		}
	}
	return res_count
} // end findall

/* Найти строки search и заменить их на replace.
При этом строка search[0] будет заменена на replace[0]. Остальные замены будут выполнены аналогично.
Если количество элементов в search больше, чем в replace, то строки (search) с индексом больше len(replace) будут заменены на последнюю строку в replace.
*/
func (this *BufferString) findReplace(search []string, replace []string) {
	if search != nil && replace != nil && len(search) > 0 && len(replace) > 0 {
		var (
			i, ir, p, delta_1            int
			s0                           string
			sRunes, rRunes               []rune
			approx_size, approx_addspace int
			find_count                   int
			new_size                     int
			err                          error
			lenSearchStr, lenReplaceStr  int
		)

		approx_size = this.size / 4
		if approx_size > 64 {
			approx_size = 64
		}
		if approx_size < 8 {
			approx_size = 8
		}
		approx_addspace = approx_size / 2
		buffPos := bufferint.New(approx_size, approx_addspace)

		for i, s0 = range search {
			sRunes = []rune(s0)
			lenSearchStr = len(sRunes)
			if lenSearchStr > 0 {
				buffPos.Empty()
				find_count = this.findAll(&sRunes, buffPos)
				if find_count > 0 {
					if i >= len(replace) {
						rRunes = []rune(replace[len(replace)-1])
					} else {
						rRunes = []rune(replace[i])
					}
					lenReplaceStr = len(rRunes)

					delta_1 = lenReplaceStr - lenSearchStr
					new_size = this.size + delta_1*find_count
					if new_size > len(this.buffer) { // new allocate buffer
						this.alloc(new_size)
					}

					if delta_1 > 0 {
						var (
							prev_p   int
							copy_pos int
						)
						for ir = find_count - 1; ir >= 0; ir-- {
							p, err = buffPos.One(ir)
							if err == nil {
								prev_p, err = buffPos.One(ir + 1)
								copy_pos = p + delta_1*ir
								if err == nil {
									copy(this.buffer[copy_pos+lenReplaceStr:], this.buffer[p+lenSearchStr:prev_p])
								} else {
									copy(this.buffer[copy_pos+lenReplaceStr:], this.buffer[p+lenSearchStr:this.size])
								}
								copy(this.buffer[copy_pos:], rRunes)
							}
						} // end for find_count
						// оконечные данные
						copy(this.buffer[0:], this.buffer[0:p])
					} else {
						var (
							next_p           int
							cumulative_delta int = 0
						)
						for ir = 0; ir < find_count; ir++ {
							p, err = buffPos.One(ir)
							if err == nil {
								copy(this.buffer[p+cumulative_delta:], rRunes)
								next_p, err = buffPos.One(ir + 1)
								if err == nil {
									copy(this.buffer[p+cumulative_delta+lenReplaceStr:], this.buffer[p+lenSearchStr:next_p])
								} else {
									copy(this.buffer[p+cumulative_delta+lenReplaceStr:], this.buffer[p+lenSearchStr:this.size])
								}
								cumulative_delta += delta_1
							}
						} // end for ir
					}
					this.size = new_size
				}
			} // end if len search
		} // end for
	}
} // end findReplace

// Подготовить маску
func (this *BufferString) prepareMask(mask string, symStar rune, symQuest rune) (*[]rune, *[]byte) {
	arMask := []rune(mask)
	sizeMask := len(arMask)
	arMaskType := make([]byte, sizeMask)
	if sizeMask > 0 { // распознаем и определим символы в маске. 0 - обычный символ, 1 - символ маски *, 2 - символ маски ?, 255 - символ не учитывать
		var (
			i int
		)
		for i = 0; i < sizeMask; i++ {
			if arMask[i] == '\\' {
				if i+1 < sizeMask {
					if arMask[i+1] == symStar || arMask[i+1] == symQuest {
						arMaskType[i] = 255
						arMaskType[i+1] = 0
						i++
					} else {
						arMaskType[i] = 0
					}
				} else {
					arMaskType[i] = 0
				}
			} else {
				switch arMask[i] {
				case symStar:
					{
						arMaskType[i] = 1
						break
					}
				case symQuest:
					{
						arMaskType[i] = 2
						break
					}
				default:
					{
						arMaskType[i] = 0
						break
					}
				}
			}
		} // end for
	}
	return &arMask, &arMaskType
} // end prepareMask

/* Ищет в буфере строковой фрагмент, определенный маской.
Символ * обозначает любые возможные символы
Символ ? обозначает один символ.
Экранирование символов маски с помощью \

arMask и arMaskType - результаты работы prepareMask

*/
func (this *BufferString) findMask(arMask *[]rune, arMaskType *[]byte, start int) (int, string) {
	if len(*arMask) > 0 {
		var (
			i           int
			maskCounter int
			resPos      int
			nextInc     int
			incI        bool
		)
		sizeMask := len(*arMask)

		resPos = -1
		i = start

		for i < this.size {
			incI = true
			for (*arMaskType)[maskCounter] == 255 && maskCounter < sizeMask {
				maskCounter++
			}
			switch (*arMaskType)[maskCounter] {
			case 0:
				{
					if this.buffer[i] != (*arMask)[maskCounter] {
						maskCounter = 0
						resPos = -1
					} else {
						if resPos == -1 {
							resPos = i
						}
						maskCounter++
					}
					break
				}
			case 2: // ?
				{
					if resPos == -1 {
						resPos = i
					}
					maskCounter++
					break
				}
			case 1: // *
				{
					if resPos == -1 {
						resPos = i
					}
					if maskCounter+1 < sizeMask {
						nextInc = 1
						if (*arMaskType)[maskCounter+nextInc] == 255 {
							nextInc = 2
						}

						switch (*arMaskType)[maskCounter+nextInc] {
						case 0:
							{
								if this.buffer[i] == (*arMask)[maskCounter+nextInc] {
									maskCounter++
									incI = false
								}
								break
							}
						case 2, 1: // если следующий элемент * или ? то переходим на него
							{
								maskCounter++
								break
							}
						}
					}
					break
				}
			}
			if maskCounter >= sizeMask {
				return resPos, string(this.buffer[resPos : i+1])
			}
			if incI {
				i++
			}
		} // end for
		if resPos != -1 && maskCounter >= sizeMask-1 {
			return resPos, string(this.buffer[resPos:])
		} else {
			return -1, ""
		}
	} else {
		return -1, ""
	}
} // end func

//******************************************************************************

// добавить строку в буфер
func (this *BufferString) AppendString(s ...string) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}

	sStrings := len(s)
	if sStrings > 0 {
		var (
			s0  string
			sr0 []rune
			r   rune
		)
		for _, s0 = range s {
			sr0 = []rune(s0)
			if len(sr0) > 0 {
				if len(this.buffer) > this.size {
					for _, r = range sr0 {
						if len(this.buffer) > this.size {
							this.buffer[this.size] = r
						} else {
							this.buffer = append(this.buffer, r)
						}
						this.size++
					}
				} else {
					this.buffer = append(this.buffer, sr0...)
					this.size = len(this.buffer)
				}
			}
		} // end for
	}
}

func (this *BufferString) AppendRune(s ...rune) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}

	sizeInput := len(s)
	if sizeInput > 0 {
		newSize := this.size + sizeInput
		if newSize > len(this.buffer) {
			this.alloc(newSize + this.addSpace)
		}
		copy(this.buffer[this.size:], s)
		this.size += sizeInput
	}
}

// Добавить в буфер значения из другого буфера. Возвращает количество добавленых записей.
func (this *BufferString) AppendBuffer(buff *BufferString, checked func(index int, value rune) bool) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	inputLength := buff.Length()
	if inputLength > 0 {
		// вычислим размер
		newSize := this.size + inputLength
		if newSize > len(this.buffer) {
			this.alloc(newSize)
		}
		if checked == nil {
			checked = func(index int, value rune) bool {
				return true
			}
		}
		var (
			i, res_count int
			v            rune
			e            error
		)
		for i = 0; i < inputLength; i++ {
			if v, e = buff.One(i); e == nil {
				if checked(i, v) {
					this.buffer[this.size] = v
					this.size++
					res_count++
				}
			} else {
				break
			}
		}
		return res_count
	} else {
		return 0
	}
}

func (this *BufferString) One(i int) (rune, error) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	return this.one(i)
}

// Выбрать и вернуть из буфера цепочку данных. Возвращаемая цепочка из буфера будет удалена.
func (this *BufferString) Pop(start int, count int) string {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	res := string(this.substr(start, count))
	_ = this.replace(start, count, nil)
	return res
}

// вставить строку s перед символом before_pos. Возвращает новый размер буфера.
func (this *BufferString) Insert(s string, before_pos int) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	r := []rune(s)
	return this.replace(before_pos, 0, &r)
}

// удалить count символов начиная с символа start. Возвращает кол-во удаленных символов.
func (this *BufferString) Delete(start int, count int) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	s := this.size
	r := []rune("")
	return s - this.replace(start, count, &r)
} // end func Delete

// заменить часть строки начиная с start длинной count на строку s. Возвращает новую длинну буфера.
func (this *BufferString) Replace(start int, count int, s string) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	r := []rune(s)
	return this.replace(start, count, &r)
}

// получить подстроку: начиная с символа start выдать count символов
func (this *BufferString) Substr(start int, count int) string {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	res := this.substr(start, count)
	if res == nil {
		return ""
	} else {
		return string(res)
	}
}

// прямой поиск подстроки needle в строке. Возвращает позицию подстроки или -1
func (this *BufferString) Find(needle string, start int) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	n_runes := []rune(needle)
	return this.find(&n_runes, start, false)
} // end Strpos

// обратный поиск подстроки needle в строке.  Возвращает позицию подстроки или -1
func (this *BufferString) FindReverse(needle string, start int) int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	n_runes := []rune(needle)
	return this.find(&n_runes, start, true)
} // end FindReverse

// Найти и заменить.
func (this *BufferString) FindReplace(search []string, replace []string) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	this.findReplace(search, replace)
} // end FindReverse

/* Ищет в буфере строковой фрагмент, определенный маской.
Символ * обозначает любые возможные символы.
Символ ? обозначает один символ.
Экранирование символов маски с помощью \.
*/
func (this *BufferString) FindMask(mask string, start int) (int, string) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	m1, m2 := this.prepareMask(mask, '*', '?')
	return this.findMask(m1, m2, start)
}

/* Ищет в буфере строковой фрагмент, определенный маской.
Символы маски - любое количество символов, одиночный символ - задаются при вызове метода.
*/
func (this *BufferString) FindMaskAdv(mask string, start int, symStar rune, symQuest rune) (int, string) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	m1, m2 := this.prepareMask(mask, symStar, symQuest)
	return this.findMask(m1, m2, start)
}

// Возвращает первую встретившуюся строку из needle. Возвращает позицию в строке и номер в needle.
func (this *BufferString) FindFirst(needle []string, start int) (int, int) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	return this.findFirst(&needle, start)
}

// Определяет идентичность двух срезов []rune.
func (this *BufferString) Equal(q, w []rune) bool {
	if q == nil && w == nil {
		return true
	}
	if q == nil || w == nil || len(q) != len(w) {
		return false
	} else {
		s := len(q)
		for i := 0; i < s; i++ {
			if q[i] != w[i] {
				return false
			}
		}
		return true
	}
} // end Equal

// Очистить буфер.
func (this *BufferString) Empty() {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	this.size = 0
}

// Выполнить функцию к каждому элементу буфера, начиная с позиции start.
func (this *BufferString) Walk(start int, count int, f func(index int, value *rune)) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	if start < this.size && f != nil {
		var (
			i int
			v *rune
		)
		for i = start; i < start+count && i < this.size; i++ {
			v = &this.buffer[i]
			f(i, v)
		}
	}
}

// Фильтрация буфера. Если функция filter возвращает false, то элемент удаляется из буфера.
func (this *BufferString) Filter(filter func(index int, value rune) bool) {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	if filter != nil {
		var (
			i, delCount int
			v           rune
		)
		for i = 0; i < this.size; i++ {
			v = this.buffer[i]
			if !filter(i+delCount, v) { // excluding from the buffer
				this.replace(i, 1, nil)
				i--
				delCount++
			}
		}
	}

}

//*****************************************************************************

// Вернуть данные в виде строки.
func (this *BufferString) String() string {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	if this.size == 0 {
		return ""
	} else {
		return string(this.buffer[0:this.size])
	}
}

// Получить данные в виде slice символов.
func (this *BufferString) GetCopy() []rune {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	res := make([]rune, this.size)
	copy(res, this.buffer[0:this.size])
	return res
}

// Вернуть текущую длину данных в буфере.
func (this *BufferString) Length() int {
	if this.modeThreadSafe {
		this.lock.Lock()
		defer this.lock.Unlock()
	}
	return this.size
}

// Установить многопоточный режим.
func (this *BufferString) SetModeThreadSafe(m bool) {
	this.modeThreadSafe = m
}

// Вернуть многопоточный режим.
func (this *BufferString) GetModeThreadSafe() bool {
	return this.modeThreadSafe
}
