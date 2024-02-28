/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package process

import (
	"fmt"
	"strconv"
	"strings"
)

//var replaceMap = map[string][]string{}

func ReplaceFields(langMap map[string][]string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing records failed: %w", err)
		}
	}()

	fromFields, fromOk := langMap["from"]
	toasciiField, toasciiOk := langMap["toascii"]

	if !fromOk && !toasciiOk {
		return err
	} else if fromOk != toasciiOk {
		err = fmt.Errorf("either \"from\" or \"toascii\" has been used while the other was not defined as a main column name")
		return err
	}

	type rep_t struct {
		original        string
		replacementCode int
	}

	replacements := make([]rep_t, 0, 20)

	for idx, original := range fromFields {

		if len(strings.TrimSpace(original)) == 0 {
			continue
		}
		if len(strings.TrimSpace(toasciiField[idx])) == 0 {
			continue
		}

		asciiCode, err := strconv.Atoi(toasciiField[idx])
		if err != nil {
			err = fmt.Errorf("could not convert replacement field to ASCII code: %s", toasciiField[idx])
			return err
		}

		if (asciiCode > 255) || (asciiCode < 0) {
			err = fmt.Errorf("toascii file contains a non ASCII code: %s", toasciiField[idx])
			return err
		}

		replacements = append(replacements, rep_t{original: original, replacementCode: asciiCode})

	}

	fmt.Println("Performing replace...")

	delete(langMap, "from")
	delete(langMap, "toascii")

	for k, v := range langMap {
		if k == "var" {
			continue
		}
		for idx, text := range v {
			for _, replacement := range replacements {
				v[idx] = strings.ReplaceAll(text, replacement.original, string(rune(replacement.replacementCode)))
			}
		}
	}

	//retMap := make(map[string][]string, len(repFields))
	// later delete these fields from the langMap
	return nil
}
