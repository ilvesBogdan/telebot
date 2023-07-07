package bot

type LanguagePackageOfMessages struct {
	language uint
}

func setLanguage(language string) LanguagePackageOfMessages {
	var objectlanguage LanguagePackageOfMessages
	objectlanguage.set(language)
	return objectlanguage
}

func (s *LanguagePackageOfMessages) set(language string) {
	switch language {
	case "ru":
		s.language = 1
	}
}

func (s *LanguagePackageOfMessages) msg(msg string) string {
	switch msg {

	// Сообщение об ошибке
	case "error":
		switch s.language {
		case 1:
			return "*Ошибка*"
		}

	// Сообщение о переименовании
	case "Rename":
		switch s.language {
		case 1:
			return "Переименовать"
		}

	// Сообщение о удалении
	case "Remove":
		switch s.language {
		case 1:
			return "Удалить"
		}

	// Сообщение о очистке
	case "Clear":
		switch s.language {
		case 1:
			return "Очистить"
		}

	// Сообщение вывода информации о контексте
	case "ContextCnfgMsg":
		switch s.language {
		case 1:
			return `*%v*
			Контекст был создан: _%v_`
		}

	// Сообщение уведомляет об отсутвии регистрации
	case "NoRegistration":
		switch s.language {
		case 1:
			return "*Вы незарегистрированны*"
		}

	// Сообщение уведомляет о том что контекст с таким
	// именем уже существует
	case "ThisContextIsAlreadySelected":
		switch s.language {
		case 1:
			return "Этот контекст уже выбран."
		}

	// Сообщение уведомляет о том что
	// контекст был успешно изменен
	case "ContextChanged":
		switch s.language {
		case 1:
			return "Контекст изменен"
		}

	// Сообщение уведомляет о том что
	// контекст не найден
	case "ContextNotFond":
		switch s.language {
		case 1:
			return "Контекст не найден"
		}

	// Сообщение с прозьбой ввести имя контекста
	case "EnterContextName":
		switch s.language {
		case 1:
			return "Введите имя нового контекста"
		}
	}

	// Если сообщение небыло найдено
	return "error text"
}
