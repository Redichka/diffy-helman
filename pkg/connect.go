package pkg

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const dataBase string = "root:root@tcp(127.0.0.1:3306)/helly-hafman" // константа для подключения к БД(конфиденциальность не так важна, так как это локальное подключение)

type Connect struct {
	Persons []User // создание структуры для хранения людей которые находятся в соединении в массиве юзеров
}

type User struct { // класс юзер, включает в себя юзернейм пользователя в телеграме и айди его чата с ботом
	Username   string
	TelegramID int64
}

func (b *Bot) register(message *tgbotapi.Message) error {
	db, err := sql.Open("mysql", dataBase) // открываем базу данных mySQL
	defer db.Close()                       // функция, которая при завершении функции закроет наше подключение к базе данных
	if err != nil {
		return err
	}
	find, _ := db.Prepare("SELECT `username` FROM `user` WHERE `username` = ?") // создание запроса к БД, который находит пользователя с таким же юзернэйм
	defer find.Close()                                                          // закрытие запроса при выходе из функции
	var usname string
	err = find.QueryRow(message.Chat.UserName).Scan(&usname) // сканируем полученные из БД данные в переменную
	if message.Chat.UserName == usname {                     // проверяем его на совпадение с нашим юзернэймом
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы уже зарегестрированы") // выдаем сообщение о том что мы уже зарегестрированы в случаи совпадения
		b.bot.Send(msg)
		return nil // выходим из функции
	}
	add, _ := db.Prepare("INSERT INTO `user` (`username`, `telegramID`) VALUES (?, ?)") // создаем запрос для добавления пользователя в БД
	defer add.Close()                                                                   // закрытие запроса при выходе из функции
	_, err = add.Exec(message.Chat.UserName, message.Chat.ID)                           // добавляем юзернэйм и айди чата телеграма в базу данных
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно зарегестрировались") // сообщаем пользователю об окончании регистрации
	b.bot.Send(msg)
	return nil
}

func (b *Bot) connect(message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) (*Connect, error) {
	db, _ := sql.Open("mysql", dataBase)                                               // подключение к БД
	defer db.Close()                                                                   // закрытие БД при выходе из функции
	stmt, _ := db.Prepare("SELECT `id`, `username` FROM `user` WHERE `username` != ?") // создание запроса к БД для нахождения всех Айди и юзернеймов в БД кроме нашего
	data, err := stmt.Query(message.Chat.UserName)                                     // выполнение запроса
	stmt.Close()                                                                       // закрытие запроса
	str := "Введите номера пользователей, с которыми вы хотите установить соединенние\n"
	for data.Next() { // создаем цикл, который поочередно будет доставать всех полученных нами пользователей
		var id int
		var username string
		err = data.Scan(&id, &username) // собственно достаем данные
		if err != nil {
			return nil, err
		}
		str += fmt.Sprintf("%d. @%s\n", id, username) // добавляем их к строке
	}
	if err != nil {
		return nil, err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, str) // отправляем сообщение со списком пользователей
	b.bot.Send(msg)
	persons := []User{{Username: message.Chat.UserName, TelegramID: message.Chat.ID}} // создаем экземпляр типа массив юзеров и добавляем туда себя же
	mutex := sync.Mutex{}                                                             // создаем мьютекс для того чтобы закрывать доступ к одним и тем же данным в одно и то же время для горутин
	updateChan := make(chan tgbotapi.Update)                                          // создаем канал для обновлений
	var wg sync.WaitGroup                                                             // создаем закрывашку, которая служит для закрытия выполнения кода пока все горутины не закончат действие
	ctx, cancel := context.WithCancel(context.Background())                           // создаем контекст для того чтобы в нужный момент закрыть горутину
	go func(ctx context.Context) {                                                    // создаем горутину и передаем в нее контекст
		for { // создаем бесконечный цикл
			select { // создаем конструкцию для того, чтобы она либо ждала обновлений от телеграма, либо закрытия контекста
			case update, ok := <-updates: // смотрит приходят ли обновления в канал
				if !ok {
					return
				}
				select { // если они пришли то кидает обновление в другой канал updateChan
				case updateChan <- update:
				}
			case <-ctx.Done(): // Ожидание закрытия контекста, после чего мы закрываем горутину
				return
			}
		}
	}(ctx) // сам переданный контекст
	for update := range updateChan { // создаем цикл, который проверяет все обновления
		if update.Message.Chat.ID == persons[0].TelegramID { // смотрим от нас ли сообщение
			numbers := strings.Fields(update.Message.Text) // если оно от нас, то делим пробелами номера всех пользователей которых мы выбрали
			for _, number := range numbers {               // пробегаемся по массиву этих номеров
				wg.Add(1)                // добавляем счетчик в блок кода wg
				go func(number string) { // запускаем горутину, передавая туда номер пользователя которого мы хотим добавить
					defer wg.Done()                                                                                                                        // при завершении работы горутины опустит счетчик на 1
					find, _ := db.Prepare("SELECT `username`, `telegramID` FROM `user` WHERE `id` = ?")                                                    // создаем запрос для поиска по айди юзернейма и айди телеграм чата
					defer find.Close()                                                                                                                     // закрываем соединение после выхода из горутины
					user := &User{}                                                                                                                        // создаем экземпляр юзера для записи данных туда
					_ = find.QueryRow(number).Scan(&user.Username, &user.TelegramID)                                                                       // записываем полученные из запроса данные
					msg = tgbotapi.NewMessage(user.TelegramID, fmt.Sprintf("С вами хочет установить соединение @%s, вы согласны?", message.Chat.UserName)) // спрашиваем у пользователя разрешение на соединение
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(                                                                                    // создаем кнопки на сообщении с выбором да/нет для пользователя
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
							tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
						),
					)
					b.bot.Send(msg)                   // отправляем сообщение
					for update1 := range updateChan { // смотрим сообщения полученные из канала updateChan
						if update1.CallbackQuery.Message.Chat.ID == user.TelegramID { // проверяем являются ли они нажатием на кнопку
							if update1.CallbackQuery.Data == "yes" { // если ответ да
								mutex.Lock()                                                                                                  // блокируем доступ к массиву персонс
								persons = append(persons, *user)                                                                              //добавляем нашего пользователя
								mutex.Unlock()                                                                                                // открываем доступ
								msg2 := tgbotapi.NewMessage(persons[0].TelegramID, fmt.Sprintf("Соединение с %s установлено", user.Username)) // отправляем обоим сообщения о том что соединение между ними установлено
								b.bot.Send(msg2)
								msg3 := tgbotapi.NewMessage(user.TelegramID, fmt.Sprintf("Соединение с %s установлено", persons[0].Username))
								b.bot.Send(msg3)
								break // выходим из горутины
							} else {
								msg4 := tgbotapi.NewMessage(message.Chat.ID, "В соединении отказано") // если пользователь нажал нет то пишем нашему пользователю что в соединении отказано
								b.bot.Send(msg4)
								break // выходим из горутины
							}
						}
					}
				}(number) // собственно сам переданный номер
			}
		} else {
			continue // в случае если сообшение не от нашего пользователя то игнорируем
		}
		wg.Wait()         // мы останавливаемся тут пока все горутины не закончат свою работу
		close(updateChan) // закрываем канал updateChan за ненадобностью
		cancel()          // закрываем контекст, после чего та горутина выше закончит свою работу и не будет принимать новые обновления
	}
	return &Connect{Persons: persons}, nil // возвращаем структуру коннект со всеми пользователями, которых мы подключили, а так же нил, так как тут нет ошибок
}

