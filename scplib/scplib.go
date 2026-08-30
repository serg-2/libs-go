package scplib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrLockFound = errors.New("remote lock file detected")

// Вспомогательная функция для проверки ответов по протоколу SCP
func checkSCPResponse(r io.Reader) error {
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return err
	}

	// 0 означает успех. 1 - предупреждение/ошибка, 2 - критическая ошибка
	if buf[0] != 0 {
		// Пытаемся прочитать текст ошибки до символа новой строки
		var errMsg []byte
		b := make([]byte, 1)
		for {
			if _, err := r.Read(b); err != nil || b[0] == '\n' {
				break
			}
			errMsg = append(errMsg, b[0]...)
		}
		return fmt.Errorf("сервер вернул код %d: %s", buf[0], string(errMsg))
	}
	return nil
}

func SendViaScp(sshUser, sshHostPort, localFilePath, remoteTargetDir, lockFile string, overwrite bool) error {
	authMethod, err := getSystemSshAuthMethod()
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", sshHostPort, sshConfig)
	if err != nil {
		return fmt.Errorf("не удалось подключиться по SSH: %w", err)
	}
	defer client.Close()

	// Для удаленного сервера ВСЕГДА используем path вместо filepath (чтобы избежать '\' на Windows)
	fileName := filepath.Base(localFilePath)

	if lockFile != "" {
		lockSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("не удалось создать сессию для проверки замка: %w", err)
		}
		lockPath := path.Join(remoteTargetDir, lockFile)
		err = lockSession.Run(fmt.Sprintf("test -f \"%s\"", lockPath))
		lockSession.Close()
		if err == nil {
			fmt.Println("[Блокировка] Обнаружена блокировка.")
			return ErrLockFound // Убедитесь, что ErrLockFound объявлен в вашем пакете
		}
	}

	remotePath := path.Join(remoteTargetDir, fileName)

	if !overwrite {
		checkSession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("не удалось создать сессию для проверки файла: %w", err)
		}
		err = checkSession.Run(fmt.Sprintf("test -f \"%s\"", remotePath))
		checkSession.Close()
		if err == nil {
			fmt.Printf("[Пропуск] Файл %s уже существует на сервере, отправка отменена.\n", fileName)
			return nil
		}
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть локальный файл: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о локальном файле: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("не удалось создать SSH сессию для SCP: %w", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания StdinPipe: %w", err)
	}
	defer stdin.Close()

	// Получаем также Stdout, чтобы читать подтверждения (ACK) от сервера
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания StdoutPipe: %w", err)
	}

	// ВАЖНО: scp -t ожидает ДИРЕКТОРИЮ, а не путь к файлу
	remoteCommand := fmt.Sprintf("scp -t \"%s\"", remoteTargetDir)
	if err := session.Start(remoteCommand); err != nil {
		return fmt.Errorf("не удалось запустить команду scp на сервере: %w", err)
	}

	// 1. Отправляем протокольное описание файла
	fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode()&os.ModePerm, fileInfo.Size(), fileName)
	if err := checkSCPResponse(stdout); err != nil {
		return fmt.Errorf("сервер отклонил заголовок файла: %w", err)
	}

	// 2. Копируем содержимое файла
	_, err = io.Copy(stdin, file)
	if err != nil {
		return fmt.Errorf("ошибка копирования содержимого файла: %w", err)
	}

	// 3. Отправляем знак окончания передачи
	fmt.Fprint(stdin, "\x00")
	if err := checkSCPResponse(stdout); err != nil {
		return fmt.Errorf("сервер отклонил завершение передачи файла: %w", err)
	}

	stdin.Close()
	if err := session.Wait(); err != nil {
		return fmt.Errorf("ошибка завершения сессии scp: %w", err)
	}

	fmt.Printf("Файл %s успешно отправлен на %s\n", fileName, sshHostPort)
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
