# BufferString

Пакет предназначен для манипуляции с текстовыми данными без выделения дополнительной памяти. Может использоваться в многопоточной программе.

### Создание

__Создание из строки__

``text := bufferstring.NewFromString(".....text.....", 5)``

Второй параметр задает шаг увеличения буфера по необходимости. 

__Создание пустого буфера__

``emptytext := bufferstring.NewFromString(100, 5)``

Создается пустой буфер размером 100 + 5 = 105 символов и с шагом увеличения 5.

## Функционал

| Функция | Пояснение | Пример |
|-------------|-----------|--------|
|AppendString(s ...string)|Добавить строку(и) в конце буфера|``text.AppendString("string1")`` ``text.AppendString("string1", "string2", "string3")``|
|AppendRune(s ...rune)|Добавить символ(ы) к буферу|``text.AppendRune('r', 'u', 'n', 'e')``|
|AppendBuffer(buff \*BufferString, checked func(index int, value rune) bool) int|Добавить в буфер значения из другого буфера. Возвращает количество добавленых записей.||
|One(i int) (rune, error)|Получить одиночный символ из буфера.|``runeSym, err := text.One(0)``|
|Pop(start int, count int) string|Выбрать и вернуть из буфера цепочку данных. Возвращаемая цепочка из буфера будет удалена.|``str := text.Pop(0, 2)``|
|Insert(s string, before_pos int) int|Вставляет строку перед символом before_pos. Возвращает новый размер буфера.|``newBuffSize := text.Insert("start text", 0)``|
|``Delete(start int, count int) int``|Удаляет цепочку из _count_ символов начиная с символа _start_. Возвращает кол-во удаленных символов.|``delCount := text.Delete(0, 10)``|
|``Replace(start int, count int, s string) int``|Заменить часть строки начиная с _start_ длинной _count_ на строку _s_. Возвращает новую длинну буфера.||
|``Substr(start int, count int) string``|Возвращает подстроку начиная с символа _start_||
|``Find(needle string, start int) int``|Возвращает позицию строки _needle_. Поиск начинается позиции _start_. Если строка не найдена, функция вернет _-1_ ||
|``FindReverse(needle string, start int) int``|Обратный поиск строки. Вернет позицию или _-1_. Поиск идет от позиции старт к началу строки.||
|``FindReplace(search []string, replace []string)``|Найти и заменить. Ищет строки _search_ и заменяет на _replace_. При этом строка _search[0]_ будет заменена на _replace[0]_. Остальные замены будут выполнены аналогично. Если количество элементов в _search_ больше, чем в replace, то строки (_search_) с индексом больше длинны _replace_ будут заменены на последнюю строку в _replace_.||
|``FindMask(mask string, start int) (int, string)``|Ищет в буфере строковой фрагмент, определенный маской. Символ * обозначает любые возможные символы. Символ ? обозначает один символ. Экранирование символов маски с помощью \.|``pos, fmask := text.FindMask("*Wh?\?*", 0)``|
|``FindMaskAdv(mask string, start int, symStar rune, symQuest rune) (int, string)``|Ищет строковой фрагмент, определенный маской. Символ * и ? задаются при вызове.||
|``FindFirst(needle []string, start int) (int, int)``|Ищет первую встретившуюся строку из _needle_. Возвращает позицию в строке и номер в _needle_.||
|``Equal(q, w []rune) bool``|Определяет идентичность двух срезов _[]rune_.||
|``Empty()``|Очистить||
|``Walk(start int, count int, f func(index int, value *rune))``|Выполнить функцию к каждому элементу буфера, начиная с позиции start.||
|``Filter(filter func(index int, value rune) bool)``|Фильтрация буфера. Если функция filter возвращает false, то элемент удаляется из буфера.||
|``String() string``|Возвращает строку.||
|``GetCopy() []rune``|||
|``Length() int``|||
|``SetModeThreadSafe(m bool)``|||
|``GetModeThreadSafe() bool``|||