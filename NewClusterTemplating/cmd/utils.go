// Copyright © 2021 Lucidworks
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var ctx = context.Background()

func templateDirectory(values interface{}, templateDir string, outputDir string) {
	var files []string
	err := filepath.Walk(templateDir, func(templateDir string, info os.FileInfo, err error) error {
		files = append(files, templateDir)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		println(file)
		if fileExists(file) {
			filename := generateTemplatedFileName(file, values, templateDir, outputDir)
			writeTemplateFile(values, file, filename)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, nil)
}

func generateTemplatedFileName(file string, values interface{}, templateDir string, outputDir string) string {
	template, _ := template.New("name").Parse(file)
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, values); err != nil {
		panic(err)
	}
	file = tpl.String()
	if strings.HasSuffix(file, ".tmpl") {
		return strings.Replace(file[0:len(file)-5], templateDir, outputDir, 1)
	}
	return strings.Replace(file, templateDir, outputDir, 1)
}

func writeTemplateFile(values interface{}, f string, nf string) {
	tpl := templateFile(values, f)
	write(nf, tpl.Bytes())
}

func templateFile(values interface{}, f string) bytes.Buffer {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["HasPrefix"] = strings.HasPrefix
	tplFuncMap["ReplaceAll"] = strings.ReplaceAll

	template, _ := template.New("file").Delims("[[", "]]").Funcs(tplFuncMap).Parse(text)
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, values); err != nil {
		panic(err)
	}
	return tpl
}

func write(filename string, data []byte) {
	dir, _ := path.Split(filename)
	os.MkdirAll(dir, 0755)
	if err := ioutil.WriteFile(filename, data, 0755); err != nil {
		panic(err)
	}
}

func getSecretFromGCP(secretName string) string {

	// Create the client.
	// ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/managed-fusion/secrets/%s/versions/latest", secretName),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data)
}


func getFilesInDir(dir string) []string {
	var filesInDir []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println(err)
            return err
        }
		if !info.IsDir() {
			filesInDir = append(filesInDir, path)
		}
        return nil
    })
	return filesInDir
}