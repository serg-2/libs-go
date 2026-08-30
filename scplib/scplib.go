package scplib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrLockFound = errors.New("remote lock file detected")

// SendViaScp проверяет наличие lockFile, если он есть, затем наличие файла если запрещён overwrite, и только потом отправляет его
func SendViaScp(sshUser, sshHostPort, localFilePath, remoteTargetDir, lockFile string, overwrite bool) error {
	// Получаем метод авторизации по системным ключам
	authMethod, err := getSystemSshAuthMethod()
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	// Конфигурация SSH клиента
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Установка SSH соединения
	client, err := ssh.Dial("tcp", sshHostPort, sshConfig)
	if err != nil {
		return fmt.Errorf("не удалось подключиться по SSH: %w", err)
	}
	defer client.Close()

	if lockFile != "" {
		lockSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("не удалось создать сессию для проверки замка: %w", err)
		}
		lockPath := filepath.Join(remoteTargetDir, lockFile)
		err = lockSession.Run(fmt.Sprintf("test -f \"%s\"", lockPath))
		lockSession.Close()
		if err == nil {
			fmt.Println("[Блокировка] Обнаружена блокировка.")
			return ErrLockFound
		}
	}

	// Извлекаем имя файла из полного пути
	fileName := filepath.Base(localFilePath)
	remotePath := filepath.Join(remoteTargetDir, fileName)

	if !overwrite {
		// Существует ли целевой файл на сервере
		checkSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("не удалось создать сессию для проверки файла: %w", err)
		}
		// Выполняем команду 'test -f'. Она вернет системный статус 0, если файл существует
		err = checkSession.Run(fmt.Sprintf("test -f \"%s\"", remotePath))
		checkSession.Close()
		if err == nil {
			fmt.Printf("[Пропуск] Файл %s уже существует на сервере, отправка отменена.\n", fileName)
			return nil
		}
	}

	// Чтение локального файла
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть локальный файл: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о локальном файле: %w", err)
	}

	// Открытие новой сессии для передачи файла по SCP
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("не удалось создать SSH сессию для SCP: %w", err)
	}
	defer session.Close()

	// Настройка stdin для сессии, чтобы передать протокольные команды SCP
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания StdinPipe: %w", err)
	}
	defer stdin.Close()

	// Запуск команды scp на удаленном сервере в режиме приема
	remoteCommand := fmt.Sprintf("scp -t \"%s\"", remotePath)
	if err := session.Start(remoteCommand); err != nil {
		return fmt.Errorf("не удалось запустить команду scp на сервере: %w", err)
	}

	// Протокол SCP: отправляем Unix-права, размер и имя файла
	fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode()&os.ModePerm, fileInfo.Size(), fileName)

	// Копируем содержимое файла в поток SSH
	_, err = io.Copy(stdin, file)
	if err != nil {
		return fmt.Errorf("ошибка копирования содержимого файла: %w", err)
	}

	// Отправляем нулевой байт (\x00) как знак окончания передачи файла
	fmt.Fprint(stdin, "\x00")

	// Закрываем stdin и ждем завершения удаленной команды scp
	stdin.Close()
	if err := session.Wait(); err != nil {
		return fmt.Errorf("ошибка завершения сессии scp (сервер отклонил файл): %w", err)
	}

	fmt.Printf("Файл %s успешно отправлен на %s\n", filepath.Base(localFilePath), sshHostPort)
	return nil
}

// getSystemSshAuthMethod ищет существующие SSH-ключи в системе и возвращает метод авторизации
func getSystemSshAuthMethod() (ssh.AuthMethod, error) {
	// Получаем путь к домашней директории текущего пользователя системы (~/)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("не удалось определить домашнюю директорию: %w", err)
	}

	// Список стандартных путей к приватным ключам в порядке приоритета
	keyPaths := []string{
		// filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_rsa"),
	}

	var keyBytes []byte
	var foundPath string

	// Ищем первый существующий ключ
	for _, path := range keyPaths {
		bytes, err := os.ReadFile(path)
		if err == nil {
			keyBytes = bytes
			foundPath = path
			break
		}
	}

	if len(keyBytes) == 0 {
		// return nil, fmt.Errorf("системные SSH-ключи не найдены в ~/.ssh/ (проверены id_ed25519 и id_rsa)")
		return nil, fmt.Errorf("системные SSH-ключи не найдены в ~/.ssh/ (проверены id_rsa)")
	}

	// Парсим приватный ключ
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		// Если ключ защищен паролем (passphrase), здесь вернется ошибка.
		return nil, fmt.Errorf("не удалось распарсить приватный ключ %s: %w", foundPath, err)
	}

	// fmt.Printf("Используется системный ключ: %s\n", foundPath)
	return ssh.PublicKeys(signer), nil
}
