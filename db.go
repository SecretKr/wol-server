package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Device struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	MAC  string `json:"mac"`
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite", "./data/devices.db")
	if err != nil {
		panic(err)
	}

	statement, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT, 
			mac TEXT UNIQUE
		)`)
	statement.Exec()
}

func GetDevices() ([]Device, error) {
	rows, err := db.Query("SELECT id, name, mac FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		err = rows.Scan(&d.ID, &d.Name, &d.MAC)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func AddDevice(name, mac string) error {
	statement, err := db.Prepare("INSERT INTO devices (name, mac) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(name, mac)
	return err
}

func DeleteDevice(id int) error {
	statement, err := db.Prepare("DELETE FROM devices WHERE id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("device not found")
	}
	return nil
}
