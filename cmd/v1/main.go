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

func generateDiff(old, new string) (map[string]*ExportDiff, error) {
	m := make(map[string]*ExportDiff)
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

		m[user[0]] = &ExportDiff{isOld: true, isNew: false, data: UserAudience{email: user[2], birthday: user[5], telefone: user[len(user)-1]}}
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
			m[user[0]].isNew = true
			m[user[0]].data = UserAudience{email: user[2], birthday: user[5], telefone: user[len(user)-1]}
			continue
		}

		m[user[0]] = &ExportDiff{isOld: false, isNew: true, data: UserAudience{email: user[2], birthday: user[5], telefone: user[len(user)-1]}}
	}

	return m, nil
}

func downloadExport(appKey, id string) (string, error) {
	url, err := getExportURL(appKey, id)
	if err != nil {
		return "", err
	}
	filepath, err := downloadFile(fmt.Sprintf("%v.csv", id), url)
	if err != nil {
		return "", err
	}

	return filepath, nil
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
