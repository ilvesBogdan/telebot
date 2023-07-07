package db

import (
	"database/sql"
	"fmt"
	"log"
)

// Получает данные о пользователе по ID.
//
// UserId - ID пользователя, который требуется проверить.
//
// Возвращает структуру с данными пользователя.
func GetUserByTelegramID(UserId int64) (user User) {
	err := DB.Get(&user, queryGetUsersByTid, UserId)
	if nil != err && sql.ErrNoRows != err {
		log.Panic("Неудалось получить данные пользователя из БД: ", err)
	}
	return
}

func GetContextsByTid(id int64) ([]Context, error) {
	var context []Context
	err := DB.Select(&context, queryGetContextByTid, id)
	return context, err
}

func GetContextsByID(ID *int) (context Context) {
	err := DB.Get(&context, queryGetContextByID, ID)
	if nil != err && sql.ErrNoRows != err {
		log.Panicf("Неудалось контекст по Context.id: '%v', из БД: %v", *ID, err)
	}
	return
}

// SetNewContextфункция предназначена для установки нового контекста для пользователя.
// Функция принимает строку contextName, содержащую имя нового контекста, и int64 UserId - ID пользователя.
func SetNewContext(contextName string, UserId int64) string {
	contexts, err := GetContextsByTid(UserId)
	if nil != err {
		log.Printf("Неудалось получить список контекстов пользователя id:'%v': %v", UserId, err)
		return "Ошибка базы данных."
	}

	if 160 < len(contextName) {
		return "Слишком длинное название для контекста."
	}

	for _, ct := range contexts {
		ct.TrimSpace()
		if ct.Name == contextName {
			return fmt.Sprintf("Контекст \"%v\" уже существует.", ct.Name)
		}
	}

	var newContextId int

	err = DB.Get(&newContextId, querySetNewContext, UserId, contextName)
	if nil != err {
		log.Printf("Неудалось записать новый контекст пользователя id:'%v': %v", UserId, err)
		return "Ошибка базы данных."
	}

	_, err = DB.Exec(queryUpdateUserContext, UserId, newContextId)
	if nil != err {
		log.Printf("Неудалось обновить новый контекст для пользователя id:'%v': %v", UserId, err)
		return "Ошибка базы данных."
	}

	return fmt.Sprintf("Текущий контекст:\n**%v**", contextName)
}

func SetUserContext(contextID int, userID int64) bool {
	result, err := DB.Exec(queryUpdateUserContext, userID, contextID)
	if nil != err {
		log.Printf("Неудалось обновить контекст для пользователя id:'%v': %v", userID, err)
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

func SetExpectedMessageType(typeMsg int, userID *int64) bool {
	result, err := DB.Exec(queryUpdateExpectedMessageType, *userID, typeMsg)
	if nil != err {
		log.Printf("Неудалось обновить контекст для пользователя id:'%v': %v", *userID, err)
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

func SetNewContextText(contextID int, text string) {
	_, err := DB.Exec(queryUpdateContextText, contextID, text)
	if nil != err {
		log.Printf("Неудалось обновить текст контекста id: '%v': %v", contextID, err)
	}
}

func SetNewUser(nameUser, tlgrmNick, firstName, lastName string, userID int64) {
	_, err := DB.Exec(querySetNewUser, userID, nameUser, tlgrmNick, firstName, lastName)
	if nil != err {
		log.Println("Неудалось зарегистрировать новог опользователя: ")
	}
}

func SetNewMessage(contextID int, text string) (err error) {
	// todo реализовать запись сообщений в БД
	// _, err = database.Exec(queryUpdateContextText, contextID, text)
	return
}

// SetUserContext позволяет установить контекст для пользователя с указанным UserId.
// Функция ищет среди существующих контекстов, привязанных к UserId, контекст с указанным contextName.
// Если такой контекст найден, функция обновляет Users.context, устанавливая Id этого контекста.
func SetUserContextByName(contextName string, userID int64) {
	contexts, err := GetContextsByTid(userID)
	if nil != err {
		log.Printf("Неудалось получить список контекстов пользователя id:'%v': %v", userID, err)
		return
	}

	for _, ct := range contexts {
		ct.TrimSpace()
		if ct.Name == contextName {
			_, err = DB.Exec(queryUpdateUserContext, userID, ct.Id)

			if nil != err {
				log.Printf("Неудалось обновить контекст для пользователя id:'%v': %v", userID, err)
			}

			return
		}
	}
}

// SetResetUserContext обнуляет контекст пользователя с указанным идентификатором.
//
// UserId - идентификатор пользователя.
//
// Возвращаемое значение: true, если успешно, false - в случае ошибки.
func ResetUserContext(UserId int64) bool {
	_, err := DB.Exec(queryUpdateUserContext, UserId, 0)

	if nil != err {
		log.Printf("Неудалось обнулить контекст для пользователя id:'%v': %v", UserId, err)
		return false
	}

	return true
}

func RenameContextByID(contextID *int, newName *string) bool {
	_, err := DB.Exec(queryRenameContextByID, *contextID, *newName)

	if nil != err {
		log.Printf("Неудалось изменить имя контекста id:'%v': %v", *contextID, err)
		return false
	}

	return true
}

func RemoveContextByID(contextID int) bool {
	result, err := DB.Exec(queryRemoveContextByID, contextID)

	if nil != err {
		log.Printf("Неудалось удалить контекст id:'%v': %v", contextID, err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func AddWaitingUser(userID *int64, byPassword bool, userName, key *string) bool {
	result, err := DB.Exec(queryAddWaitingUser, *userName, *key, byPassword, *userID)

	if nil != err {
		log.Printf("Неудалось записать на ожидание пользователя:'%v': %v", *userName, err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func GetWaitingUser(key *string, byPassword bool) (name string) {
	err := DB.Get(&name, queryGetWaitingUser, *key, byPassword)
	if nil != err && sql.ErrNoRows != err {
		log.Panicf("Неудалось запросить пользователя '%v', из ожидания: %v", *key, err)
	}
	return
}
