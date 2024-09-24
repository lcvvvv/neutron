package main

import (
	"flag"
	"fmt"
	"github.com/chainreactors/neutron/protocols"
	http2 "github.com/chainreactors/neutron/protocols/http"
	"github.com/chainreactors/neutron/templates"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var ExecuterOptions *protocols.ExecuterOptions

func main() {
	// 定义命令行参数
	proxyAddr := flag.String("proxy", "", "Proxy address (e.g., http://127.0.0.1:8080)")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: shot [-proxy <proxy_address>] <path_or_file> <target_url>")
		return
	}

	targetPath := flag.Arg(0)
	targetURL := flag.Arg(1)

	if *proxyAddr != "" {
		fmt.Println("Using proxy:", *proxyAddr)
		proxyURL, err := url.Parse(*proxyAddr)
		if err != nil {
			fmt.Printf("Invalid proxy address: %s\n", err.Error())
			return
		}
		http2.DefaultTransport.Proxy = http.ProxyURL(proxyURL)
	}

	var yamlFiles []string
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		return
	}

	for _, yamlFile := range yamlFiles {
		content, err := os.ReadFile(yamlFile)
		if err != nil {
			fmt.Printf("Error reading %s: %s\n", yamlFile, err.Error())
			continue
		}

		t := &templates.Template{}
		err = yaml.Unmarshal(content, t)
		if err != nil {
			fmt.Printf("Error unmarshalling %s: %s\n", yamlFile, err.Error())
			continue
		}

		err = t.Compile(ExecuterOptions)
		if err != nil {
			fmt.Printf("Error compiling %s: %s\n", yamlFile, err.Error())
			continue
		}

		fmt.Printf("Load success for %s\n", yamlFile)
		start := time.Now()
		res, err := t.Execute(targetURL, nil)
		if err == nil {
			fmt.Println("OK:", res)
		} else {
			fmt.Println("Error:", res)
		}
		fmt.Println("Execution time:", time.Since(start))
	}
}
