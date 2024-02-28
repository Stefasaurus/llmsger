<h1 align="center">
  llmsger
</h2>
<h3 align="center">
  A CSV localization and text management CLI tool.
</h3>
<p align="center">
  <a href="https://github.com/Stefasaurus/llmsger" target="_blank" rel="noopener noreferrer"><img src="https://img.shields.io/badge/GitHub-Source-success" alt="View project on GitHub" /></a>&nbsp;
</p>

## Overview

This CLI tool allows you to generate files recquired for the localization of strings in a
project (currently only C projects), all from a CSV file. 


### Features include:

- Creation of C files for the purpose of localization
- Merging localization files
- Customizing many aspects of auto-generated code
- Creating dynamically usable localization files, which enable using all localized strings in runtime!
- Ability to replace all user defined UTF-8 characters with ASCII codes
<!-- [lock:donate] 🚫--------------------------------------- -->


## Table of Contents

- [Overview](#overview)
  - [Features include:](#features-include)
- [Table of Contents](#table-of-contents)
- [Clone the project](#clone-the-project)
- [Install dependencies](#install-dependencies)
- [Basic usage](#basic-usage)
- [Advanced usage](#advanced-usage)
- [Using the UTF-8 to ASCII code replacer](#using-the-utf-8-to-ascii-code-replacer)
- [Contributing](#contributing)
- [⭐ Found It Helpful? Star It!](#-found-it-helpful-star-it)
- [License](#license)

## Clone the project

```bash
git clone https://github.com/Stefasaurus/llmsger.git
```

## Install dependencies

To be able to buld the project, you must install [Go](https://go.dev/doc/install).\
Then using *make*, run the following in the directory of the Makefile:
```bash
make
```
For windows, it is recommended to run make using git bash.

## Basic usage

Preparation:
It is suggested to use a CSV file editor that allows viewing CSV files as spreadsheets. Something such as https://www.libreoffice.org/ is recommended.

1. Create a CSV (.csv) file such as the "template.csv" included in the project directory (shown below):


| var            |   | en              | fr               | de              | custom     |
|----------------|---|-----------------|------------------|-----------------|-----------------|
| GREET_TEXT     |   | Hi              | Salut            | Hallo           | Ciao            |
| EXIT_TEXT      |   | Bye             | Au revoir        | Tschuss         | Adios amigo     |
| BUTTON_PRESSED |   | button pressed! | bouton enfonce ! | Knopf gedruckt! | knoppie gedruk! |

The "var" field is mandatory, and <span style="color:red">*must*</span> be the first field in the CSV!
This example only uses ASCII characters, but UTF-8 characters can also be written in the language definition fields (i.e. Tschüss).

2. To create the default split localization headers, run:
  ```bash
$ .\build\windows\amd64\llmsger.exe -f "tests/template.csv" -o "tests"
  ```
Modify this step so that you use the built executable appropriate for your machine.
You should now see four new header files created in the tests directory:
- en_lang.h
- fr_lang.h
- de_lang.h
- custom_lang.h

These files define the localized strings and their appropriate variable macro name. This default use of llmsger is less useful since it allows the user to only use one of these files in his/her project.\
\
Another use of the default llmsger command, is to merge all the language headers into one file. This way the user can define the language set to use in the project by editing the correct macro in the created file.\
An example to do so would be:
```bash
$ .\build\windows\amd64\llmsger.exe -f "tests/template.csv" -o "tests" --mrg -n my_merged
```
This creates the file "*my_merged.h*" which combines all the language definitions, which are selected through "*#ifdef*" statements.

## Advanced usage
Most modern projects have to have the ability to use all the localized strings in the project during runtime. For this reason, the llmsger includes the command *dyngen*, which creates a '.c' and '.h' file that allows the user to use all the strings in his/her project.\
An example of calling this command would be the following:
```bash
$ .\build\windows\amd64\llmsger.exe dyngen -f tests/template.csv -o "tests" --basename stef
```
This creates "*stef.c*" and "*stef.h*" in the tests directory. The source file should look like this (*stef.c*):
```c
#include "stef.h"

#define _STEF_H_C_LLMSGER_ // Guard from extern in header file

#define EN_TEXTS_LLMSGR { "Hi" /*1*/, "Bye" /*2*/, "button pressed!" /*3*/, }

#define FR_TEXTS_LLMSGR { "Salut" /*1*/, "Au revoir" /*2*/, "bouton enfonce !" /*3*/, }

#define DE_TEXTS_LLMSGR { "Hallo" /*1*/, "Tschuss" /*2*/, "knoppie gedruk!" /*3*/, }

#define CUSTOM_TEXTS_LLMSGR { "Ciao" /*1*/, "Adios amigo" /*2*/, "knoppie gedruk!" /*3*/, }



char* gp_lang_texts[MAX_STEF_LANG_OPT][MAX_STEF_VAR_OPT] =
{
	[EN] = EN_TEXTS_LLMSGR,
	[FR] = FR_TEXTS_LLMSGR,
	[DE] = DE_TEXTS_LLMSGR,
	[CUSTOM] = CUSTOM_TEXTS_LLMSGR,

};
```
These files define the variable named gp_lang_texts (this is the default name which can be modified with other flags):
```c
char* gp_lang_texts[MAX_STEF_LANG_OPT][MAX_STEF_VAR_OPT];
```
Through this variable, the user can access all the strings in a similar way as:
```c
#include <stdio.h>
#include "your_PATH\llmsger.h"

int main() {
  
		printf("EN:\t%s\n", gp_lang_texts[EN][GREET_TEXT]);
		printf("FR:\t%s\n", gp_lang_texts[FR][GREET_TEXT]);
		printf("GERMAN:\t%s\n", gp_lang_texts[DE][GREET_TEXT]);

/* Prints out:
        EN:     Hi
        FR:     Salut
        GERMAN: Hallo
*/
	return 0;
}
```
You can view the corresponding enums for your localizations in the header file.\
\
Please feel free to make as many localization strings as you would like!
\
Other helpful information can be seen when running the *help* command.
## Using the UTF-8 to ASCII code replacer
Todo readme text

## Contributing

Open source software is important for the community.

Feel free to submit a pull request for bugs or additions, and make sure to update tests as appropriate or include the files recquired for bug recreation. If you find a mistake in the docs, send a PR! Even the smallest changes help.

For major changes, open an issue first to discuss what you'd like to change.

<!-- [/lock:contributing] --------------------------------------🚫 -->

## ⭐ Found It Helpful? Star It!

If you found this project helpful, let the community know by giving it a star: ⭐

## License

See [LICENSE.md](https://github.com/Stefasaurus/llmsger/blob/main/LICENSE).