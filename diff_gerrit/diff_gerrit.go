// @filename           :  diff_gerrit.go
// @author             :  church.zhong@hmdglobal.com
// @date               :  Sat Mar  6 10:23:14 HKT 2021
// @function           :  parse and diff json of gerrit.
// @see                :  https://golang.org/pkg/io/ioutil/#ReadDir
// @require            :  golang 1.16
package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var sep string = string(os.PathSeparator)

type OwnerType struct {
	Email    string `json:"email"`    // church.zhong@hmdglobal.com
	Username string `json:"username"` // church.zhong
}

type Record struct {
	Project       string    `json:"project"` // HMD/platform/vendor/opensource/audio-kernel
	Branch        string    `json:"branch"`  // sm4350_TF_AKT_DEV
	Id            string    `json:"id"`      // I3dcca06f99ce25c06cf79803f84b4a913ee6ebe1
	Number        int       `json:"number"`  // 40868
	Subject       string    `json:"subject"` // <48265><1_1><fengpengbo><aoki><bsp><NA>porting changes of scw to aoki
	Owner         OwnerType `json:"owner"`   //  email: church.zhong@hmdglobal.com;  username: church.zhong
	Url           string    `json:"url"`     // http://hmdgerritserver.southeastasia.cloudapp.azure.com/c/HMD/platform/vendor/opensource/audio-kernel/+/40868
	CommitMessage string    `json:"commitMessage"`
	CreatedOn     int       `json:"createdOn"`   // 2021-03-08 07:04:28 UTC
	LastUpdated   int       `json:"lastUpdated"` // 2021-03-08 07:04:48 UTC
	Open          bool      `json:"open"`        // false
	Status        string    `json:"status"`      // MERGED
}

//{"type":"stats","rowCount":27,"runTimeMilliseconds":16,"moreChanges":false}
type Tail struct {
	Type                string `json:"type"`                // "stats"
	RowCount            int    `json:"rowCount"`            // 27
	RunTimeMilliseconds int    `json:"runTimeMilliseconds"` // 16
	MoreChanges         bool   `json:"moreChanges"`         // false
}

func getVendor(branch string) string {
	var vendor string = ""
	f, err := os.Open("vendor.config")
	if err != nil {
		fmt.Println(err.Error())
	}
	buf := bufio.NewReader(f)
	for {
		b, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			fmt.Println(errR.Error())
		}
		line := string(b)
		//fmt.Println(line)
		nv := strings.Split(line, ":")
		if branch == nv[0] {
			vendor = strings.TrimSpace(nv[1])
			break
		}
	}
	return vendor
}

func getRecords(path string, data []Record) ([]Record, int) {
	var rowCount int = 0
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf := bufio.NewReader(f)
	for {
		b, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			fmt.Println(errR.Error())
		}
		line := string(b)
		//fmt.Println(line)
		var j Record
		errR = json.Unmarshal([]byte(line), &j)
		if 0 == j.Number {
			var r Tail
			json.Unmarshal([]byte(line), &r)
			rowCount = r.RowCount
			//fmt.Println(r)
		} else {
			data = append(data, j)
			//fmt.Println(j)
		}
	}

	return data, rowCount
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Project struct {
	XmlName   xml.Name `xml:"project"`
	Name      string   `xml:"name,attr"`     // "device/common"
	Path      string   `xml:"path,attr"`     // "device/common"
	Revision  string   `xml:"revision,attr"` // "sm4350_TF_AKT_DEV"
	InnerText string   `xml:",innerxml"`
}

type Manifests struct {
	XmlName  xml.Name  `xml:"manifest"`
	Projects []Project `xml:"project"`
}

var gManifest Manifests

func getManifests(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	err = xml.Unmarshal(content, &gManifest)
	if err != nil {
		panic(err.Error())
	}

	// don't print
	//fmt.Println("manifest:", gManifest)
	//fmt.Println("XmlName:", gManifest.XmlName)
	//for _, p := range gManifest.Projects {
	//fmt.Println("Tag:", p.XmlName)
	//fmt.Println("Name:", p.Name)
	//fmt.Println("Path:", p.Path)
	//fmt.Println("Revision:", p.Revision)
	//fmt.Println("InnerText:", p.InnerText)
	//}
}

