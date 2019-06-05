package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ExportURL ...
type ExportURL struct {
	URL string
}

// UserAudience ...
type UserAudience struct {
	email    string
	birthday string
	telefone string
}

// ExportDiff ...
type ExportDiff struct {
	isNew bool
	isOld bool
	data  UserAudience
}

func main() {
	appKey := os.Getenv("APP_KEY")
	m := make(map[string]*ExportDiff)

	diff, err := handleDiff(appKey, os.Getenv("OLD_EXPORT_ID"), true, false, m)
	if err != nil {
		panic(err)
	}

	diff, err = handleDiff(appKey, os.Getenv("NEW_EXPORT_ID"), false, true, diff)
	if err != nil {
		panic(err)
	}

	add, remove := splitDiff(diff)
	fmt.Println("************* Count add *************")
	fmt.Println(len(add))
	fmt.Println("************* Array add *************")
	fmt.Println(add)
	fmt.Println("************ Count remove ***********")
	fmt.Println(len(remove))
	fmt.Println("************ Array remove ***********")
	fmt.Println(remove)
}

func splitDiff(diff map[string]*ExportDiff) (add []interface{}, remove []interface{}) {
	for _, u := range diff {
		if !u.isOld && u.isNew {
			add = append(add, u.data)
		}

		if u.isOld && !u.isNew {
			remove = append(remove, u.data)
		}
	}
	return add, remove
}

func handleDiff(appKey, id string, isOld, isNew bool, m map[string]*ExportDiff) (map[string]*ExportDiff, error) {
	url, err := getExportURL(appKey, id)
	if err != nil {
		return nil, err
	}
	diff, err := proccessRemoteCsv(url, isOld, isNew, m)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func getExportURL(appKey, id string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf(
		"%v/%v/url?appKey=%v",
		os.Getenv("EXPORT_URL"),
		id,
		appKey,
	)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	e := ExportURL{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&e); err != nil {
		return "", err
	}

	return e.URL, err
}

func proccessRemoteCsv(url string, isOld, isNew bool, m map[string]*ExportDiff) (map[string]*ExportDiff, error) {
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

func generateDiff(arq io.ReadCloser, isOld, isNew bool, m map[string]*ExportDiff) (map[string]*ExportDiff, error) {
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
			m[user[0]].isNew = isNew
			m[user[0]].data = UserAudience{
				email:    user[2],
				birthday: user[5],
				telefone: user[len(user)-1],
			}
			continue
		}

		m[user[0]] = &ExportDiff{
			isOld: isOld,
			isNew: isNew,
			data: UserAudience{
				email:    user[2],
				birthday: user[5],
				telefone: user[len(user)-1],
			},
		}
	}

	return m, nil
}
