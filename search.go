package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

type NixPackage struct {
	Meta struct {
		Description string   `json:"description"`
		Homepage    string   `json:"homepage"`
		Platforms   []string `json:"platforms"`
		License     []struct {
			URL       string `json:"url"`
			ShortName string `json:"shortName"`
			SpdxID    string `json:"spdxId"`
			FullName  string `json:"fullName"`
		} `json:"license"`
		Maintainers []interface{} `json:"maintainers"`
		Position    string        `json:"position"`
	} `json:"meta"`
	System string `json:"system"`
	Name   string `json:"name"`
}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func nixos_hash() string {
	out, cmd_err := exec.Command("nixos-version", "--hash").Output()
	// TODO: Remove double quotes in result
	//out, cmd_err := exec.Command("nix-instantiate", "--eval", "-E", "(import <nixpkgs> {}).lib.nixpkgsVersion").Output()
	check(cmd_err)
	result := fmt.Sprintf("%s", out)
	return strings.TrimSpace(result)
}

func Search(pattern string) {
	version := nixos_hash()
	cache_file := fmt.Sprintf("/tmp/nixie-%s.json", version)

	// TODO: Flag to update cache
	if _, err := os.Stat(cache_file); os.IsNotExist(err) {
		fmt.Println("Building cache", cache_file)

		out, cmd_err := exec.Command("nix-env", "-qa", "--json", "*").Output()
		check(cmd_err)

		write_err := ioutil.WriteFile(cache_file, out, 0644)
		check(write_err)
	}

	packages_json, read_err := ioutil.ReadFile(cache_file)
	check(read_err)

	var packages map[string]NixPackage
	json_err := json.Unmarshal(packages_json, &packages)
	if json_err != nil {
		// TODO: Figure out which thing is not matching the expected struct
		log.Println(json_err)
	}

	// Sort packages by key
	keys := make([]string, len(packages))
	i := 0
	for k, _ := range packages {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// TODO: Validate input
	r, _ := regexp.Compile(pattern)

	//for key, value := range packages {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, key := range keys {
		if r.MatchString(key) || r.MatchString(packages[key].Name) {
			fmt.Fprintf(writer, "%v\t%v\t%v\t", packages[key].Name, key, packages[key].Meta.Description)
			fmt.Fprintln(writer)
		}
	}
	fmt.Fprintln(writer)
	writer.Flush()
}
