package db

// # Запрос получает пользователя по Telegram id.
//
// * В качестве аргумента принимается телеграм id пользователя.
//
// Возвращает все поля указанного пользователя.
const queryGetUsersByTid = `SELECT tid, name, user_name, first_name, last_name, context, expected_message_type, date_of_registration FROM User_info WHERE tid = $1 LIMIT 1`

// # Запрос получает контексты пользователя.
//
// * В качестве аргумента принимается телеграм id пользователя.
//
// Возвращает id и name контекстов указанного пользователя.
const queryGetContextByTid = `SELECT Contexts.id, Contexts.name, Contexts.text, User_info.context = Contexts.id as selected, Contexts."dateOfCreation" FROM Contexts JOIN User_info ON User_info.id = Contexts."userId" WHERE User_info.tid = $1 ORDER BY "dateOfCreation" DESC`

// # Запрос получает контекст по его id.
//
// * В качестве аргумента принимается id контекста.
//
// Возвращает id и name контекстов указанного пользователя.
const queryGetContextByID = `SELECT Contexts.id, Contexts.name, Contexts.text, true as selected, Contexts."dateOfCreation" FROM Contexts WHERE id = $1`

// # Запрос записывает новый контекст в базу данных.
//
// * Первым аргументом принимает телеграм id пользователя, кому пренадлежыит контекст.
//
// * Вторым аргументом принимает название контекста.
//
// Возвращает id контакста записанного в базу данных.
const querySetNewContext = `INSERT INTO Contexts ( "userId", name ) VALUES ( (SELECT id FROM User_info WHERE tid = $1 LIMIT 1), $2 ) RETURNING id`

// # Запрос записывает нового пользователя в базу данных.
//
// * Первым аргументом принимает телеграм id пользователя.
//
// * Вторым аргументом принимает внутренее имя пользователя.
//
// * Третим аргументом принимает ник в телеграме @nick.
//
// * Четвертым аргументом принимает имя пользователя.
//
// * Пятым аргументом принимает фамилию пользователя.
const querySetNewUser = `INSERT INTO User_info ( tid, name, user_name, first_name, last_name ) VALUES ( $1, $2, $3, $4, $5 )`

// # Обновляет информацию о контексте пользователя.
//
// * Первым аргументом принимает телеграм id пользователя.
//
// * Вторым аргументом принимает новый id контекста.
const queryUpdateUserContext = "UPDATE User_info SET context = $2 WHERE tid = $1"

// # Изменяем тип ожидаемого сообщения от пользователя.
//
// * Первым аргументом принимает телеграм id пользователя.
//
// * Вторым аргументом принимает номер ожидаемого сообщения.
const queryUpdateExpectedMessageType = "UPDATE User_info SET expected_message_type = $2 WHERE tid = $1"

// # Обновляет информацию текста контекста.
//
// * Первым аргументом принимает id контекста.
//
// * Вторым аргументом принимает новый текст контекста.
const queryUpdateContextText = "UPDATE Contexts SET text = $2 WHERE id = $1"

// # Обновляет информацию текста контекста.
//
// * Первым аргументом принимает id контекста.
//
// * Вторым аргументом принимает новое название контекста.
const queryRenameContextByID = "UPDATE Contexts SET name = $2 WHERE id = $1"

// # Удаляет контекст.
//
// * Аргументом принимает id удаляемого контекста.
const queryRemoveContextByID = "DELETE FROM Contexts WHERE id = $1"

// # Запрос ставит на ожидание регистрации нового пользователя.
//
// * Первым аргументом принимает внутренее имя нового пользователя.
//
// * Вторым аргументом клуч доступа.
//
// * Третим аргументом bool является ли ключ паролем.
//
// * Четвертым аргументом TID пользователя создавшего запись.
const queryAddWaitingUser = `INSERT INTO Waiting ( user_name, key, by_password, id_admin ) VALUES ( $1, $2, $3, (SELECT id FROM User_info WHERE tid = $4) )`

// # Запрос получает пользователя из таблицы ожидания регистрации.
//
// * В качестве аргумента принимается ключ доступа.
//
// * Вторым аргументом принимается bool является ли ключь паролем.
//
// Возвращает внутренее имя пользователя.
const queryGetWaitingUser = `SELECT user_name FROM Waiting WHERE key = $1 AND by_password = $2`
