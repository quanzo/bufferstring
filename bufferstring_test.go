package bufferstring

import (
	//	"fmt"
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/quanzo/gservice/bufferint"
)

type TestingInsert struct {
	name   string
	pos    int
	source string
	insert string
	result string
}
type TestingDelete struct {
	name   string
	pos    int
	count  int
	source string
	result string
}
type TestingReplace struct {
	name    string
	pos     int
	count   int    // кол-во заменяемых символов
	replace string // строка для замены
	source  string
	result  string
}
type TestingFindReplace struct {
	name    string
	find    []string
	replace []string
	source  string
	result  string
}
type TestingFindMask struct {
	source    string
	mask      string
	res_index int
	res_str   string
}

type TestingFindFirst struct {
	source    string
	find      []string
	res       int
	res_index int
}

type TestingPop struct {
	source       string
	start, count int
	result       string
	dest         string
}

type TestingWalk struct {
	source       string
	start, count int
	result       string
	walkFunc     func(index int, value *rune)
}

type TestingFilter struct {
	source     string
	result     string
	filterFunc func(index int, value rune) bool
}

var tInsert = []TestingInsert{
	TestingInsert{"вставка 10 чисел в начало строки", 0, "Hello! Привет!", "1234567890", "1234567890Hello! Привет!"},
	TestingInsert{"вставка 10 чисел в конец строки", 24, "1234567890Hello! Привет!", "1234567890", "1234567890Hello! Привет!1234567890"},
	TestingInsert{"вставка === в середину строки", 10, "1234567890Hello! Привет!0123456789", "===", "1234567890===Hello! Привет!0123456789"},
}

var tDelete = []TestingDelete{
	TestingDelete{"Del 2 symb in 2 pos", 2, 2, "Hello! Привет!", "Heo! Привет!"},
	TestingDelete{"Del symbols in endline", 8, 10, "Heo! Привет!", "Heo! При"},
}

var tReplace = []TestingReplace{
	TestingReplace{"", 0, 0, "abcde", "0123456789 0123456789 0123456789", "abcde0123456789 0123456789 0123456789"},
	TestingReplace{"Del 5 sym from begining line", 0, 5, "", "abcde0123456789 0123456789 0123456789", "0123456789 0123456789 0123456789"},
	TestingReplace{"", 10, 0, "abcde", "0123456789 0123456789 0123456789", "0123456789abcde 0123456789 0123456789"},
	TestingReplace{"", 10, 5, "qwerty", "0123456789abcde 0123456789 0123456789", "0123456789qwerty 0123456789 0123456789"},
	TestingReplace{"Delete substring", 10, 6, "", "0123456789qwerty 0123456789 0123456789", "0123456789 0123456789 0123456789"},
	TestingReplace{"", 32, 5, "qwerty", "0123456789 0123456789 0123456789", "0123456789 0123456789 0123456789qwerty"},
	TestingReplace{"", 32, 10, "", "0123456789 0123456789 0123456789qwerty", "0123456789 0123456789 0123456789"},
	TestingReplace{"Trim line", 10, 100, "", "0123456789 0123456789 0123456789qwerty", "0123456789"},
}

