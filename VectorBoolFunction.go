package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

//block - целочисленное значение, в котормо хранится значение функции.
//может быть длиной 8, 16 или 32 бита (64 тоже можно, но генерация обратимой n*m функции сломается)
//значения хранятся "слева направо". Значение нулевой (первой) функции хранится в старшем разряде
type block uint32

//blockSize() - функция, возвращающая длину блока в битах
func blockSize() int {
	return int(unsafe.Sizeof(block(0)) * 8)
}

//VectorBoolFunction - структура, хранящая векторную булеву функцию
type VectorBoolFunction struct {
	//value - массив значений функции
	value []block
	n     int
	m     int
	rows  int

	wastedBits int //wastedBits - сколько в блоке незначащих бит в конце
}

//getWeight() - подсчитывает вес каждой кординаты. Возвращает массив с полученными значениями
func (bf VectorBoolFunction) getWeight() []int {
	weight := make([]int, bf.m)
	for _, v := range bf.value {
		for i := 0; i < bf.m; i++ {
			weight[i] += int(v>>(blockSize()-1-i)) & 1
		}
	}
	return weight
}

//Реализация Интерфейса stringer. Аналог перегрузки оператора вывода
//Представляет функцию в виде таблицы с указанием значений переменных
func (bf VectorBoolFunction) String() string {
	res := ""
	formatString2 := "%0" + strconv.Itoa(bf.n) + "b : "                     //Форматная строка для таблицы переменных
	formatString3 := "%0" + strconv.Itoa(blockSize()-bf.wastedBits) + "b\n" //Форматная строка для значений функции

	for i, v := range bf.value {
		res += fmt.Sprintf(formatString2, i) // строку можно убрать, если не нужен вывод значений переменных
		res += fmt.Sprintf(formatString3, v>>bf.wastedBits)
	}
	return res
}

func (bf VectorBoolFunction) valueOf(i int) string {
	formatString3 := "%0" + strconv.Itoa(blockSize()-bf.wastedBits) + "b\n" //Форматная строка для значений функции
	return fmt.Sprintf(formatString3, bf.value[i]>>bf.wastedBits)
}

//newRandomVBF() - генерирует случайную векторную булеву функцию по заданным n и m
func newRandomVBF(n, m int) (VectorBoolFunction, error) {
	rand.Seed(time.Now().UnixNano())

	//если n или m больше размера блока, то возвращает ошибку
	if n > blockSize() || m > blockSize() {
		return VectorBoolFunction{}, errors.New("n or m is too big")
	}

	//создание объекта функции
	bf := VectorBoolFunction{
		value:      make([]block, 1<<n), //создаём массив на 2^n элементов
		rows:       1 << n,
		n:          n,
		m:          m,
		wastedBits: blockSize() - m,
	}

	for i := 0; i < bf.rows; i++ {
		bf.value[i] = block(rand.Intn(1<<m)) << bf.wastedBits //генерирует m бит и сдвигает их к старшим разрядам
	}
	return bf, nil
}

//newRevVBF() - генерирует случайную обратимую векторную булеву функкцию по заданным n и m
func newRevVBF(n, m int) (VectorBoolFunction, error) {
	rand.Seed(time.Now().UnixNano())

	//если n или m больше размера блока, то возвращает ошибку
	if n > blockSize() || m > blockSize() {
		return VectorBoolFunction{}, errors.New("n or m is too big")
	}

	//создание объекта функции
	bf := VectorBoolFunction{
		value:      make([]block, 1<<n), //создаём массив на 2^n элементов
		rows:       1 << n,
		n:          n,
		m:          m,
		wastedBits: blockSize() - m,
	}

	if n < m {
		//Если n<m, то генерируем m случайных неповторяющихся значений и сохраняем
		l := ((uint64(1) << m) + 64 - 1) / 64 //[2^m/64] - сколько блоков по 64 надо для 2^m бит
		isTaken := make([]uint64, l)          //вектор с флагами, указывающими на то, какие значения уже были использованы

		for i := 0; i < bf.rows; i++ {
			x := rand.Intn(1 << m)
			//проверка был ли x использован раньше
			if (isTaken[x/64]>>(64-x%64-1))&1 == 1 {
				i--
			} else {
				bf.value[i] = block(x << bf.wastedBits)
				isTaken[x/64] |= uint64(1) << (64 - x%64 - 1)
			}
		}
		return bf, nil
	} else if m == n {
		//Если m == n, то сначала таблица заполняется значениями из множества {0,1)^n по порядку
		//а затем выполняется случайная перестановка
		for i := 0; i < bf.rows; i++ {
			bf.value[i] = block(i) << bf.wastedBits
		}
		for i := bf.rows - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			t := bf.value[i]
			bf.value[i] = bf.value[j]
			bf.value[j] = t
		}
		return bf, nil
	}
	return VectorBoolFunction{}, errors.New("n > m")
}

