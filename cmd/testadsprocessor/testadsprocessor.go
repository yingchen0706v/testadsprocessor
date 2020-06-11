package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/ini.v1"
)

func main() {

	conf, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	rootPath :=conf.Section("").Key("root").String()
	newRootDir := path.Base(rootPath) + "_processed"

	newRoot := path.Join(path.Dir(rootPath) , newRootDir)
	fmt.Printf("Begin processing, target folder: %s\n", newRoot)
	if _, err:= os.Stat(newRoot); !os.IsNotExist(err) {
		fmt.Printf("Target folder exists, please backup, clear and retry!\n")
		return
	}

	filepath.Walk(rootPath, func(admPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			newFolderPath := strings.Replace(admPath, path.Base(rootPath), newRootDir, 1)
			if _, err:= os.Stat(newFolderPath); os.IsNotExist(err) {
				os.Mkdir(newFolderPath, 0777)
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}

		testAds, err := ioutil.ReadFile(admPath)
		if err != nil {
			fmt.Printf("Error reading file %s\n", admPath)
			return nil
		}

		adm:= gjson.Get(string(testAds), "ads.0.ad_markup")
		if adm.Exists() {
			var out bytes.Buffer
			err := json.Indent(&out, []byte(adm.String()), "", "  ")
			if err != nil {
				fmt.Printf("Error set indent for ad markup for file %s\n", admPath)
				return nil
			}

			newAdmPath := strings.Replace(admPath, path.Base(rootPath), newRootDir, 1)
			fileName := path.Base(newAdmPath)
			reg, err := regexp.Compile("[^a-zA-Z0-9_.]+")
			if err != nil {
				panic(err)
			}
			fileName = reg.ReplaceAllString(fileName, "_")
			ioutil.WriteFile(path.Join(path.Dir(newAdmPath), fileName), out.Bytes(), 0777)
			fmt.Printf("Processed file %s\n", admPath)
		} else {
			fmt.Printf("Error processing file %s, ad_markup does not exist.\n", admPath)
		}
		return nil
	})

	fmt.Println("Done processing!")
}
