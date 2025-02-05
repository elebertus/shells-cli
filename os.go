package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KarpelesLab/rest"
)

type ShellOs struct {
	Id   string `json:"Shell_OS__"`
	Name string

	URL     string
	Default string   // Y|N
	Ready   string   // Y|N
	Visible string   // Y|N
	Beta    string   // Y|N
	Public  string   // Y|N
	Family  string   // linux|windows|macos|android|unknown
	Boot    string   // guest-linux|bios|efi
	CPU     string   // x86_64
	Purpose string   // unknown|desktop|server|mobile
	Cmdline string   // cmdline for guest-linux
	Flags   struct{} // byol_warning
}

type ShellOsImage struct {
	Id       string `json:"Shell_OS_Image__"`
	Version  string
	QA       string `json:"QA_Passed"` // Y, P or N
	Filename string
	Format   string
	Source   string
	Status   string
	Size     string // as string because might be too large to fit
	Hash     string
	// Created timestamp
}

func osList(ri *runInfo) error {
	// list available shells
	var list []ShellOs

	err := ri.auth.Apply(context.Background(), "Shell/OS", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shos := range list {
		fmt.Fprintf(os.Stdout, "%s %s\r\n", shos.Id, shos.Name)
	}

	return nil
}

func osImgList(ri *runInfo) error {
	// list available shells
	var list []ShellOsImage

	osId := ri.flags["os"]

	err := ri.auth.Apply(context.Background(), "Shell/OS/"+osId+"/Image", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, img := range list {
		fmt.Fprintf(os.Stdout, "%s %s %s\r\n", img.Id, img.Version, img.Filename)
	}

	return nil
}

func osImgUpload(ri *runInfo) error {
	osId := ri.flags["os"]
	fn := ri.flags["file"]

	fp, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("while trying to open file to upload: %w", err)
	}
	defer fp.Close()

	log.Printf("Uploading %s ...", fn)

	res, err := rest.Upload(ri.auth.token.Use(context.Background()), "Shell/OS/"+osId+"/Image:upload", "POST", map[string]interface{}{"filename": filepath.Base(fn)}, fp, "application/octet-stream")
	if err != nil {
		return err
	}

	var img *ShellOsImage
	err = res.Apply(&img)
	if err != nil {
		return err
	}

	// show info
	fmt.Fprintf(os.Stdout, "%s %s %s\r\n", img.Id, img.Version, img.Filename)

	return nil
}
