package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/pdropaiva/poc-diff-csv/cmd/util"
	"github.com/pdropaiva/poc-diff-csv/domain"
)

func main() {
	appKey := os.Getenv("APP_KEY")
	m := make(map[string]*domain.ExportDiff)

	diff, err := handleDiff(appKey, os.Getenv("OLD_EXPORT_ID"), true, false, m)
	if err != nil {
		panic(err)
	}

	diff, err = handleDiff(appKey, os.Getenv("NEW_EXPORT_ID"), false, true, diff)
	if err != nil {
		panic(err)
	}

	add, remove := util.SplitDiff(diff)
	util.PrintDiff(add, remove)
}

func handleDiff(appKey, id string, isOld, isNew bool, m map[string]*domain.ExportDiff) (map[string]*domain.ExportDiff, error) {
	e := domain.Export{}
	url, err := e.AssignedURL(appKey, id)
	if err != nil {
		return nil, err
	}
	diff, err := proccessRemoteCsv(url, isOld, isNew, m)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func proccessRemoteCsv(url string, isOld, isNew bool, m map[string]*domain.ExportDiff) (map[string]*domain.ExportDiff, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	diff, err := generateDiff(resp.Body, isOld, isNew, m)
	if err != nil {
		return nil, err
	}

	return diff, err
}

func generateDiff(arq io.ReadCloser, isOld, isNew bool, m map[string]*domain.ExportDiff) (map[string]*domain.ExportDiff, error) {
	defer arq.Close()
	r := csv.NewReader(arq)
	for {
		user, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if m[user[0]] != nil {
			m[user[0]].IsNew = isNew
			m[user[0]].Data = domain.UserAudience{
				Email:    user[2],
				Birthday: user[5],
				Telefone: user[len(user)-1],
			}
			continue
		}

		m[user[0]] = &domain.ExportDiff{
			IsOld: isOld,
			IsNew: isNew,
			Data: domain.UserAudience{
				Email:    user[2],
				Birthday: user[5],
				Telefone: user[len(user)-1],
			},
		}
	}

	return m, nil
}