var tFindReplace = []TestingFindReplace{
	TestingFindReplace{"", []string{"o", "i"}, []string{"-", "_"}, "Lorem Ipsum is simply dummy text Lorem Ipsum has been the industry desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.", "L-rem Ipsum _s s_mply dummy text L-rem Ipsum has been the _ndustry deskt-p publ_sh_ng s-ftware l_ke Aldus PageMaker _nclud_ng vers_-ns -f L-rem Ipsum."},
	TestingFindReplace{"", []string{"o", "i"}, []string{"---", "___"}, "Lorem Ipsum is simply dummy text Lorem Ipsum has been the industry desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.", "L---rem Ipsum ___s s___mply dummy text L---rem Ipsum has been the ___ndustry deskt---p publ___sh___ng s---ftware l___ke Aldus PageMaker ___nclud___ng vers___---ns ---f L---rem Ipsum."},
	TestingFindReplace{"", []string{"a"}, []string{"--a--"}, "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc,", "Lorem ipsum dolor sit --a--met, consectetuer --a--dipiscing elit. Aene--a--n commodo ligul--a-- eget dolor. Aene--a--n m--a--ss--a--. Cum sociis n--a--toque pen--a--tibus et m--a--gnis dis p--a--rturient montes, n--a--scetur ridiculus mus. Donec qu--a--m felis, ultricies nec, pellentesque eu, pretium quis, sem. Null--a-- consequ--a--t m--a--ss--a-- quis enim. Donec pede justo, fringill--a-- vel, --a--liquet nec, vulput--a--te eget, --a--rcu. In enim justo, rhoncus ut, imperdiet --a--, venen--a--tis vit--a--e, justo. Null--a--m dictum felis eu pede mollis pretium. Integer tincidunt. Cr--a--s d--a--pibus. Viv--a--mus elementum semper nisi. Aene--a--n vulput--a--te eleifend tellus. Aene--a--n leo ligul--a--, porttitor eu, consequ--a--t vit--a--e, eleifend --a--c, enim. Aliqu--a--m lorem --a--nte, d--a--pibus in, viverr--a-- quis, feugi--a--t --a--, tellus. Ph--a--sellus viverr--a-- null--a-- ut metus v--a--rius l--a--oreet. Quisque rutrum. Aene--a--n imperdiet. Eti--a--m ultricies nisi vel --a--ugue. Cur--a--bitur ull--a--mcorper ultricies nisi. N--a--m eget dui. Eti--a--m rhoncus. M--a--ecen--a--s tempus, tellus eget condimentum rhoncus, sem qu--a--m semper libero, sit --a--met --a--dipiscing sem neque sed ipsum. N--a--m qu--a--m nunc, bl--a--ndit vel, luctus pulvin--a--r, hendrerit id, lorem. M--a--ecen--a--s nec odio et --a--nte tincidunt tempus. Donec vit--a--e s--a--pien ut libero venen--a--tis f--a--ucibus. Null--a--m quis --a--nte. Eti--a--m sit --a--met orci eget eros f--a--ucibus tincidunt. Duis leo. Sed fringill--a-- m--a--uris sit --a--met nibh. Donec sod--a--les s--a--gittis m--a--gn--a--. Sed consequ--a--t, leo eget bibendum sod--a--les, --a--ugue velit cursus nunc,"},
	TestingFindReplace{"", []string{"A", "I", "   ", "  ", ".", ","}, []string{"a", "iiiii", " ", " ", ""}, "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque   penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis   vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus.  Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque  rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies  nisi. Nam eget dui.      Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero  venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo.     Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget  bibendum sodales, augue velit cursus nunc,", "Lorem ipsum dolor sit amet consectetuer adipiscing elit aenean commodo ligula eget dolor aenean massa Cum sociis natoque penatibus et magnis dis parturient montes nascetur ridiculus mus Donec quam felis ultricies nec pellentesque eu pretium quis sem Nulla consequat massa quis enim Donec pede justo fringilla vel aliquet nec vulputate eget arcu iiiiin enim justo rhoncus ut imperdiet a venenatis vitae justo Nullam dictum felis eu pede mollis pretium iiiiinteger tincidunt Cras dapibus Vivamus elementum semper nisi aenean vulputate eleifend tellus aenean leo ligula porttitor eu consequat vitae eleifend ac enim aliquam lorem ante dapibus in viverra quis feugiat a tellus Phasellus viverra nulla ut metus varius laoreet Quisque rutrum aenean imperdiet Etiam ultricies nisi vel augue Curabitur ullamcorper ultricies nisi Nam eget dui Etiam rhoncus Maecenas tempus tellus eget condimentum rhoncus sem quam semper libero sit amet adipiscing sem neque sed ipsum Nam quam nunc blandit vel luctus pulvinar hendrerit id lorem Maecenas nec odio et ante tincidunt tempus Donec vitae sapien ut libero venenatis faucibus Nullam quis ante Etiam sit amet orci eget eros faucibus tincidunt Duis leo  Sed fringilla mauris sit amet nibh Donec sodales sagittis magna Sed consequat leo eget bibendum sodales augue velit cursus nunc"},
	//TestingFindReplace{"", []string{}, []string{}, "", ""},
}
var tFindMask []TestingFindMask = []TestingFindMask{
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789", "1*7", 19, "1234567"},
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789", "*w*d", 0, "Hello world"},
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789", "w*\\*", 6, "world om*"},
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789", "H???o", 0, "Hello"},
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789", "0??", 18, "012"},
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789 pix  ", "p??*9", -1, ""}, //49
	TestingFindMask{"Hello world om*ga 012345678901234567890123456789 pix  ", "p??*", 49, "pix  "},
}

