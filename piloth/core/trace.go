/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  portions taken from:
  - https://github.com/sabhiram/go-tracey
    Copyright (c) 2015 Shaba Abhiram
  - https://github.com/firnsan/file-rotator
*/

package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TraFx func(string)
type CeFx func(...interface{}) string

var (
	TRA TraFx
	CE  CeFx
)

func NewTracer(enabled bool) (TraFx, CeFx) {
	var logger *log.Logger
	// check if tracing is enabled
	if enabled {
		// if the trace folder does not exist, create it
		_, err := os.Stat(path.Join(DataPath(), "trace"))
		if err != nil {
			err = os.MkdirAll(path.Join(DataPath(), "trace"), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		// creates a file rotation logger
		fr, err := newLogFileRotator(path.Join(DataPath(), "trace", "pilot.log"))
		if err != nil {
			log.Fatal(err.Error())
		}
		logger = log.New(fr, "", log.LstdFlags)
	}
	// return tracing functions
	return newTrace(&Options{
		DisableTracing:    !enabled,
		CustomLogger:      logger,
		DisableDepthValue: false,
		DisableNesting:    false,
		SpacesPerIndent:   0,
		EnterMessage:      "",
		ExitMessage:       "",
		currentDepth:      0,
	})
}

// FileRotator It writes messages by lines limit, file size limit, or time frequency.
type FileRotator struct {
	sync.Mutex // write file order by order and  atomic incr maxLinesCurLines and maxSizeCurSize
	// The opened file
	Filename   string
	fileWriter *os.File

	// Rotate at line
	MaxLines         int
	maxLinesCurLines int

	// Rotate at size
	MaxSize        int
	maxSizeCurSize int

	// Rotate daily
	Daily         bool
	MaxDays       int64
	dailyOpenDate int

	Rotate bool

	Perm os.FileMode

	fileNameOnly, suffix string // like "project.log", project is fileNameOnly and .log is suffix
}

func newLogFileRotator(filePath string) (*FileRotator, error) {
	var err error
	w := &FileRotator{
		Filename: filepath.Clean(filePath),
		MaxLines: 1000000,
		MaxSize:  1 << 24, // 16 MB
		Daily:    false,
		MaxDays:  7,
		Rotate:   true,
		Perm:     0660,
	}

	w.suffix = filepath.Ext(w.Filename)
	w.fileNameOnly = strings.TrimSuffix(w.Filename, w.suffix)
	if w.suffix == "" {
		w.suffix = ".log"
	}

	err = w.doRotate()
	return w, err
}

// start file rotator. create file and set to locker-inside file writer.
func (w *FileRotator) startRotater() error {
	file, err := w.createFile()
	if err != nil {
		return err
	}
	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
	w.fileWriter = file
	return w.initFd()
}

func (w *FileRotator) needRotate(size int) bool {
	var day int
	if w.Daily {
		_, _, day = time.Now().Date()
	}

	return (w.MaxLines > 0 && w.maxLinesCurLines >= w.MaxLines) ||
		(w.MaxSize > 0 && w.maxSizeCurSize >= w.MaxSize) ||
		(w.Daily && day != w.dailyOpenDate)

}

// WriteMsg write bytes into file.
func (w *FileRotator) Write(b []byte) (n int, err error) {
	if w.Rotate {
		if w.needRotate(len(b)) {
			w.Lock()
			if w.needRotate(len(b)) {
				if err := w.doRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileRotator.Write: rotate failed: %s, path: %s\n", err, w.Filename)
				}
			}
			w.Unlock()
		}
	}

	w.Lock()
	n, err = w.fileWriter.Write(b)
	if err == nil {
		w.maxLinesCurLines++
		w.maxSizeCurSize += len(b)
	}
	w.Unlock()
	return
}

func (w *FileRotator) createFile() (*os.File, error) {
	// Open the file
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, w.Perm)
	return fd, err
}

func (w *FileRotator) initFd() error {
	fd := w.fileWriter
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("FileRotator.initFd: get stat err: %s\n", err)
	}
	w.maxSizeCurSize = int(fInfo.Size())
	w.dailyOpenDate = time.Now().Day()
	w.maxLinesCurLines = 0
	if fInfo.Size() > 0 {
		count, err := w.lines()
		if err != nil {
			return err
		}
		w.maxLinesCurLines = count
	}
	return nil
}

func (w *FileRotator) lines() (int, error) {
	fd, err := os.Open(w.Filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// doRotate means it need to write file in new file.
// new file name like xx.2013-01-01.log (daily) or xx.2013-01-01.001.log (by line or size)
func (w *FileRotator) doRotate() error {
	var err error
	now := time.Now()

	// Find the next available number
	num := 1
	fName := ""
	if w.MaxLines > 0 || w.MaxSize > 0 {
		for ; err == nil && num <= 9999; num++ {
			fName = w.fileNameOnly + fmt.Sprintf(".%s.%04d%s", now.Format("2006-01-02"), num, w.suffix)
			_, err = os.Lstat(fName)
		}
	} else {
		fName = fmt.Sprintf("%s.%s%s", w.fileNameOnly, now.Format("2006-01-02"), w.suffix)
		_, err = os.Lstat(fName)
	}
	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("FileRotator.doRotate: cannot find free file name number to rename %s\n", w.Filename)
	}

	// close fileWriter before rename
	if w.fileWriter != nil {
		w.fileWriter.Close()
	}

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to restart new rotator
	renameErr := os.Rename(w.Filename, fName)
	// re-start rotator
	startErr := w.startRotater()
	go w.deleteOldFiles()

	if startErr != nil {
		return fmt.Errorf("FileRotator.doRotate: restart rotator failed: %s\n", startErr)
	}
	if renameErr != nil && !os.IsNotExist(err) {
		return fmt.Errorf("FileRotator.doRotate: rename failed: %s\n", renameErr)
	}
	return nil

}

func (w *FileRotator) deleteOldFiles() {
	dir := filepath.Dir(w.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		if path == w.Filename {
			// We don't need to delete the w.Filename, because it is always up to date,
			// and the w.Filename may not exsit because some race condition
			return
		}
		if err != nil {
			// Because some race condition, the file may not exsit now
			fmt.Fprintf(os.Stderr, "FileRotator.deleteOldFiles: unable to get file info: %s, path: %s\n", err, path)
			return
		}

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.MaxDays) {
			if strings.HasPrefix(path, w.fileNameOnly) &&
				strings.HasSuffix(path, w.suffix) {
				os.Remove(path)
			}
		}
		return
	})
}

