package storage

import (
	"file_mgmt_system/internal/models"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import for MySQL driver
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Conn *sqlx.DB
}

func NewDB(user, password, host, dbname string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)

	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to the database successfully")
	return &DB{Conn: conn}, nil
}

func (d *DB) InsertUser(data *models.Input) (int, error) {
	query := `
		INSERT INTO users (email, first_name, last_name, phone) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		first_name = VALUES(first_name),
		last_name = VALUES(last_name),
		phone = VALUES(phone)`

	result, err := d.Conn.Exec(query, data.Email, data.FirstName, data.LastName, data.Phone)
	if err != nil {
		log.Printf("Database execution error: %v", err)
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (db *DB) SaveFileMetadata(file models.FileMetadata) error {
	query := `
		INSERT INTO uploaded_files (file_name, unique_name, file_type, file_size, email, upload_time, oci_reference , file_id)
		VALUES (?, ?, ?, ?, ?, ?, ? , ?)
	`

	// Execute the query with the provided file metadata
	_, err := db.Conn.Exec(
		query,
		file.FileName,
		file.UniqueName,
		file.FileType,
		file.FileSize,
		file.Email,
		file.UploadTime,
		file.OCIReference,
		file.FileID,
	)

	return err
}

func (db *DB) GetFileList(email string) ([]models.FileMetadata, error) {
	query := `
		SELECT file_name, unique_name, file_type, file_size, email, upload_time, oci_reference ,file_id
		FROM uploaded_files 
		WHERE email = ?;`

	rows, err := db.Conn.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.FileMetadata

	for rows.Next() {
		var file models.FileMetadata
		err := rows.Scan(
			&file.FileName,
			&file.UniqueName,
			&file.FileType,
			&file.FileSize,
			&file.Email,
			&file.UploadTime,
			&file.OCIReference,
			&file.FileID,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