var tFindFirst []TestingFindFirst = []TestingFindFirst{
	TestingFindFirst{"Hello world., | $ Good morning ", []string{"$", "|", ".", ","}, 11, 2},
	TestingFindFirst{" ololololol", []string{"ol", "lolo", "olo"}, 1, 0},
	TestingFindFirst{" ololololol", []string{"lolo", "ol", "olo"}, 1, 1},
}

var tPop []TestingPop = []TestingPop{
	TestingPop{"Hello world", 0, 5, "Hello", " world"},
	TestingPop{"Hello world", 2, 2, "ll", "Heo world"},
	TestingPop{"Hello world", 11, 5, "", "Hello world"},
}

var tWalk []TestingWalk = []TestingWalk{
	TestingWalk{
		"Hello world", 0, 100, "Hekko workd",
		func(i int, v *rune) {
			if *v == 'l' {
				*v = 'k'
			}
		},
	},
}

var tFilter []TestingFilter = []TestingFilter{
	TestingFilter{
		"0123456789", "012345",
		func(i int, v rune) bool {
			if i > 5 {
				return false
			}
			return true
		},
	},
}

var (
	defStr   string = "Hello! Привет!"
	str_0    string = "Heo! Привет!"
	str_1    string = "Hello! Приве"
	defStr_2 string = "Hello! Привет! Hello! Привет!"
)

func TestInsert(t *testing.T) {
	for _, v := range tInsert {
		buff := NewFromString(v.source, 1)
		buff.Insert(v.insert, v.pos)
		if buff.String() != v.result {
			t.Error("Error test: ", v.name)
			t.Error("Result:", buff.String(), "Etalon:", v.result)
		}
	}
}
func TestCreateBuffer(t *testing.T) {
	buff := NewFromString(defStr, 0)
	if buff.String() != defStr {
		t.Error("The original and the resulting line is not the same.")
	}
	if buff.Length() != utf8.RuneCountInString(defStr) {
		t.Error("The length of the string and the length of the result is not the same.")
	}
}

func TestDelete(t *testing.T) {
	for _, v := range tDelete {
		buff := NewFromString(v.source, 1)
		buff.Delete(v.pos, v.count)
		if buff.String() != v.result {
			t.Error("Error test: ", v.name)
			t.Error("Result:", buff.String(), "Etalon:", v.result)
		}
	}
}

func TestFindReplace(t *testing.T) {
	for _, v := range tFindReplace {
		buff := NewFromString(v.source, 1)
		buff.FindReplace(v.find, v.replace)
		if buff.String() != v.result {
			t.Error("Error test: ", v.name)
			t.Error("Result:", buff.String(), "Etalon:", v.result)
		}
	}
}

func TestSubstr_1(t *testing.T) {
	t.Log("Получение подстроки из центра")
	buff := NewFromString(defStr, 10)
	if buff.Substr(1, 4) != "ello" {
		t.Error("Подстрока не правильная <" + buff.Substr(1, 4) + ">")
	}
	t.Log("Получение подстроки с конца")
	if buff.Substr(12, 5) != "т!" {
		t.Log([]byte(buff.Substr(12, 5)))
		t.Error("Подстрока не правильная <" + buff.Substr(12, 5) + ">")
	}
}

func TestReplace(t *testing.T) {
	for _, v := range tReplace {
		buff := NewFromString(v.source, 1)
		buff.Replace(v.pos, v.count, v.replace)
		if buff.String() != v.result {
			t.Error("Error test: ", v.name)
			t.Error("Result:", buff.String(), "Etalon:", v.result)
		}
	}
}

func TestStrpos(t *testing.T) {
	buff := NewFromString("0123456789 0123456789 0123456789", 3)
	//***
	t.Log("Find 012 the beginning of the line")
	if buff.Find("012", -1) != 0 {
		t.Error("Error find 012")
	}
	//***
	t.Log("Find 012 the ending of the line")
	if buff.FindReverse("012", -1) != 22 {
		t.Error("Error find 012")
	}
	//***
	t.Log("Find 89 the ending of the line")
	t.Log(buff.Length())
	if buff.FindReverse("89", -1) != 30 {
		t.Error("Error find 89")
	}
	//***
	t.Log("Find 012 from the position 22")
	if buff.FindReverse("012", 22) != 11 {
		t.Error("Error find")
	}
}

