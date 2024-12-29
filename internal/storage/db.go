package storage

import (
	"database/sql"
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

func (d *DB) InsertUser(data *models.Input) (int, bool, error) {
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
		return 0, false, err
	}
	rowsAffected, _ := result.RowsAffected()
	fileExistsQuery := `select exists(select * from uploaded_files where email = ? limit 1)`
	var fileExists bool
	err = d.Conn.QueryRow(fileExistsQuery, data.Email).Scan(&fileExists)
	if err != nil {
		log.Printf("Database execution error: %v", err)
		return 0, false, err
	}
	return int(rowsAffected), fileExists, nil
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
		WHERE email = ? and is_deleted = "N";`

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

func (db *DB) DeleteFile(fileID string) (string, string, error) {
	tx, err := db.Conn.Begin()
	if err != nil {
		return "", "", err
	}

	var fileName string
	var uniqueName string
	querySelectUpdate := `
		SELECT file_name , unique_name
		FROM uploaded_files
		WHERE file_id = ? 
		FOR UPDATE
	`
	err = tx.QueryRow(querySelectUpdate, fileID).Scan(&fileName, &uniqueName)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	queryUpdate := `
		UPDATE uploaded_files
		SET is_deleted = 'Y', deleted_time = CURRENT_TIMESTAMP
		WHERE file_id = ?
	`
	_, err = tx.Exec(queryUpdate, fileID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", "", err
	}
	return fileName, uniqueName, nil
}

func (db *DB) GetOCIFileName(fileID string) (*models.FileName, error) {
	query := `select file_name , unique_name from uploaded_files where file_id = ?`

	row := db.Conn.QueryRow(query, fileID)

	var metadata models.FileName
	err := row.Scan(&metadata.FileName, &metadata.OCIFileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to fetch file metadata: %w", err)
	}

	return &metadata, nil

}
