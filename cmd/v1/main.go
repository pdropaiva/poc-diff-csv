package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pdropaiva/poc-diff-csv/cmd/util"
	"github.com/pdropaiva/poc-diff-csv/domain"
)

func main() {
	appKey := os.Getenv("APP_KEY")

	oldExpo, err := downloadExport(appKey, os.Getenv("OLD_EXPORT_ID"))
	if err != nil {
		panic(err)
	}

	newExpo, err := downloadExport(appKey, os.Getenv("NEW_EXPORT_ID"))
	if err != nil {
		panic(err)
	}

	diff, err := generateDiff(oldExpo, newExpo)
	if err != nil {
		panic(err)
	}

	add, remove := util.SplitDiff(diff)
	util.PrintDiff(add, remove)
}

func downloadExport(appKey, id string) (string, error) {
	e := domain.Export{}
	url, err := e.AssignedURL(appKey, id)
	if err != nil {
		return "", err
	}
	filepath, err := downloadFile(fmt.Sprintf("%v.csv", id), url)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func downloadFile(filepath string, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filepath = fmt.Sprintf("./csv/%v", filepath)
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filepath, err
}

func generateDiff(old, new string) (map[string]*domain.ExportDiff, error) {
	m := make(map[string]*domain.ExportDiff)
	oldFile, err := os.Open(old)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(oldFile)
	for {
		user, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		m[user[0]] = &domain.ExportDiff{
			IsOld: true,
			IsNew: false,
			Data: domain.UserAudience{
				Email:    user[2],
				Birthday: user[5],
				Telefone: user[len(user)-1],
			},
		}
	}

	newFile, err := os.Open(new)
	if err != nil {
		return nil, err
	}

	r = csv.NewReader(newFile)
	for {
		user, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if m[user[0]] != nil {
			m[user[0]].IsNew = true
			m[user[0]].Data = domain.UserAudience{
				Email:    user[2],
				Birthday: user[5],
				Telefone: user[len(user)-1],
			}
			continue
		}

		m[user[0]] = &domain.ExportDiff{
			IsOld: false,
			IsNew: true,
			Data: domain.UserAudience{
				Email:    user[2],
				Birthday: user[5],
				Telefone: user[len(user)-1],
			},
		}
	}

	return m, nil
}
