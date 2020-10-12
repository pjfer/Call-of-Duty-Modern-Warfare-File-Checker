package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tadvi/winc"
)

func fileToMD5(filename *string) (string, error) {
	f, err := os.Open(*filename)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer f.Close()

	h := md5.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func computeFolder(folderPath *string, statusReporter *winc.MultiEdit, progressBarPercentage *winc.Edit) (map[string]string, error) {
	files := []string{}
	filesAsMD5Values := make(map[string]string)

	statusReporter.AddLine("Hashing the files!")

	err := filepath.Walk(*folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				statusReporter.AddLine(err.Error())
				return err
			}

			if !info.IsDir() {
				files = append(files, path)
			}

			return nil
		})

	if err != nil {
		statusReporter.AddLine(err.Error())
		return nil, err
	}

	for idx := range files {
		start := time.Now()

		if res, err := fileToMD5(&files[idx]); err == nil {
			filesAsMD5Values[files[idx]] = res
			t := time.Now()
			elapsed := t.Sub(start)
			statusReporter.AddLine(fmt.Sprintf("File %s converted in %s", files[idx], elapsed.String()))
			progressBarPercentage.SetText(fmt.Sprintf("%d%%", ((idx+1)*100)/len(files)))
		} else {
			statusReporter.AddLine(err.Error())
			return nil, err
		}
	}

	statusReporter.AddLine("Folder hashing complete!")

	return filesAsMD5Values, nil
}

func compareMD5Values(myMD5Values map[string]string, rightMD5Values map[string]string, statusReporter *winc.MultiEdit, progressBarPercentage *winc.Edit) map[string]string {
	i := 1
	diffFiles := make(map[string]string)

	statusReporter.AddLine("Comparing the hashes!")

	if len(myMD5Values) == len(rightMD5Values) {
		for path, MD5Value := range rightMD5Values {
			if myMD5Values[path] != MD5Value {
				diffFiles[path] = MD5Value
			}

			progressBarPercentage.SetText(fmt.Sprintf("%d%%", (i*100)/len(rightMD5Values)))
			i++
		}
	} else {
		for path, MD5Value := range rightMD5Values {
			if _, found := myMD5Values[path]; !found {
				diffFiles[path] = MD5Value
			}

			progressBarPercentage.SetText(fmt.Sprintf("%d%%", (i*100)/len(rightMD5Values)))
			i++
		}
	}

	statusReporter.AddLine("MD5 hashes successful compared!")

	return diffFiles
}

func readMD5Values(filename *string, statusReporter *winc.MultiEdit) (map[string]string, error) {
	file, err := os.Open(*filename)

	statusReporter.AddLine("Writing the hashes into the file!")

	if err != nil {
		statusReporter.AddLine(err.Error())
		return nil, err
	}

	defer file.Close()

	rightMD5Values := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		MD5Value := strings.Split(scanner.Text(), ",")
		rightMD5Values[MD5Value[0]] = MD5Value[1]
	}

	if err := scanner.Err(); err != nil {
		statusReporter.AddLine(err.Error())
		return nil, err
	}

	statusReporter.AddLine("MD5 hashes successfully read from the file!")

	return rightMD5Values, nil
}

func writeMD5Values(filename *string, myMD5Values map[string]string, statusReporter *winc.MultiEdit, progressBarPercentage *winc.Edit) error {
	file, err := os.Create(*filename)

	if err != nil {
		statusReporter.AddLine(err.Error())
		return err
	}

	defer file.Close()

	i := 1

	for path, MD5Value := range myMD5Values {
		file.WriteString(path)
		file.WriteString(",")
		file.WriteString(MD5Value)
		file.WriteString("\n")
		progressBarPercentage.SetText(fmt.Sprintf("%d%%", (i*100)/len(myMD5Values)))
		i++
	}

	statusReporter.AddLine("MD5 hashes successfully written to the file!")

	return nil
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}

