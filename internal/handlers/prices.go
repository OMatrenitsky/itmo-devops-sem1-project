package handlers

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"io"
	"net/http"
	"os"
)

func PostPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		tmp, err := os.CreateTemp("", "*.zip")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer os.Remove(tmp.Name())

		io.Copy(tmp, file)
		tmp.Close()

		zipReader, err := zip.OpenReader(tmp.Name())
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer zipReader.Close()

		for _, f := range zipReader.File {
			if f.Name != "data.csv" {
				continue
			}

			rc, _ := f.Open()
			reader := csv.NewReader(rc)

			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

				_, err = db.Exec(
					`INSERT INTO prices (id, created_at, name, category, price)
					 VALUES ($1, $2, $3, $4, $5)`,
					record[0], record[1], record[2], record[3], record[4],
				)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
			rc.Close()
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(
			`SELECT id, created_at, name, category, price FROM prices`,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")

		zipWriter := zip.NewWriter(w)
		csvFile, _ := zipWriter.Create("data.csv")
		writer := csv.NewWriter(csvFile)

		for rows.Next() {
			var id, createdAt, name, category, price string
			rows.Scan(&id, &createdAt, &name, &category, &price)
			writer.Write([]string{id, createdAt, name, category, price})
		}

		writer.Flush()
		zipWriter.Close()
	}
}
