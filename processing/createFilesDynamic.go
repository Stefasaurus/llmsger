/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
header guard,
header guard,
enum lang values,
enum var values,
header guard,
global variable name
*/
const dynamicHeaderText = `
#ifndef %s
#define %s

typedef enum
{
	%s
	MAX_LANG_OPT,
}langs_t;

typedef enum
{
	%s,
	MAX_VAR_OPT,
}vars_t;

#ifndef _%s_C_LLMSGER
extern char* %s[MAX_LANG_OPT][MAX_VAR_OPT];
#endif


#endif
`

/*
header file name (i.e. head.h),
header guard,
lang values list i.e.

	#define EN_TEXTS_LLMSGR { "Hi" , "Good day",}
	#define FR_TEXTS_LLMSGR { "Bonjour" , "Bonne journée" ,}

global variable name
lang assignments,
*/
const dynamicSrcText = `
#include "%s"

#define _%s_C_LLMSGER // Guard from extern in header file

#define %s %s_LLMSGR

char* %s[MAX_LANG_OPT][MAX_VAR_OPT] =
{
	%s
};

/*
enum lang value,
lang values
*/`
const langAssignTemplate = `[%s] = %s_LLMSGR,`

// Struct containing all required parse information for templates
type templateInfo_t struct {
	baseName        string
	headName        string
	srcName         string
	headFilepath    string
	srcFilepath     string
	headguardDefine string
	varname         string
}

func CreateFilesDynamic(langMap map[string][]string, outDir string, varname string, basename string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing records failed: %w", err)
		}
	}()

	fileInfo, err := os.Stat(outDir)
	if err != nil {
		err = fmt.Errorf("unable to get file information for %s", outDir)
		return err
	}

	if !fileInfo.IsDir() {
		err = fmt.Errorf("output path is not a directory %s", outDir)
		return err
	}

	varNames := langMap["var"]
	if len(varNames) == 0 {
		err = fmt.Errorf("no var fields assigned in input file")
		return err
	}

	if strings.Contains(basename, ".") {
		err = fmt.Errorf("basename cannot have an extension")
		return err
	}

	if strings.TrimSpace(basename) == "" {
		err = fmt.Errorf("basename cannot be empty")
		return err
	}

	if strings.TrimSpace(varname) == "" {
		err = fmt.Errorf("varname cannot be empty")
		return err
	}

	// Create filenames
	basename = strings.TrimSpace(basename)
	basename = strings.ReplaceAll(basename, " ", "_")
	headName := basename + ".h"
	srcName := basename + ".c"
	headFilepath := filepath.Join(outDir, filepath.Base(headName))
	srcFilepath := filepath.Join(outDir, filepath.Base(srcName))
	headguardStr := strings.ToUpper(basename + "_H")

	varname = strings.TrimSpace(varname)
	varname = strings.ReplaceAll(varname, " ", "_")

	ptemplateInfo := &templateInfo_t{
		baseName:        basename,
		headName:        headName,
		srcName:         srcName,
		headFilepath:    headFilepath,
		srcFilepath:     srcFilepath,
		headguardDefine: headguardStr,
	}

	// First create the header file
	err = createHeader(langMap, ptemplateInfo)
	/*
		for k, v := range langMap {
			if k == "var" {
				continue
			}

			f, err := os.Create(newFilepath)
			if err != nil {
				return err
			}

			defer f.Close()

			headguardStr := (strings.ToUpper(strings.TrimSpace(k) + "_LANG_H"))

			writeStr := fmt.Sprintf(startText, newFilename, headguardStr, headguardStr)

			_, err = f.WriteString(writeStr)
			if err != nil {
				return err
			}

			// Actual localizations written here

			for i := 0; i < len(varNames); i++ {
				varNames[i] = strings.TrimSpace(varNames[i])
				varNames[i] = strings.ReplaceAll(varNames[i], " ", "_")
				varNames[i] = strings.ToUpper(varNames[i])
				if varNames[i] == "" {
					continue
				}

				if len(v) == i {
					break
				}

				writeStr = fmt.Sprintf(langText, varNames[i], v[i])

				_, err = f.WriteString(writeStr)
				if err != nil {
					return err
				}

			}

			// localization writing end

			writeStr = fmt.Sprintf(endText, headguardStr)

			_, err = f.WriteString(writeStr)
			if err != nil {
				return err
			}

			log.Println("Done generating:", srcName, headName)

		}
	*/
	return err
}

func createHeader(langMap map[string][]string, templateInfo *templateInfo_t) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing records failed: %w", err)
		}
	}()

	f, err := os.Create(templateInfo.headFilepath)
	if err != nil {
		return err
	}

	defer f.Close()

	//writeStr := fmt.Sprintf(startText, newFilename, headguardStr, headguardStr)

	//_, err = f.WriteString(writeStr)
	if err != nil {
		return err
	}

	fmt.Println("Done creating:", templateInfo.headFilepath)

	return err
}