// Prints how much time has elapsed after main() exit.
func elapsedTime(start time.Time) {
	fmt.Println(time.Since(start))
}
func main() {
	defer elapsedTime(time.Now())

	// input
	var username = flag.String("username", "church.zhong", "your gerrit username")
	var bspBranch = flag.String("bsp", "sm4350_TF_AKT_BSP", "BSP branch name")
	var devBranch = flag.String("dev", "sm4350_TF_AKT_DEV", "DEV branch name")
	var repoPath = flag.String("repo", "/data/church/", "your repo path")
	flag.Parse()
	fmt.Println("username=" + *username)
	fmt.Println("bspBranch=" + *bspBranch)
	fmt.Println("devBranch=" + *devBranch)
	fmt.Println("repoPath=" + *repoPath)

	bspBranchJson := *bspBranch + ".json"
	devBranchJson := *devBranch + ".json"

	if exist, _ := pathExists(bspBranchJson); !exist {
		fmt.Printf("no such file or directory:bspBranchJson=%32s\n", bspBranchJson)
		return
	}
	if exist, _ := pathExists(devBranchJson); !exist {
		fmt.Printf("no such file or directory:devBranchJson=%32s\n", devBranchJson)
		return
	}
	if exist, _ := pathExists(*repoPath); !exist {
		fmt.Printf("no such file or directory:repoPath=%32s\n", *repoPath)
		return
	}

	// work start
	var bspRecords []Record
	var bspRowCount int = 0
	var devRecords []Record
	var devRowCount int = 0
	//bsp
	bspJson := *bspBranch + ".json"
	bspRecords, bspRowCount = getRecords(bspJson, bspRecords)
	if len(bspRecords) != bspRowCount {
		fmt.Printf("len=%v,row=%v\n", len(bspRecords), bspRowCount)
		return
	}

	//dev
	devJson := *devBranch + ".json"
	devRecords, devRowCount = getRecords(devJson, devRecords)
	if len(devRecords) != devRowCount {
		fmt.Printf("len=%v,row=%v\n", len(devRecords), devRowCount)
		return
	}

	//diff
	var diffRecords []Record
	for _, b := range bspRecords {
		//fmt.Println(b)
		found := false
		for _, d := range devRecords {
			//fmt.Println(d)
			if b.Id == d.Id {
				found = true
			}
		}
		if !found {
			diffRecords = append(diffRecords, b)
		}
	}
	diffRecordsLen := len(diffRecords)
	fmt.Printf("len=%v\n", diffRecordsLen)

	if 0 == diffRecordsLen {
		fmt.Println("Do nothing!")
		return
	}

	bspBranchVendor := getVendor(*bspBranch)
	bspBranchVendorLen := len(bspBranchVendor)
	fmt.Println("bspBranchVendor=" + bspBranchVendor)
	devBranchVendor := getVendor(*devBranch)
	//devBranchVendorLen := len(devBranchVendor)
	fmt.Println("devBranchVendor=" + devBranchVendor)

	//gerrit,slash
	host := "hmdgerritserver.southeastasia.cloudapp.azure.com"
	port := "29418"
	pushPrefix := "git push ssh://" + *username + "@" + host + ":" + port + "/" + devBranchVendor + "/"

	query := "ssh -p " + port + " " + *username + "@hmdgerritserver.hmdglobal.com gerrit query "

	//output
	//bspBranchManifestFile := *repoPath + ".repo/manifests/" + *bspBranch + ".xml"
	devBranchManifestFile := *repoPath + ".repo/manifests/" + *devBranch + ".xml"
	//getManifests(bspBranchManifestFile)
	getManifests(devBranchManifestFile)

	for _, p := range gManifest.Projects {
		//fmt.Println("Tag:", p.XmlName)
		//fmt.Println("Name:", p.Name)
		//fmt.Println("Path:", p.Path)
		//fmt.Println("Revision:", p.Revision)
		//fmt.Println("InnerText:", p.InnerText)

		for _, r := range diffRecords {
			projectName := r.Project[bspBranchVendorLen+1:]
			if p.Name == projectName {
				fmt.Println("###")
				fmt.Println(strings.Replace(r.Url, "http://hmdgerritserver.southeastasia.cloudapp.azure.com", "https://hmdgerritserver.hmdglobal.com", 1))
				fmt.Println(r.Project)
				fmt.Println(r.Id)
				fmt.Println(r.Subject)
				fmt.Println("cd " + *repoPath + sep + p.Path)
				fmt.Println("git checkout -b " + *bspBranch + " " + "origin" + sep + *bspBranch + ";")
				fmt.Println("git pull -f" + ";")
				fmt.Println("git checkout -b " + *devBranch + " " + "origin" + sep + *devBranch + ";")
				fmt.Println("git pull -f" + ";")
				commitId := "$(" + query + " --current-patch-set " + r.Id + "| awk -F':' '{if($1~/revision/) print $2}')"
				fmt.Println("git cherry-pick " + commitId + ";")
				fmt.Println(pushPrefix + projectName + " " + "HEAD:refs/for/" + *devBranch + " --no-follow-tags" + ";")
				fmt.Println("###")
			}
		}
	}
}
