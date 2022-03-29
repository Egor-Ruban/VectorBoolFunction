package main

import (
	"errors"
	"fmt"
	"math"
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

//newRandomVBF() - генерирует случайную векторную булеву функцию по заданным n и m
func newConstVBF(n, m int, b bool) (VectorBoolFunction, error) {
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
		k := 0
		if b {
			k = 1<<blockSize() - 1
		}
		bf.value[i] = block(k) << bf.wastedBits //генерирует m бит и сдвигает их к старшим разрядам
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
	return 0
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

// WHT - выполняет преобразования Уолша-Адамара. Возвращает m массивов с результатами
func (bf VectorBoolFunction) WHT() [][]int {
	//сначала создание квадратной матрицы и заполнение её 1 и -1, в зависимости от функции
	wht := make([][]int, bf.m)
	for i := range wht {
		wht[i] = make([]int, bf.rows)
		for j := 0; j < bf.rows; j++ {
			if ((bf.value[j] >> (blockSize() - i - 1)) & 1) == 0 {
				wht[i][j] = 1
			} else {
				wht[i][j] = -1
			}
		}
	}

	//дальше преобразование для каждой функции идёт по очереди в цикле
	for m := 0; m < bf.m; m++ {
		//l - через сколько элементов находится тот, с которым мы будем складывать или вычитать
		for l := 1; l < bf.rows; l *= 2 {
			//j указывает на начало блока. Можно было бы и убрать этот цикл, но тогда надо делать сложные условия для k
			for j := 0; j < bf.rows; j += 2 * l {
				//k указывает какой элемент мы сейчас высчитываем
				for k := j; k < j+l; k++ {
					x := wht[m][k]
					y := wht[m][k+l]
					wht[m][k] = x + y
					wht[m][k+l] = x - y
				}
			}
		}
	}

	return wht
}

//cor вычисляет степень корреляционной имунности. Для этого проверяет все наборы с весами от 1 до n и если
//в преобразовании Уэлша-Адамара на наборе есть не только нули, то принимает за результат вес векторов из прошлого набора
func (bf VectorBoolFunction) cor() []int {
	res := make([]int, bf.m)
	for i, v := range bf.WHT() {
		for j := bf.n - 1; j >= 0; j-- {
			b := isCorJ(bf.rows-1, j, 0, v, bf.n)
			if !b {
				break
			}
			res[i] = bf.n - j
		}
	}
	return res
}

//isCorJ функция, возвращающая false,
//если на наборе векторов с весом step в преобразовании Уэлша-Адамара присутствуют не только нули
func isCorJ(initial int, step int, startFrom int, whtVector []int, n int) bool {
	if startFrom+step > n {
		return true
	}
	if step == 0 {
		if whtVector[initial] != 0 {
			return false
		}
	} else {
		for i := startFrom; i < n; i++ {
			if !isCorJ(initial-(1<<i), step-1, i+1, whtVector, n) {
				return false
			}
		}
	}
	return true
}

//Affine - возвращает нелинейность каждой координаты
// и векторную функцию, каждая кордината которой является НАП для соответствующей координаты исходной функции
func (bf VectorBoolFunction) Affine() ([]int, VectorBoolFunction) {
	nonlinearity := make([]int, bf.m)
	maxPlaces := make([]int, bf.m)
	isInversed := make([]int, bf.m)

	//цикл по всем координатам
	for i, v := range bf.WHT() {
		max := 0
		//находим максимальное абсолютное значение в преобразовании Уолша-Адамара.
		//Сохраняем значение, позицию и отрицательное ли оно
		for j, el := range v {
			if math.Abs(float64(el)) > math.Abs(float64(max)) {
				isInversed[i] = 0
				if el < 0 {
					isInversed[i] = 1
				}
				maxPlaces[i] = j
				max = el
			}
		}
		//Когда нашли максимальное - вычисляем нелинейность
		nonlinearity[i] = (1 << (bf.n - 1)) - max/2
		if max < 0 {
			nonlinearity[i] = (1 << (bf.n - 1)) + max/2
		}
	}
	return nonlinearity, bf.getAffineFunction(maxPlaces, isInversed)
}

//getAffineFunction - Получает на вход исходную функцию, НАП координат в числовом виде, а также добавочное слагаемое
// 0, если значение в ПУА было положительным и 1, если было отрицательным для каждой координаты
func (bf VectorBoolFunction) getAffineFunction(coordinated []int, isInversedCoordinates []int) VectorBoolFunction {
	t := VectorBoolFunction{
		value:      make([]block, bf.rows),
		rows:       bf.rows,
		n:          bf.n,
		m:          bf.m,
		wastedBits: blockSize() - bf.m,
	}
	for i, p := range coordinated {
		for j := 0; j < bf.rows; j++ {
			vars := p & j
			res := 0
			for k := 0; k < bf.n; k++ {
				res ^= (vars >> k) & 1
			}
			res ^= isInversedCoordinates[i]
			t.value[j] |= block(res << (blockSize() - i - 1))
		}
	}
	return t
}

//isCoordinatesDegenerate проверят все координаты векторной булевой функции на наличие фиктивных переменных
//Если у координаты есть фиктивные переменные, то она помечается единицей в результирующем векторе
//и нулём, если она существенно зависит от всех своих переменных
func (bf VectorBoolFunction) isCoordinatesDegenerate() block {
	res := block(0)
	anf := bf.Moebius()
	for j := 0; j < bf.m; j++ {
		significantVars := 0
		for i := 1; i < bf.rows; i++ {
			if (anf.value[i] >> (blockSize() - 1 - j) & 1) == 1 {
				significantVars |= i
			}
		}
		if significantVars != bf.rows-1 {
			res |= 1 << (blockSize() - 1 - j)
		}
	}
	return res
}

//isNonDegenerate - показывает является ли данная векторная функция невырожденной
func (bf VectorBoolFunction) isNonDegenerate() bool {
	return bf.isCoordinatesDegenerate() == 0
}

func newFunctionByHand() VectorBoolFunction {
	fmt.Print("n? ")
	n, m := 0, 0
	fmt.Scanln(&n)
	fmt.Print("m? ")
	fmt.Scanln(&m)
	t := VectorBoolFunction{
		value:      make([]block, 1<<n),
		rows:       1 << n,
		n:          n,
		m:          m,
		wastedBits: blockSize() - m,
	}
	for i := 0; i < t.rows; i++ {
		formatString := "value at %0" + strconv.Itoa(n) + "b? "
		fmt.Printf(formatString, i)
		temp := ""
		fmt.Scanln(&temp)
		if len(temp) != m {
			fmt.Println("try again")
			i--
			continue
		}
		for j := 0; j < m; j++ {
			if temp[j] != '0' && temp[j] != '1' {
				fmt.Println("try again")
				i--
				break
			}
			t.value[i] |= block((temp[j]-'0')&1) << (blockSize() - j - 1)
		}
	}
	fmt.Println()
	return t
}

//functionByHand - возволяет создать функцию, вводя её на клавиатуре
func functionByHand() VectorBoolFunction {
	fmt.Print("n ? ")
	n := 0
	fmt.Scanln(&n)
	fmt.Print("m ? ")
	m := 0
	fmt.Scanln(&m)
	t := VectorBoolFunction{
		value:      make([]block, 1<<n),
		rows:       1 << n,
		n:          n,
		m:          m,
		wastedBits: blockSize() - m,
	}

	vars := "%0" + strconv.Itoa(n) + "b ? "
	for i := 0; i < (1 << n); i++ {
		fmt.Printf(vars, i)
		value := ""
		fmt.Scanln(&value)
		if len(value) != m {
			fmt.Println("try again")
			i--
		}
		for j := 0; j < m; j++ {
			if value[j] != '0' && value[j] != '1' {
				fmt.Println("try again")
				i--
				break
			}
			t.value[i] |= block(value[j]-'0') << (32 - 1 - j)
		}
	}

	fmt.Println()

	return t
}
