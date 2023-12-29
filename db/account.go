package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func AddAccountRecord(uuid []byte, username string, key, salt []byte) error {
	_, err := handle.Exec("INSERT INTO accounts (uuid, username, key, salt, registered) VALUES (?, ?, ?, ?, UTC_TIMESTAMP())", uuid, username, key, salt)
	if err != nil {
		return err
	}

	return nil
}

func AddAccountSession(username string, token []byte) error {
	_, err := handle.Exec("INSERT INTO sessions (token, uuid, expire) SELECT a.uuid, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 WEEK) FROM accounts a WHERE a.username = ?", token, username)
	if err != nil {
		return err
	}

	return nil
}

func GetUsernameFromToken(token []byte) (string, error) {
	var username string
	err := handle.QueryRow("SELECT a.username FROM accounts a JOIN sessions s ON s.uuid = a.uuid WHERE s.token = ? AND s.expire > UTC_TIMESTAMP()").Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}

		return "", err
	}

	return username, nil
}

func GetAccountKeySaltFromUsername(username string) ([]byte, []byte, error) {
	var key, salt []byte
	err := handle.QueryRow("SELECT key, salt FROM accounts WHERE username = ?", username).Scan(&key, &salt)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

func GetUUIDFromToken(token []byte) ([]byte, error) {
	var uuid []byte
	err := handle.QueryRow("SELECT uuid FROM sessions WHERE token = ? AND expire > UTC_TIMESTAMP()", token).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	return uuid, nil
}

func RemoveSessionFromToken(token []byte) error {
	_, err := handle.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return err
	}

	return nil
}
