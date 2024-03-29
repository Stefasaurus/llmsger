/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

/*
header guard,
header guard,
enum lang values,
uppercase basename,
basename,
enum var values,
uppercase basename,
basename
header guard,
global variable name,
uppercase basename,
uppercase basename,
*/
const dynamicHeaderText = `
#ifndef %s
#define %s

typedef enum
{
%s
	MAX_%s_LANG_OPT,
}%s_langs_t;

typedef enum
{
%s
	MAX_%s_VAR_OPT,
}%s_vars_t;

#ifndef _%s_C_LLMSGER_
extern char* %s[MAX_%s_LANG_OPT][MAX_%s_VAR_OPT];
#endif


#endif
`

/*
header file name (i.e. head.h),
header guard,
lang values list i.e.

	#define EN_TEXTS_LLMSGR { "Hi" , "Good day",}
	#define FR_TEXTS_LLMSGR { "Bonjour" , "Bonne journée" ,}

global variable name,
uppercase basename,
uppercase basename,
lang assignments,
*/
const dynamicSrcText = `
#include "%s"

#define _%s_C_LLMSGER_ // Guard from extern in header file

%s

char* %s[MAX_%s_LANG_OPT][MAX_%s_VAR_OPT] =
{
%s
};
`

// Struct containing all required parse information for templates

type langInfo_t struct {
	langVars     []string
	langOpts     []string
	langOptEnums []string
	langOptTexts [][]string // index correlates with langOpts and langOptEnums
}
type templateInfo_t struct {
	baseName        string
	headName        string
	srcName         string
	headFilepath    string
	srcFilepath     string
	headguardDefine string
	varname         string
	langInfo_t
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
		varname:         varname,
	}
	ptemplateInfo.langOpts = make([]string, 0, 100)
	ptemplateInfo.langOptEnums = make([]string, 0, 100)
	ptemplateInfo.langOptTexts = make([][]string, 0, 100)

	for k, v := range langMap {
		if k == "var" {
			ptemplateInfo.langVars = v // Get user variables names
			continue
		}

		ptemplateInfo.langOpts = append(ptemplateInfo.langOpts, k)

		temp := strings.ToUpper(strings.ReplaceAll(k, " ", "_"))
		ptemplateInfo.langOptEnums = append(ptemplateInfo.langOptEnums, temp)

	}

	for _, opt := range ptemplateInfo.langOpts {
		ptemplateInfo.langOptTexts = append(ptemplateInfo.langOptTexts, langMap[opt]) // Get the translated texts
	}
	// First create the header file
	err = createHeader(langMap, ptemplateInfo)
	err = createSrc(langMap, ptemplateInfo)

	return err
}

func createHeader(langMap map[string][]string, templateInfo *templateInfo_t) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("creating header (.h) file failed: %w", err)
		}
	}()

	f, err := os.Create(templateInfo.headFilepath)
	if err != nil {
		return err
	}

	defer f.Close()

	// creation
	var langOptEnumText string = ""
	var VarEnumText string = ""

	for _, v := range templateInfo.langOptEnums {
		langOptEnumText = langOptEnumText + fmt.Sprintf("\t%s,\n", v)
	}

	for _, v := range templateInfo.langVars {
		VarEnumText = VarEnumText + fmt.Sprintf("\t%s,\n", v)
	}

	writeStr := fmt.Sprintf(dynamicHeaderText,
		templateInfo.headguardDefine,
		templateInfo.headguardDefine,
		langOptEnumText,
		strings.ToUpper(templateInfo.baseName),
		templateInfo.baseName,
		VarEnumText,
		strings.ToUpper(templateInfo.baseName),
		templateInfo.baseName,
		templateInfo.headguardDefine,
		templateInfo.varname,
		strings.ToUpper(templateInfo.baseName),
		strings.ToUpper(templateInfo.baseName))

	_, err = f.WriteString(writeStr)
	if err != nil {
		return err
	}
	// creation done

	fmt.Println("Done creating:", templateInfo.headFilepath)

	return err
}

func createSrc(langMap map[string][]string, templateInfo *templateInfo_t) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("creating source (.c) file failed: %w", err)
		}
	}()

	f, err := os.Create(templateInfo.srcFilepath)
	if err != nil {
		return err
	}

	defer f.Close()

	// creation

	var wg sync.WaitGroup
	results := make(chan string, len(templateInfo.langOptTexts))
	//fmt.Println("Processing:", templateInfo.srcFilepath)
	bar := progressbar.NewOptions((len(templateInfo.langOptTexts)),
		progressbar.OptionSetDescription("Processing langs:"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetDescription("[cyan] Processing langs...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for idx, v := range templateInfo.langOptTexts {
		wg.Add(1)
		go func(idx int, v []string) {
			defer wg.Done()
			//log.Println("Processing lang option:", idx+1)
			langValuesDefinitionsText := fmt.Sprintf("#define %s_TEXTS_LLMSGR ", templateInfo.langOptEnums[idx])
			langValuesDefinitionsText += "{ "
			for subidx, str := range v {
				langValuesDefinitionsText += fmt.Sprintf("\"%s\" /*%d*/, ", str, subidx+1)
			}
			langValuesDefinitionsText = langValuesDefinitionsText + "}\n\n"
			results <- langValuesDefinitionsText
			bar.Add(1)
		}(idx, v)
	}

	wg.Wait()
	close(results)
	var langValuesDefinitionsTextComb string = "" //for combining all the results
	for result := range results {
		langValuesDefinitionsTextComb += result
	}

	var varDefinitions string
	for idx, v := range templateInfo.langOptEnums {
		varDefinitions += fmt.Sprintf("\t[%s] = %s_TEXTS_LLMSGR,\n", v, templateInfo.langOptEnums[idx])
	}

	//for _, v := range templateInfo.langOptEnums {
	//	VarEnumText = VarEnumText + fmt.Sprintf("\t#define %s_TEXTS_LLMSGR ,\n", v)
	//}
	//log.Print(VarEnumText)

	writeStr := fmt.Sprintf(dynamicSrcText,
		templateInfo.headName,
		templateInfo.headguardDefine,
		langValuesDefinitionsTextComb,
		templateInfo.varname,
		strings.ToUpper(templateInfo.baseName),
		strings.ToUpper(templateInfo.baseName),
		varDefinitions)

	_, err = f.WriteString(writeStr)
	if err != nil {
		return err
	}
	// creation done

	fmt.Println("Done creating:", templateInfo.srcFilepath)

	return err
}