func TestAppend(t *testing.T) {
	buff := NewFromString("", 1)
	buff.AppendString("123", "456", "789", "===", "+++++")
	if buff.String() != "123456789===+++++" {
		t.Error("Error append data.")
		t.Log(buff.String())
	}
}
func TestAppend2(t *testing.T) {
	buff := NewFromString("000", 10)
	buff.AppendString("123", "456", "789", "===", "+++++")
	if buff.String() != "000123456789===+++++" {
		t.Error("Error append data.")
		t.Log(buff.String())
	}
}
func TestFindAll(t *testing.T) {
	buff := NewFromString("Lorem Ipsum is simply dummy text Lorem Ipsum has been the industry desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.", 10)
	res := bufferint.New(10, 10)
	f := []rune("Lorem")
	buff.findAll(&f, res)
	if !res.Equal(res.GetCopy(), []int{0, 33, 138}) {
		t.Error("Error find 'Lorem'")
	}
	res.Empty()
	f = []rune("Ipsum")
	buff.findAll(&f, res)
	if !res.Equal(res.GetCopy(), []int{6, 39, 144}) {
		t.Error("Error find 'Ipsum'")
	}
	res.Empty()
	f = []rune("o")
	buff.findAll(&f, res)
	if !res.Equal(res.GetCopy(), []int{1, 34, 72, 87, 131, 135, 139}) {
		t.Error("Error find 'o'")
		t.Log(res.GetCopy())
	}
}

func TestAppendBuffer(t *testing.T) {
	buff := NewFromString("012345678901234567890123456789", 3)
	buff2 := NewFromString("qweqweqweqweqweqwe", 10)
	c := buff.AppendBuffer(buff2, func(i int, v rune) bool {
		if v == 'w' {
			return false
		} else {
			return true
		}
	})
	if c != 12 {
		t.Error("Error count append rune.", c)
	}
	if buff.String() != "012345678901234567890123456789qeqeqeqeqeqe" {
		t.Error("Error result string.")
	}
}

func TestFindMask(t *testing.T) {
	for i, v := range tFindMask {
		buff := NewFromString(v.source, 3)
		p, s := buff.FindMask(v.mask, 0)
		if p != v.res_index || s != v.res_str {
			t.Error("Error test #", i, p, s)
		}
		t.Log(p, s)
	}
}

func TestFindFirst(t *testing.T) {
	for i, v := range tFindFirst {
		buff := NewFromString(v.source, 3)
		p, s := buff.findFirst(&v.find, 0)
		if p != v.res || s != v.res_index {
			t.Error("Error test #", i, p, s)
		}
		t.Log(p, s)
	}
}

func TestPop(t *testing.T) {
	for i, v := range tPop {
		buff := NewFromString(v.source, 3)
		r := buff.Pop(v.start, v.count)
		if r != v.result || buff.String() != v.dest {
			t.Error("Error test #", i, "|", r, "|", buff.String())
		}
	}
}

func TestWalk(t *testing.T) {
	for testNum, test := range tWalk {
		buff := NewFromString(test.source, 3)
		buff.Walk(test.start, test.count, test.walkFunc)
		if test.result != buff.String() {
			t.Error("Error test #", testNum, test.result, " != ", buff.String())
		}
	}
}

func TestFilter(t *testing.T) {
	for testNum, test := range tFilter {
		buff := NewFromString(test.source, 3)
		buff.Filter(test.filterFunc)
		if test.result != buff.String() {
			t.Error("Error test #", testNum, test.result, " != ", buff.String())
		}
	}
}

// BENCHMARK

func BenchmarkSimpleAppendString(b *testing.B) {
	res := ""
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res += "a"
	}
	b.StopTimer()
}

func BenchmarkBufferAppendString(b *testing.B) {
	buff := NewFromString("", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buff.AppendString("a")
	}
	b.StopTimer()
}

func BenchmarkBufferInsertString(b *testing.B) {
	buff := NewFromString("", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buff.Insert("a", buff.Length())
	}
	b.StopTimer()
}

func BenchmarkAppendByteBuffer(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")
	}
	b.StopTimer()
}

func BenchmarkIncreaseBuffSizeMake(b *testing.B) {
	var bf []rune
	for n := 0; n < b.N; n++ {
		bf = make([]rune, n)
	}
	b.StopTimer()
	_ = (len(bf))
}

func BenchmarkIncreaseBuffSizeAppend(b *testing.B) {
	var bf []rune = make([]rune, 0)
	for n := 0; n < b.N; n++ {
		bf = append(bf, make([]rune, 1)...)
	}
	b.StopTimer()
	_ = (len(bf))
}

func BenchmarkConvertStr2Byte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = []byte("Строка")
	}
	b.StopTimer()
}

func BenchmarkConvertStr2Rune(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = []rune("Строка")
	}
	b.StopTimer()
}
