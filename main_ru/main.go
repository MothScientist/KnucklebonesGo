package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"os/exec"
	"sync"

	"golang.org/x/exp/slices" // костыль для версий до 1.21, т.к. потом модуль slices является встроенным
)

type player struct {
	dice          [][]int // Текущее разложение костей
	points        int     // Сумма очков (текущая)
	name          string  // Имя игрока (вводится в начале игры)
	desk_position bool    // Сверху или снизу отображает в консоли игрока
}

func main(){
	// Точка входа

	player_0 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points: 0,
		name:   "",
		desk_position: false, // 0 -> выводим поле этого игрока первым
	}

	player_1 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points: 0,
		name:   "",
		desk_position: true, // 1
	}

	fmt.Println("Имя 1-го игрока")
	fmt.Scanf("%s\n", &player_0.name)

	fmt.Println("Имя 2-го игрока")
	fmt.Scanf("%s\n", &player_1.name)
	fmt.Print("\n")

	step := true // переменная, определяющая, какой игрок ходит сейчас

	var cur_player player // храним ссылку на игрока, который сейчас делает ход
	var other_player player // ссылка на второго игрока

	var number int // число, которое получаем в каждом раунде

	var new_state int // номер столбца, в который игрок хочет записать число

	var wg sync.WaitGroup // Переменная-счетчик для ожидания горутин


	for { // запускаем бесконечный цикл, из которого выйдем, только если diceIsFull() вернет true
		clearConsole() // Очищаем консоль

		step = !step // передаем ход другому игроку

		// печатаем поля игроков:
		fmt.Println(player_0.name)
		for column := range player_0.dice {
			for row := range player_0.dice[column] {
				state := player_0.dice[row][column]
				if state != 0 {
					fmt.Print(state)
				} else {
					fmt.Print(" ")
				}
				fmt.Print(" ")
			}
			fmt.Print("\n")
		}

		fmt.Print("\n")
		fmt.Printf("Счёт: %d", player_0.points)

		fmt.Print("\n\n\n")

		fmt.Println(player_1.name)
		for column := range player_1.dice {
			for row := range player_1.dice[column] {
				state := player_1.dice[row][2 - column] // "2 - column" нужно, так как мы выводим матрицу второго игрока направлением к матрице первого игрока
				if state != 0 {
					fmt.Print(state)
				} else {
					fmt.Print(" ")
				}
				fmt.Print(" ")
			}
			fmt.Print("\n")
		}

		fmt.Print("\n")
		fmt.Printf("Счёт: %d", player_1.points)

		if step {
			cur_player = player_1
			other_player = player_0
		} else {
			cur_player = player_0
			other_player = player_1
		}

		fmt.Print("\n\n\n")

		// Теперь печатаем имя игрока, который ходит и его новое число:
		fmt.Printf("%s, твой ход!", cur_player.name)
		fmt.Print("\n")
		number = getRandomDice()
		fmt.Printf("Ваше число: %d", number)
		fmt.Print("\n")
		availableColumns := cur_player.getAvailableColumns()
		for {
			fmt.Print("В какую колонку поместить данный кубик [1-3]: ")
			fmt.Fscan(os.Stdin, &new_state)
			new_state -= 1 // Массив начинается с нуля
			if (slices.Contains(availableColumns, new_state)) {
				cur_player.dice[new_state][getIndexToInsertDice(cur_player.dice[new_state])] = number
				break
			} else {
				fmt.Println("Такой столбец не существует, либо уже заполнен")
			}
		}

		// Смотрим, если у пользователя "сбили" кости
		other_player.reCalcDice(new_state, cur_player.dice[new_state])

		// Пересчитываем points для обоих игроков после хода и "выбивания" костей соперника
		// TODO - Почему не работает через указатели?
		for i := 0; i <= 1; i++ {
			wg.Add(1)
			if i == 0 {
				go player_0.calcPoints(&wg)
			} else {
				go player_1.calcPoints(&wg)
			}
		}
		wg.Wait() // Ожидаем выполнения всех горутин (обнуление счетчика)

		// TODO Тут тоже можно сделать горутину и в канал помещать общую переменную
		if player_0.diceIsFull() || player_1.diceIsFull() {
			break // Если у кого-то заполнено поле, то выходим из игры
		}

	}

	clearConsole() // Очищаем консоль перед объявлением победителя

	fmt.Printf("%s - %d\n", player_0.name, player_0.points)
	fmt.Printf("%s - %d\n", player_1.name, player_1.points)

	if player_0.points > player_1.points {
        fmt.Printf("Победитель: %s\n", player_0.name)
    } else if player_1.points > player_0.points {
        fmt.Printf("Победитель: %s\n", player_1.name)
    } else {
        fmt.Println("Ничья!")
    }
}

func (p *player) calcPoints(wg *sync.WaitGroup) {
	defer wg.Done() // после выполнения горутины вычитаем единицу из счетчика

	points := 0

	for column := range p.dice {
		slice_copy := Dedupe(p.dice[column]) // необходим для исключения повторов, чтобы корректно посчитать множители
		fmt.Println("-->>")
		fmt.Println(slice_copy)
		for _, value := range slice_copy {
			factor := count(value, p.dice[column])
			if factor == 3 {
				points += (value * 3) * 3
			} else if factor == 2 {
				points += (value * 2) * 2
			} else {
				points += value
			}
		}
	}

	p.points = points
}

func (p *player) reCalcDice(column int, opponentColumnDice []int) {
	// Пересчитывает dice после хода противника
	for idx, value := range p.dice[column] {
		if slices.Contains(opponentColumnDice, value) {
			p.dice[column][idx] = 0
		}
	}
}

func getRandomDice() int {
	// Получает случайное число 1-6 (включительно)
	// обновляем seed перед использованием - иначе всё свалится к одному числу
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(6) + 1
}

func (p *player) diceIsFull() bool {
	// Определяет, что игра закончена -> у игрока после хода заполнены все поля внутри dice
	for column := range p.dice {
		for row := range p.dice[column] {
			if p.dice[column][row] == 0 {
				return false
			}
		}
	}
	return true
}

func (p *player) getAvailableColumns() []int {
	// Получаем список доступных для записи столбцов
	var res []int // в данный срез записываем доступные столбцы
	for column := range p.dice {
		for row := range p.dice[column] {
			if p.dice[column][row] == 0 {
				res = append(res, column)
				break
			}
		}
	}
	return res
}

func getIndexToInsertDice(dice []int) int {
	// Возвращает позицию для записи в строку([0,2])
	var res int
	for idx, value := range dice {
		if value == 0 {
			res = idx
		}
	}
	return res
}

func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func count(value int, arr []int) int {
	// Возвращает количество вхождения числа в массив
	res := 0

	for _, arr_val := range arr {
		if arr_val == value {
			res += 1
		}
	}

	return res
}

func Dedupe[T comparable](arr []T) []T {
	// Возвращаем slice без повторяющихся элементов
	m, uniq := make(map[T]struct{}), make([]T, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v], uniq = struct{}{}, append(uniq, v)
		}
	}
	return uniq
}