// Close Destroy the file description, close file writer.
func (w *FileRotator) Close() {
	w.fileWriter.Close()
}

// Flush file
// there are no buffering messages in file rotator in memory.
// flush file means sync file from disk.
func (w *FileRotator) Flush() {
	w.fileWriter.Sync()
}

// Define a global regex for extracting function names
var RE_stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)
var RE_detectFN = regexp.MustCompile(`\$FN`)

// Options These options represent the various settings which tracey exposes.
// A pointer to this structure is expected to be passed into the
// `newTrace(...)` function below.
type Options struct {
	// Setting "DisableTracing" to "true" will cause tracey to return
	// no-op'd functions for both exit() and enter(). The default value
	// for this is "false" which enables tracing.
	DisableTracing bool

	// Setting the "CustomLogger" to nil will cause tracey to log to
	// os.Stdout. Otherwise, this is a pointer to an object as returned
	// from `log.New(...)`.
	CustomLogger *log.Logger

	// Setting "DisableDepthValue" to "true" will cause tracey to not
	// prepend the printed function's depth to enter() and exit() messages.
	// The default value is "false", which logs the depth value.
	DisableDepthValue bool

	// Setting "DisableNesting" to "true" will cause tracey to not indent
	// any messages from nested functions. The default value is "false"
	// which enables nesting by prepending "SpacesPerIndent" number of
	// spaces per level nested.
	DisableNesting  bool
	SpacesPerIndent int `default:"2"`

	// Setting "EnterMessage" or "ExitMessage" will override the default
	// value of "Enter: " and "EXIT:  " respectively.
	EnterMessage string `default:"ENTER: "`
	ExitMessage  string `default:"EXIT:  "`

	// Private member, used to keep track of how many levels of nesting
	// the current trace functions have navigated.
	currentDepth int
}

// newTrace Main entry-point for the tracey lib. Calling New with nil will
// result in the default options being used.
func newTrace(opts *Options) (func(string), func(...interface{}) string) {
	var options Options
	if opts != nil {
		options = *opts
	}

	// If tracing is not enabled, just return no-op functions
	if options.DisableTracing {
		return func(string) {}, func(...interface{}) string { return "" }
	}

	// Revert to stdout if no logger is defined
	if options.CustomLogger == nil {
		options.CustomLogger = log.New(os.Stdout, "", 0)
	}

	// Use reflect to deduce "default" values for the
	// Enter and Exit messages (if they are not set)
	reflectedType := reflect.TypeOf(options)
	if options.EnterMessage == "" {
		field, _ := reflectedType.FieldByName("EnterMessage")
		options.EnterMessage = field.Tag.Get("default")
	}
	if options.ExitMessage == "" {
		field, _ := reflectedType.FieldByName("ExitMessage")
		options.ExitMessage = field.Tag.Get("default")
	}

	// If nesting is enabled, and the spaces are not specified,
	// use the "default" value
	if options.DisableNesting {
		options.SpacesPerIndent = 0
	} else if options.SpacesPerIndent == 0 {
		field, _ := reflectedType.FieldByName("SpacesPerIndent")
		options.SpacesPerIndent, _ = strconv.Atoi(field.Tag.Get("default"))
	}

	//
	// Define functions we will use and return to the caller
	//
	_spacify := func() string {
		spaces := strings.Repeat(" ", options.currentDepth*options.SpacesPerIndent)
		if !options.DisableDepthValue {
			return fmt.Sprintf("[%2d]%s", options.currentDepth, spaces)
		}
		return spaces
	}

	// Increment function to increase the current depth value
	_incrementDepth := func() {
		options.currentDepth += 1
	}

	// Decrement function to decrement the current depth value
	//  + panics if current depth value is < 0
	_decrementDepth := func() {
		options.currentDepth -= 1
		if options.currentDepth < 0 {
			panic("Depth is negative! Should never happen!")
		}
	}

	// Enter function, invoked on function entry
	_enter := func(args ...interface{}) string {
		defer _incrementDepth()

		// Figure out the name of the caller and use that
		fnName := "<unknown>"
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fnName = RE_stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
		}

		traceMessage := fnName
		if len(args) > 0 {
			if fmtStr, ok := args[0].(string); ok {
				// We have a string leading args, assume its to be formatted
				traceMessage = fmt.Sprintf(fmtStr, args[1:]...)
			}
		}

		// "$FN" will be replaced by the name of the function (if present)
		traceMessage = RE_detectFN.ReplaceAllString(traceMessage, fnName)

		options.CustomLogger.Printf("%s%s%s\n", _spacify(), options.EnterMessage, traceMessage)
		return traceMessage
	}

	// Exit function, invoked on function exit (usually deferred)
	_exit := func(s string) {
		_decrementDepth()
		options.CustomLogger.Printf("%s%s%s\n", _spacify(), options.ExitMessage, s)
	}

	return _exit, _enter
}