func displayGUI() {
	var myMD5Values map[string]string

	font := winc.NewFont("Times New Roman", 12, 0)

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(800, 600)
	mainWindow.SetText("Call of Duty Modern Warfare File Checker")
	mainWindow.EnableSizable(false)
	mainWindow.EnableMaxButton(false)

	progressBar := winc.NewProgressBar(mainWindow)
	progressBar.SetPos(180, 120)
	progressBar.SetSize(350, 20)
	progressBarX, progressBarY := progressBar.Pos()
	progressBarWidth, progressBarHeight := progressBar.Size()

	progressBarPercentage := winc.NewEdit(mainWindow)
	progressBarPercentage.SetPos(progressBarWidth+progressBarX+5, progressBarY)
	progressBarPercentage.SetSize(40, progressBarHeight)
	progressBarPercentage.SetReadOnly(true)
	progressBarPercentage.SetText("0%")
	progressBarPercentage.OnChange().Bind(func(e *winc.Event) {
		percentage := progressBarPercentage.Text()
		num, _ := strconv.Atoi(percentage[:len(percentage)-1])
		progressBar.SetValue(num)
	})

	statusReporterLabel := winc.NewLabel(mainWindow)
	statusReporterLabel.SetPos(40, 240)
	statusReporterLabel.SetSize(300, 20)
	statusReporterLabel.SetFont(font)
	statusReporterLabel.SetText("Messages about the program's status")
	statusReporterLabelX, statusReporterLabelY := statusReporterLabel.Pos()
	_, statusReporterLabelHeight := statusReporterLabel.Size()

	statusReporter := winc.NewMultiEdit(mainWindow)
	statusReporter.SetPos(statusReporterLabelX, statusReporterLabelY+statusReporterLabelHeight)
	statusReporter.SetSize(700, 120)
	statusReporter.SetReadOnly(true)
	statusReporter.SetFont(font)
	statusReporterWidth, statusReporterHeight := statusReporter.Size()

	resultsReporterLabel := winc.NewLabel(mainWindow)
	resultsReporterLabel.SetPos(statusReporterLabelX, 400)
	resultsReporterLabel.SetSize(300, statusReporterLabelHeight)
	resultsReporterLabel.SetFont(font)
	resultsReporterLabel.SetText("Messages about the files' hashes values")
	resultsReporterLabelX, resultsReporterLabelY := resultsReporterLabel.Pos()
	_, resultsReporterLabelHeight := resultsReporterLabel.Size()

	resultsReporter := winc.NewMultiEdit(mainWindow)
	resultsReporter.SetPos(resultsReporterLabelX, resultsReporterLabelY+resultsReporterLabelHeight)
	resultsReporter.SetSize(statusReporterWidth, statusReporterHeight)
	resultsReporter.SetReadOnly(true)
	resultsReporter.SetFont(font)

	folderLabel := winc.NewLabel(mainWindow)
	folderLabel.SetPos(20, 20)
	folderLabel.SetSize(320, 20)
	folderLabel.SetFont(font)
	folderLabel.SetText("Call of Duty Modern Warfare Game folder")
	folderLabelX, folderLabelY := folderLabel.Pos()
	folderLabelWidth, folderLabelHeight := folderLabel.Size()

	saveFileLabel := winc.NewLabel(mainWindow)
	saveFileLabel.SetPos(folderLabelX, folderLabelY+folderLabelHeight+10)
	saveFileLabel.SetSize(folderLabelWidth, folderLabelHeight)
	saveFileLabel.SetFont(font)
	saveFileLabel.SetText("File to write the hashes obtained from your files")
	saveFileLabelX, saveFileLabelY := saveFileLabel.Pos()
	_, saveFileLabelHeight := saveFileLabel.Size()

	readFileLabel := winc.NewLabel(mainWindow)
	readFileLabel.SetPos(saveFileLabelX, saveFileLabelY+saveFileLabelHeight+10)
	readFileLabel.SetSize(folderLabelWidth, saveFileLabelHeight)
	readFileLabel.SetFont(font)
	readFileLabel.SetText("File to read the hashes of the functional files")
	_, readFileLabelY := readFileLabel.Pos()
	_, readFileLabelHeight := readFileLabel.Size()

	folderPath := winc.NewEdit(mainWindow)
	folderPath.SetReadOnly(true)
	folderPath.SetPos(folderLabelWidth+folderLabelX, folderLabelY)
	folderPath.SetSize(400, folderLabelHeight)
	folderPath.SetFont(font)
	folderPathWidth, _ := folderPath.Size()

	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		statusReporter.AddLine(err.Error())
		winc.MsgBoxOk(mainWindow, "ERROR", err.Error())
		return
	}

	saveFilePath := winc.NewEdit(mainWindow)
	saveFilePath.SetReadOnly(true)
	saveFilePath.SetPos(folderLabelWidth+folderLabelX, saveFileLabelY)
	saveFilePath.SetSize(folderPathWidth, saveFileLabelHeight)
	saveFilePath.SetFont(font)
	saveFilePath.SetText(currentDir + "\\results\\myMD5Values.txt")

	readFilePath := winc.NewEdit(mainWindow)
	readFilePath.SetReadOnly(true)
	readFilePath.SetPos(folderLabelWidth+folderLabelX, readFileLabelY)
	readFilePath.SetSize(folderPathWidth, readFileLabelHeight)
	readFilePath.SetFont(font)
	readFilePath.SetText(currentDir + "\\results\\rightMD5Values.txt")

	hashButton := winc.NewPushButton(mainWindow)
	saveButton := winc.NewPushButton(mainWindow)
	compareButton := winc.NewPushButton(mainWindow)
	fullRunButton := winc.NewPushButton(mainWindow)

	folderBrowserButton := winc.NewPushButton(mainWindow)
	folderBrowserButton.SetText("...")
	folderBrowserButton.SetPos(folderLabelWidth+folderPathWidth+folderLabelX+5, folderLabelY)
	folderBrowserButton.SetSize(20, folderLabelHeight)
	folderBrowserButton.OnClick().Bind(func(e *winc.Event) {
		if folder, accepted := winc.ShowBrowseFolderDlg(mainWindow, "Select Folder"); accepted {
			folderPath.SetText(string(folder))

			hashButton.SetEnabled(true)
			saveButton.SetEnabled(false)
			compareButton.SetEnabled(false)
			fullRunButton.SetEnabled(true)
		}
	})

	saveFileFolderBrowserButton := winc.NewPushButton(mainWindow)
	saveFileFolderBrowserButton.SetText("...")
	saveFileFolderBrowserButton.SetPos(folderLabelWidth+folderPathWidth+folderLabelX+5, saveFileLabelY)
	saveFileFolderBrowserButton.SetSize(20, saveFileLabelHeight)
	saveFileFolderBrowserButton.OnClick().Bind(func(e *winc.Event) {
		if folder, accepted := winc.ShowOpenFileDlg(mainWindow, "Select the file to write the MD5 values", "", 0, "C:\\"); accepted {
			saveFilePath.SetText(string(folder))
		}
	})

	readFileFolderBrowserButton := winc.NewPushButton(mainWindow)
	readFileFolderBrowserButton.SetText("...")
	readFileFolderBrowserButton.SetPos(folderLabelWidth+folderPathWidth+folderLabelX+5, readFileLabelY)
	readFileFolderBrowserButton.SetSize(20, readFileLabelHeight)
	readFileFolderBrowserButton.OnClick().Bind(func(e *winc.Event) {
		if folder, accepted := winc.ShowOpenFileDlg(mainWindow, "Select the file to read the MD5 values", "", 0, "C:\\"); accepted {
			readFilePath.SetText(string(folder))
		}
	})

	hashButton.SetFont(font)
	hashButton.SetText("Hash the files")
	hashButton.SetPos(140, 160)
	hashButton.SetSize(150, 30)
	hashButton.SetEnabled(false)
	hashButtonX, hashButtonY := hashButton.Pos()
	hashButtonWidth, hashButtonHeight := hashButton.Size()
	hashButton.OnClick().Bind(func(e *winc.Event) {
		hashButton.SetEnabled(false)
		saveButton.SetEnabled(false)
		compareButton.SetEnabled(false)
		fullRunButton.SetEnabled(false)

		go func() {
			start := time.Now()
			folder := folderPath.Text()
			folder = strings.TrimSuffix(folder, "\r\n")
			out, err := computeFolder(&folder, statusReporter, progressBarPercentage)

			if err != nil {
				statusReporter.AddLine(err.Error())
				winc.MsgBoxOk(mainWindow, "ERROR", err.Error())
			} else {
				myMD5Values = out

				resultsReporter.AddLine("Hashes generated from your game files:")

				for path, MD5Value := range myMD5Values {
					resultsReporter.AddLine(fmt.Sprintf("%s: %s", path, MD5Value))
				}

				resultsReporter.AddLine("_" + strings.Repeat("_", 100))

				t := time.Now()
				elapsed := t.Sub(start)

				winc.MsgBoxOk(mainWindow, "Info", "Folder hashing completed in "+elapsed.String()+"!")
			}

			hashButton.SetEnabled(true)
			saveButton.SetEnabled(true)
			compareButton.SetEnabled(true)
			fullRunButton.SetEnabled(true)
		}()
	})

	saveButton.SetFont(font)
	saveButton.SetText("Save the hashes")
	saveButton.SetPos(hashButtonX+hashButtonWidth+10, hashButtonY)
	saveButton.SetSize(hashButtonWidth, hashButtonHeight)
	saveButton.SetEnabled(false)
	saveButtonX, saveButtonY := saveButton.Pos()
	saveButtonWidth, saveButtonHeight := saveButton.Size()
	saveButton.OnClick().Bind(func(e *winc.Event) {
		hashButton.SetEnabled(false)
		saveButton.SetEnabled(false)
		compareButton.SetEnabled(false)
		fullRunButton.SetEnabled(false)

		go func() {
			start := time.Now()
			filename := saveFilePath.Text()
			filename = strings.TrimSuffix(filename, "\r\n")
			writeMD5Values(&filename, myMD5Values, statusReporter, progressBarPercentage)

			t := time.Now()
			elapsed := t.Sub(start)

			winc.MsgBoxOk(mainWindow, "Info", "Hashes successfully saved in "+elapsed.String()+"!")

			hashButton.SetEnabled(true)
			saveButton.SetEnabled(true)
			compareButton.SetEnabled(true)
			fullRunButton.SetEnabled(true)
		}()
	})

	compareButton.SetFont(font)
	compareButton.SetText("Compare the hashes")
	compareButton.SetPos(saveButtonX+saveButtonWidth+10, saveButtonY)
	compareButton.SetSize(saveButtonWidth, saveButtonHeight)
	compareButton.SetEnabled(false)
	compareButton.OnClick().Bind(func(e *winc.Event) {
		hashButton.SetEnabled(false)
		saveButton.SetEnabled(false)
		compareButton.SetEnabled(false)
		fullRunButton.SetEnabled(false)

		go func() {
			start := time.Now()
			filename := readFilePath.Text()
			filename = strings.TrimSuffix(filename, "\r\n")
			rightMD5Values, err := readMD5Values(&filename, statusReporter)

			if err != nil {
				statusReporter.AddLine(err.Error())
				winc.MsgBoxOk(mainWindow, "ERROR", err.Error())
			} else {
				faultyMD5Values := compareMD5Values(myMD5Values, rightMD5Values, statusReporter, progressBarPercentage)

				resultsReporter.AddLine("Files that might be corrupted and their respective hashes:")

				for path, MD5Value := range faultyMD5Values {
					resultsReporter.AddLine(fmt.Sprintf("%s: %s", path, MD5Value))
				}

				resultsReporter.AddLine("_" + strings.Repeat("_", 100))

				t := time.Now()
				elapsed := t.Sub(start)

				winc.MsgBoxOk(mainWindow, "Info", "Hash successfully compared in "+elapsed.String()+"!")
			}

			hashButton.SetEnabled(true)
			saveButton.SetEnabled(true)
			compareButton.SetEnabled(true)
			fullRunButton.SetEnabled(true)
		}()
	})

	fullRunButton.SetFont(font)
	fullRunButton.SetText("Full Run")
	fullRunButton.SetPos(hashButtonX+hashButtonWidth+10, hashButtonY+hashButtonHeight+10)
	fullRunButton.SetSize(hashButtonWidth, hashButtonHeight)
	fullRunButton.SetEnabled(false)
	fullRunButton.OnClick().Bind(func(e *winc.Event) {
		hashButton.SetEnabled(false)
		saveButton.SetEnabled(false)
		compareButton.SetEnabled(false)
		fullRunButton.SetEnabled(false)

		go func() {
			start := time.Now()
			folder := folderPath.Text()
			folder = strings.TrimSuffix(folder, "\r\n")
			out, err := computeFolder(&folder, statusReporter, progressBarPercentage)

			if err != nil {
				statusReporter.AddLine(err.Error())
				winc.MsgBoxOk(mainWindow, "ERROR", err.Error())
			} else {
				myMD5Values = out

				resultsReporter.AddLine("Hashes generated from your game files:")

				for path, MD5Value := range myMD5Values {
					resultsReporter.AddLine(fmt.Sprintf("%s: %s", path, MD5Value))
				}

				resultsReporter.AddLine("_" + strings.Repeat("_", 100))

				filename := saveFilePath.Text()
				filename = strings.TrimSuffix(filename, "\r\n")
				writeMD5Values(&filename, myMD5Values, statusReporter, progressBarPercentage)

				filename = readFilePath.Text()
				filename = strings.TrimSuffix(filename, "\r\n")
				rightMD5Values, err := readMD5Values(&filename, statusReporter)

				if err != nil {
					statusReporter.AddLine(err.Error())
					winc.MsgBoxOk(mainWindow, "ERROR", err.Error())
				} else {
					faultyMD5Values := compareMD5Values(myMD5Values, rightMD5Values, statusReporter, progressBarPercentage)

					resultsReporter.AddLine("Files that might be corrupted and their respective hashes:")

					for path, MD5Value := range faultyMD5Values {
						resultsReporter.AddLine(fmt.Sprintf("%s: %s", path, MD5Value))
					}

					resultsReporter.AddLine("_" + strings.Repeat("_", 100))

					t := time.Now()
					elapsed := t.Sub(start)

					winc.MsgBoxOk(mainWindow, "Info", "Full Run successfully completed in "+elapsed.String()+"!")
				}
			}

			hashButton.SetEnabled(true)
			saveButton.SetEnabled(true)
			compareButton.SetEnabled(true)
			fullRunButton.SetEnabled(true)
		}()
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop()
}

func main() {
	displayGUI()
}
