package usecase

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/kyoh86/my-nerds/driver/source"
	"github.com/kyoh86/my-nerds/model"
)

type RenameComicOption struct {
	ReplaceMap  map[*regexp.Regexp]string
	Force       bool
	ForceNoDiff bool
	SkipInvalid bool
	Dryrun      bool
}

func RenameComic(server *source.FTPServer, pathFrom string, option RenameComicOption) error {
	var valid string
	var doer string
	base := path.Base(pathFrom)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	comic, err := model.ParseComicFuzzy(name)
	if err != nil {
		replaced := name
		for reg, rep := range option.ReplaceMap {
			newone := reg.ReplaceAllString(replaced, rep)
			if newone != replaced {
				fmt.Printf("%q is replaced with %q -> %q to %q\n", replaced, reg.String(), rep, newone)
				replaced = newone
			}
		}
		if replaced != name {
			fmt.Printf("%q is replaced to %q\n", name, replaced)
		}
		comic, err = model.ParseComicFuzzy(replaced)
	}
	if err != nil {
		fmt.Println(pathFrom)
		if option.SkipInvalid {
			fmt.Printf("%q is invalid name: %s\n", name, err)
			return nil
		}
		if option.Dryrun {
			fmt.Printf("\tinvalid: %q (%s)\n", name, err)
			return nil
		}
		if option.Force {
			return fmt.Errorf("%q is invalid name: %w", name, err)
		}
		if surveyErr := survey.AskOne(&survey.Select{
			Message: fmt.Sprintf("%q is invalid name.\nRename it?", name),
			Options: []string{"Skip", "EDIT"},
		}, &doer); surveyErr != nil {
			return surveyErr
		}
		if doer == "Skip" {
			return nil
		}
	} else {
		valid = comic.String()
		diff := cmp.Diff(name+ext, valid+ext)
		validComic, _ := model.ParseComicName(valid)

		prop, _ := yaml.Marshal(comic)
		if diff != "" {
			if option.Dryrun {
				fmt.Println(pathFrom)
				fmt.Println(diff)
				return nil
			} else if option.ForceNoDiff && comic != nil && validComic != nil && *comic == *validComic {
				fmt.Println(pathFrom)
				fmt.Println(diff)
				doer = "Yes"
			} else if option.Force {
				doer = "Yes"
			} else {
				fmt.Println(pathFrom)
				if err := survey.AskOne(&survey.Select{
					Message: fmt.Sprintf("Rename it?\n%s\n%s", diff, string(prop)),
					Options: []string{"Skip", "Yes", "EDIT"},
				}, &doer); err != nil {
					return err
				}
				if doer == "Skip" {
					return nil
				}
			}
		}
	}
	if doer == "EDIT" {
		if comic == nil {
			comic = new(model.Comic)
		}

		if err := survey.AskOne(&survey.Input{
			Message: "Author",
			Default: comic.Author,
		}, &comic.Author); err != nil {
			return err
		}

		if err := survey.AskOne(&survey.Input{
			Message: "Title",
			Default: comic.Title,
		}, &comic.Title); err != nil {
			return err
		}

		def := "0"
		if comic.Number != 0 {
			def = strconv.FormatFloat(comic.Number, 'f', -1, 64)
		}
		if err := survey.AskOne(&survey.Input{
			Message: "Number (if it has no number, put 0)",
			Default: def,
		}, &comic.Number); err != nil {
			return err
		}

		if err := survey.AskOne(&survey.Input{
			Message: "Volume",
			Default: comic.Volume,
		}, &comic.Volume); err != nil {
			return err
		}

		if err := survey.AskOne(&survey.Input{
			Message: "Append",
			Default: comic.Append,
		}, &comic.Append); err != nil {
			return err
		}

		valid = comic.String()
		diff := cmp.Diff(name+ext, valid+ext)
		var yes bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Rename it?\n%s", diff),
		}, &yes); err != nil {
			return err
		}
		if !yes {
			fmt.Println("cancelled")
			return nil
		}
	}
	return server.Rename(pathFrom, path.Join(path.Dir(pathFrom), valid+ext))
}