func (b *Bot) chat(connect *Connect, updates tgbotapi.UpdatesChannel) error {
	for update := range updates { // создаем цикл для получения обновлений
		if update.CallbackQuery != nil { // игнорируем нажатия на кнопки за ненадобностью
			continue
		}
		if update.Message.IsCommand() { // проверка на то является ли обновление командой
			switch update.Message.Command() { // для выбора между вариантами команд
			case "disconnect": // команда для отключения пользователя
				{
					for index, person := range connect.Persons { // перебираем всех пользователей в соединении
						if update.Message.Chat.ID == person.TelegramID { // находим в массиве человека который хочет прервать соединение
							msg := tgbotapi.NewMessage(person.TelegramID, "Соединение c чатом прервано") // пишем ему что соединение прервано
							b.bot.Send(msg)
							connect.Persons = append(connect.Persons[:index], connect.Persons[index+1:]...) // удаляем этого пользователя из массива
							continue                                                                        // переходим на следующую итерацию цикла, так как мы уже удалили пользователя
						}
						msg := tgbotapi.NewMessage(person.TelegramID, fmt.Sprintf("Соединение c %s прервано", update.Message.Chat.UserName)) // присылаем всем остальным участникам сообщение об отключении этого пользователя
						b.bot.Send(msg)
						break // выходим из цикла
					}
					if len(connect.Persons) == 1 { // если в соединении остался 1 пользователь то разрываем соединение полностью
						msg := tgbotapi.NewMessage(connect.Persons[0].TelegramID, "Соединение полностью прекращено")
						b.bot.Send(msg)
						return nil
					}
				}
			case "sendpm": // отправка личного сообщения пользователю
				{
					commandArguments := strings.Split(update.Message.CommandArguments(), " ") // делим получаемые с командой аргументы
					username := commandArguments[0]                                           // юзернейм того кому мы хотим отправить сообщение передается первым аргументом
					commandArguments = commandArguments[1:]                                   // удаляем юзернейм с массива, чтобы осталось только наше сообщение в массиве
					text := strings.Join(commandArguments, " ")                               // воссоздаем из оставшегося массива текст нашего сообщения
					message := fmt.Sprintf("[pm]%s : %s", update.Message.Chat.UserName, text) // создаем полный экземпляр нашего сообщения, [pm] у никнейма означает то что сообщение личное
					for _, person := range connect.Persons {                                  // перебираем всех пользователей
						if person.Username == username { // находим нужного пользователя
							msg := tgbotapi.NewMessage(person.TelegramID, message) // отправляем ему наше сообщение
							b.bot.Send(msg)
						}
					}
					log.Println(message) // выводим его в логи для имитации тайного считывания третьим лицом
				}
			case "send": // отправка сообщения всем пользователям
				{
					message := fmt.Sprintf("%s : %s", update.Message.Chat.UserName, update.Message.CommandArguments()) // создаем экземпляр сообщения
					for _, person := range connect.Persons {                                                           // пробегаемся по всем пользователям
						if update.Message.Chat.ID == person.TelegramID { // если это мы то пропускаем
							continue
						}
						msg := tgbotapi.NewMessage(person.TelegramID, message) // остальным отправляем наше сообщение
						b.bot.Send(msg)
					}
					log.Println(message) // лог для имитации тайного считывания сообщений третьим лицом
				}
			case "generateNumber":
				{
					numbers := strings.Split(update.Message.CommandArguments(), " ") // делим полученные аргументы команды
					first1, _ := strconv.ParseInt(numbers[0], 10, 64)
					second2, _ := strconv.ParseInt(numbers[1], 10, 64) // конвертируем их в числа размера 64 байта
					third3, _ := strconv.ParseInt(numbers[2], 10, 64)
					code := generateYouNumber(first1, second2, third3)                              // передаем их в функцию для генерации своего уникального номера
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(code, 10)) // высылаеем сообщение, конвертируя полученное число в строку
					b.bot.Send(msg)
				}
			case "decryptSecretNumber":
				{
					numbers := strings.Split(update.Message.CommandArguments(), " ") // делим полученные аргументы команды
					first1, _ := strconv.ParseInt(numbers[0], 10, 64)
					second2, _ := strconv.ParseInt(numbers[1], 10, 64) // конвертируем их в числа размера 64 байта
					third3, _ := strconv.ParseInt(numbers[2], 10, 64)
					code := decryptSecretNumber(first1, second2, third3)                            // передаем их в функцию для дешифрования номеров
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(code, 10)) /// высылаем сообщения, конвертируя полученное число в строку
					b.bot.Send(msg)
				}
			case "encrypt": // шифрование
				{
					arguments := strings.Split(update.Message.CommandArguments(), " ") // делим аргументы
					second2, _ := strconv.ParseInt(arguments[1], 10, 64)               // конвертируем в число полученное значение
					code := encrypt(arguments[0], second2)                             // передаем значения в функцию для шифрования
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, code)           // высылаем полученный результат
					b.bot.Send(msg)
				}
			case "decrypt": // дешифровка
				{
					arguments := strings.Split(update.Message.CommandArguments(), " ") // делим полученные значения
					second2, _ := strconv.ParseInt(arguments[1], 10, 64)               // конвертируем в число полученное значение
					code := decrypt(arguments[0], second2)                             // передаем значеня в функцию для дешифровки
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, code)           // высылаем полученный результат
					b.bot.Send(msg)
				}
			case "generateRandomNumber": // генератор случайного числа
				{
					arguments := strings.Split(update.Message.CommandArguments(), " ") // мы передаем 4 аргумента: первое число, степень первого числа, второе число, степень второго числа, делим их
					numbers := []int{}                                                 // создаем массив для хранения этих чисел
					for _, num := range arguments {                                    // перебираем все 4 числа и конвертируем в тип интеджер, добавляя в массив
						number, _ := strconv.Atoi(num)
						numbers = append(numbers, number)
					}
					min := new(big.Int).Exp(big.NewInt(int64(numbers[0])), big.NewInt(int64(numbers[1])), nil) // возводим первое число в степень второго, так как мы работаем с библиотекой больших чисел то переводим int в int64
					max := new(big.Int).Exp(big.NewInt(int64(numbers[2])), big.NewInt(int64(numbers[3])), nil) // возводим третье число в степень четвертого, так как мы работаем с библиотекой больших чисел то переводим int в int64
					source := rand.New(rand.NewSource(time.Now().UnixNano()))                                  // задаем источником для генерации чисел нынешнее время в наносекундах
					diff := new(big.Int).Sub(max, min)                                                         // вычисляем разницу между максимальным и минимальным числом
					randomDiff := new(big.Int).Rand(source, diff)                                              // генерируем случайное число размером с разницу между мин и макс числами
					randomNumber := new(big.Int).Add(randomDiff, min)                                          // добавляем к получившемуся числу минимальное, чтобы оно входило в диапозон
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, randomNumber.String())                  // отправляем полученный результат пользователю
					b.bot.Send(msg)
				}
			}
		}
	}
	return nil
}