//shiftDown() - метод, сдвигающий значения функции "Вниз" по таблице
func (bf VectorBoolFunction) shiftDown(k int) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := bf.rows - 1; i >= k; i-- {
		t.value[i] = bf.value[i-k]
	}
	return t
}

//shiftUp() - метод, сдвигающий значения функции "Вверх" по таблице
func (bf VectorBoolFunction) shiftUp(k int) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows-k; i++ {
		t.value[i] = bf.value[i+k]
	}
	return t
}

//xor() - покомпонентное сложение по модулю два двух таблиц значений ВБФ
func (bf VectorBoolFunction) xor(bf2 VectorBoolFunction) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] ^ bf2.value[i]
	}
	return t
}

//and() - покомпонентное умножение двух таблиц значений ВБФ
func (bf VectorBoolFunction) and(bf2 VectorBoolFunction) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] & bf2.value[i]
	}
	return t
}

//and() - покомпонентное сложение двух таблиц значений ВБФ
func (bf VectorBoolFunction) or(bf2 VectorBoolFunction) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := 0; i < bf.rows; i++ {
		t.value[i] = bf.value[i] | bf2.value[i]
	}
	return t
}

// Moebius - метод, выполняющий преобразование мёбиуса над заданной функцией.
// возвращает полученную после преобразования функцию
func (bf VectorBoolFunction) Moebius() VectorBoolFunction {
	anf := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}

	for i := range bf.value {
		anf.value[i] = bf.value[i]
	}

	//doWeNeedToAdd показывает надо ли нам складывать значения или нет
	doWeNeedToAdd := uint8(255)
	for i := 0; i < anf.n; i++ {
		for j := 0; j < anf.rows; j++ {
			if j%(1<<i) == 0 {
				doWeNeedToAdd ^= 255
			}
			if doWeNeedToAdd > 0 {
				anf.value[j] ^= anf.value[j-(1<<i)]
			}

		}
	}
	return anf
}

// printANF() - возвращает строку, в которой перечислены АНФ для каждой координаты ВБФ
func (bf VectorBoolFunction) printANF() string {
	//Если переменных больше 26, то в английском языке не хватит букв.
	//Хотя смысла от такого представления нет уже при n = 26
	if bf.n > 26 {
		return "too many variables in ANF"
	}

	//создаём массив строк, по строке для каждой координаты
	ANFs := make([]string, bf.m)

	//Для каждой координаты првоеряем. нужна ли единица в сумме
	for j := 0; j < bf.m; j++ {
		if (bf.value[0]>>(blockSize()-j-1))&1 == 1 {
			ANFs[j] += "1"
		}
	}

	//Идём по значениями функций
	for i := 1; i < bf.rows; i++ {
		//Для каждой координаты првоеряем, нужно ли слагаемое на данном значении переменных
		for j := 0; j < bf.m; j++ {
			if (bf.value[i]>>(blockSize()-j-1))&1 == 1 {
				if len(ANFs[j]) != 0 {
					ANFs[j] += "+"
				}
				for k := 0; k < bf.n; k++ {
					if (i>>k)&1 == 1 {
						ANFs[j] += string(rune('z' - k))
					}
				}
			}
		}
	}

	res := ""
	for i := range ANFs {
		res += "anf" + strconv.Itoa(i) + " = "
		if len(ANFs[i]) == 0 {
			res += "0" //Если координата состоит из нулей, то её АНФ = 0
		} else {
			res += ANFs[i]
		}
		res += "\n"
	}

	return res
}

//isEqual() - Проверяет две функции на равенство
func (bf VectorBoolFunction) isEqual(bf2 VectorBoolFunction) bool {
	if bf.n != bf2.n || bf.m != bf2.m {
		return false
	}
	for i := range bf.value {
		if bf.value[i] != bf2.value[i] {
			return false
		}
	}
	return true
}

// degree() - метод, возвращающий степень нелинейности функции.
// Задача выполняется в цикле по кличеству нулей среди значений переменных
// Если при k нулей среди значений функции есть ненулевое, то степень равна n-k
func (bf VectorBoolFunction) degree() int {
	m := bf.Moebius()
	for i := 0; i <= bf.n; i++ {
		b := m.isNotNull(bf.rows-1, i, 0)
		if b {
			return bf.n - i
		}
	}
	return -1
}

// isNotNull() - метод, получающий на вход вектор переменных и информацию о том, сколько единиц из него нужно убрать
// Возвращает true, если после удаления единиц нашлось значение функции не равное нулю
func (bf VectorBoolFunction) isNotNull(initial, step, startFrom int) bool {
	if startFrom+step > bf.n {
		return false
	}
	if step == 0 {
		if bf.value[initial] > 0 {
			return true
		}
	} else {
		for i := startFrom; i < bf.n; i++ {
			if bf.isNotNull(initial-(1<<i), step-1, i+1) {
				return true
			}
		}
	}
	return false
}
