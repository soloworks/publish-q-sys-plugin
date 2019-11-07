package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/soloworks/go-nuget-utils"
	"github.com/soloworks/go-nuspec"
)

func main() {

	// make new Action
	a := newAction()

	// Get and verify the Version
	ver := a.GetInput("version")
	if ver == "commit-tag" {
		if strings.HasPrefix(os.Getenv("GITHUB_REF"), "refs/tags/") {
			// Extract tag from end of Envar
			ver = strings.TrimPrefix(os.Getenv("GITHUB_REF"), "refs/tags/")
		} else {
			ver = ""
		}
	}

	// If no Version is found, exit with error
	if ver == "" {
		log.Fatalln("No version supplied or found ($GITHUB_REF="+os.Getenv("GITHUB_REF"), ")")
	}

	// Check ver against RegEx to be in format x.x.x.x only
	if match, _ := regexp.MatchString(`^(([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+))$`, ver); !match {
		log.Fatalln("SemVer format incorrect: ", ver)
	}

	// Get nuspec file - if not arg then first in directory
	nsfilename := ""
	if a.GetInput("nuspec-file") != "" {
		nsfilename = a.GetInput("nuspec-file")
	} else {
		m, err := filepath.Glob("*.nuspec")
		if err != nil {
			log.Fatal(err)
		}
		for _, x := range m {
			nsfilename = x
			break
		}
	}

	// If no NuSpec is found, exit with error
	if nsfilename == "" {
		log.Fatalln("No NuSpec supplied or found")
	}

	// Load the NuSpec file
	nsf, err := nuspec.FromFile(nsfilename)
	if err != nil {
		log.Fatalln(err)
	}

	// Adjust Version to match provided
	nsf.Meta.Version = ver

	// Edit all .qplug files to update version from 0.0.0.0-master
	m, err := filepath.Glob("content/*.qplug")
	if err != nil {
		log.Fatalln(err)
	}
	for _, x := range m {
		read, err := ioutil.ReadFile(x)
		if err != nil {
			log.Fatalln(err)
		}
		newContents := strings.Replace(string(read), "0.0.0.0-master", ver, -1)
		err = ioutil.WriteFile(x, []byte(newContents), 0)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Gather Git Release notes

	// Convert MarkDown to Description
	if x := a.GetInput("md-to-desc"); x != "" {
		// Load the file
		read, err := ioutil.ReadFile(x)
		if err != nil {
			log.Fatalln(err)
		}
		// Do conversion to NuSpec formats - Images
		reImg := regexp.MustCompile(`(?:!\[(.*?)\]\((.*?)\))`)
		reBraces := regexp.MustCompile(`\((.*)\)`)
		nsf.Meta.Description = reImg.ReplaceAllStringFunc(string(read), func(s string) string {
			log.Println(s)
			return reBraces.ReplaceAllStringFunc(string(s), func(s string) string {
				p := path.Clean(s[1 : len(s)-1])
				return "(http://" + p + ")"
			})
		})
	}
	log.Println(nsf.Meta.Description)
	// Pack it up
	_, err = nuget.PackNupkg(nsf, ".", ".")
	if err != nil {
		log.Fatalln(err)
	}

	// // Push it up
	// status, _, err := nuget.PushNupkg(npkg, a.GetInput("Api-Key"), a.GetInput("nuget-host"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// if status == 201 {
	// 	os.Exit(0)
	// } else {
	// 	log.Fatalln("Failed with HTTP Status:", status)
	// }

	// Set Output variables
	//println("::set-output name=duration::", strconv.Itoa(dur))
}
