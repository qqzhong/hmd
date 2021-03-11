// @filename           :  diff_vendor_bsp.go
// @author             :  church.zhong@hmdglobal.com
// @date               :  Sat Mar  6 10:23:14 HKT 2021
// @function           :  diff codes of vendor and bsp then output.
// @see                :  https://golang.org/pkg/io/ioutil/#ReadDir
// @require            :  golang 1.16
package main

import (
	//"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var sep string = string(os.PathSeparator)
var BUFFERSIZE = 8388608 // 8*1024*1024
var gBspPath string
var gBspPathLength int
var gQcomPath string
var gQcomPathLength int

type FileMd5MapType map[string]string

var gBspModifiedOutput string
var gQcomBaselineOutput string

// Free memory
var gBspFileMd5Map FileMd5MapType = make(FileMd5MapType)

type scanFunc func(string)

func getPathFilename(v string) (string, string, string, string) {
	d, filename, f, extension := "", "", "", ""
	if 0 != len(strings.TrimSpace(v)) {
		//filename = filepath.Base(v)
		d, filename = filepath.Split(v)
		extension = path.Ext(v)
		f = filename[0 : len(filename)-len(extension)]
	}
	//fmt.Printf("d=%s,filename=%s,f=%s,extension=%s\n", d, filename, f, extension)
	return d, filename, f, extension
}

func getFileMode(path string) (bool, bool) {
	dir, file := false, false
	fi, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
		return dir, file
	}
	// 0400, 0777, etc.
	//fmt.Printf("permissions: %#o\n", fi.Mode().Perm())
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		//fmt.Println("regular file")
		file = true
	case mode.IsDir():
		//fmt.Println("directory")
		dir = true
	case mode&os.ModeSymlink != 0:
		//fmt.Println("symbolic link")
	case mode&os.ModeNamedPipe != 0:
		//fmt.Println("named pipe")
	}
	return dir, file
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

func getFileMd5sum(filePath string) (string, error) {
	md5sum := ""
	file, err := os.Open(filePath)
	if err != nil {
		return md5sum, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5sum, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5sum = hex.EncodeToString(hashInBytes)
	return md5sum, nil
}

func copyFile(srcPath string, dstPath string, name string) {
	//fmt.Println("srcPath=" + srcPath)
	//fmt.Println("dstPath=" + dstPath)
	src := srcPath + sep + name
	dst := dstPath + sep + name
	//fmt.Println("src=" + src)
	//fmt.Println("dst=" + dst)

	dstDir, _, _, _ := getPathFilename(dst)
	if exist, _ := pathExists(dstDir); !exist {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	source, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer destination.Close()

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
			return
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			log.Fatal(err)
			return
		}
	}
	//fmt.Println("copy:" + dst)
}

func scanAndroid(path string, fun scanFunc) {
	//"""Scan android"""
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()
		entryPath := path + sep + fileName
		isDir, isFile := getFileMode(entryPath)
		if isDir && fileName != ".git" && fileName != ".svn" && fileName != "CVS" {
			scanAndroid(entryPath, fun)
		} else if isFile {
			//fmt.Println(entryPath)
			if exist, _ := pathExists(entryPath); exist {
				fun(entryPath)
			}
		} else {
			//pass
		}
	}
}

func bspFunc(entryPath string) {
	name := entryPath[gBspPathLength:]
	md5sum, _ := getFileMd5sum(entryPath)
	gBspFileMd5Map[name] = md5sum
}

func scanBsp(path string) {
	scanAndroid(path, bspFunc)
}

func qcomFunc(entryPath string) {
	name := entryPath[gQcomPathLength:]
	md5sum, _ := getFileMd5sum(entryPath)

	if _, ok := gBspFileMd5Map[name]; ok {
		if gBspFileMd5Map[name] != md5sum {
			//fmt.Printf("%s=%s\n", name, gBspFileMd5Map[name])
			copyFile(gBspPath, gBspModifiedOutput, name)
			copyFile(gQcomPath, gQcomBaselineOutput, name)
			fmt.Println("copy:" + name)
		}
		gBspFileMd5Map[name] = ""
	}
}

func scanQcom(path string) {
	scanAndroid(path, qcomFunc)
}

// Prints how much time has elapsed after main() exit.
func elapsedTime(start time.Time) {
	fmt.Println(time.Since(start))
}

func main() {
	defer elapsedTime(time.Now())

	var bspPath = flag.String("bsp", "bspPath", "BSP android path")
	var qcomPath = flag.String("qcom", "qcomPath", "Qcom android path")
	var android = flag.String("list", "build,device,frameworks,hardware,kernel,vendor", "list of android path")
	flag.Parse()
	fmt.Println("bspPath=" + *bspPath)
	fmt.Println("qcomPath=" + *qcomPath)
	fmt.Println("android=" + *android)
	// we're happy if android source are given.
	list := strings.Split(*android, ",")
	androidSrc := make(map[string]int, 8)
	for _, f := range list {
		key := strings.Trim(f, " ")
		if len(key) != 0 {
			if _, ok := androidSrc[key]; !ok {
				fmt.Println(key)
				androidSrc[key] = 0
			}
		}
	}

	if exist, _ := pathExists(*bspPath); !exist {
		fmt.Printf("no such file or directory:bspPath=%32s\n", *bspPath)
		return
	}
	if exist, _ := pathExists(*qcomPath); !exist {
		fmt.Printf("no such file or directory:qcomPath=%32s\n", *qcomPath)
		return
	}

	gBspPath = *bspPath
	//fmt.Println("gBspPath=" + gBspPath)
	gBspPathLength = len(gBspPath)
	gQcomPath = *qcomPath
	//fmt.Println("gQcomPath=" + gQcomPath)
	gQcomPathLength = len(gQcomPath)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	now := time.Now()
	moment := fmt.Sprintf("%04d_%02d%02d_%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	extract := pwd + sep + "extract_" + moment
	fmt.Println("diff=" + extract)
	gBspModifiedOutput = extract + sep + "bspModified"
	//fmt.Println("gBspModifiedOutput=" + gBspModifiedOutput)
	err = os.MkdirAll(gBspModifiedOutput, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}
	gQcomBaselineOutput = extract + sep + "qcomBaseline"
	//fmt.Println("gQcomBaselineOutput=" + gQcomBaselineOutput)
	err = os.MkdirAll(gQcomBaselineOutput, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}

	for f, _ := range androidSrc {
		scanBsp(*bspPath + sep + f)
	}
	for f, _ := range androidSrc {
		scanQcom(*qcomPath + sep + f)
	}

	// BSP added something
	for key, value := range gBspFileMd5Map {
		if value != "" {
			copyFile(gBspPath, gBspModifiedOutput, key)
			fmt.Println("BSP added=" + key)
		}
	}

	//gBspFileMd5Map = make(FileMd5MapType)
}
