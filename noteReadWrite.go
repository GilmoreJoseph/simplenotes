//Package simplenotes provides simple writing and searching of notes.
//Notes are plaintext with a date to start them (Format: Jan 2, 2006) and tags.
//Tags are simply part of the note and anything in square brackets [].
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

//WriteNote writes the given note(string) to the given file.
//It will make new file if fileName is not found.
//The date is automatically added to beginning of the note in the format: Jan 2, 2006.
//Any tags should be in square brackets [] and are simply part of the note.
//It returns error if file cannot be opened.
func WriteNote(fileName, note string) error {

	var f *os.File

	_, err := os.Stat("notes/" + fileName)
	if os.IsNotExist(err) {
		f, err = os.OpenFile("notes/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("Tags: \n")
	} else {
		f, err = os.OpenFile("notes/"+fileName, os.O_APPEND|os.O_WRONLY, 0644)
	}

	if err != nil {
		return err
	}

	defer f.Close()

	t := time.Now().Format("Jan 2, 2006")

	//probably faster to write each of these instead of concatenting strings
	f.WriteString("\n" + t + "\n" + note + "\n")

	return nil
}

//SplitByDate will split the given file into a slice of strings.
//Dates in the file should be of format: 'Jan 2, 2006'.
//The First thing in a file cannot be a date (otherwise segfault). It is assumed there is a list of tags at the top, but not neccesary.
//It returns a reference to array of strings, as well as an error if opening the file fails.
func SplitByDate(fileName string) (*[]string, error) {

	dat, err := ioutil.ReadFile("notes/" + fileName)
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`[A-Z][a-z][a-z] [0-9]+, [0-9][0-9][0-9][0-9]`)

	dates := r.FindAllString(string(dat), -1)
	text := r.Split(string(dat), -1)

	//drop first element because it is list of tags in the file(probably created by cronjob)
	text = text[1:]

	notes := []string{}

	for i, t := range text {
		notes = append(notes, dates[i]+t)
	}

	return &notes, nil
}

//FilterDates removes all notes not within given start and end time (inclusively).
//It returns an error if time cannot be parsed correctly.
func FilterDates(start time.Time, end time.Time, notes *[]string) (err error) {

	r := regexp.MustCompile(`[A-Z][a-z][a-z] [0-9]+, [0-9][0-9][0-9][0-9]`)

	length := len(*notes)
	for i := 0; i < length; i++ {

		date, err := time.Parse("Jan 2, 2006", r.FindString((*notes)[i]))
		if err != nil {
			return err
		}

		if date.After(end) || date.Before(start) {
			//this does not perserve order of notes, but is faster
			(*notes)[i] = (*notes)[len(*notes)-1]
			*notes = (*notes)[:len(*notes)-1]
			i--
			length--
		}

	}

	return nil
}

//FilterTags removes any notes that do not contain one of the given tags.
//It escapes any neccesary characters in the tags to make into regex.
//It returns an error if the tags cannot be converted to regex.
func FilterTags(tags []string, notes *[]string) (err error) {

	tagString := ""

	for i, t := range tags {
    t = regexp.QuoteMeta(t)
		tagString += "\\[" + t + "\\]"
		if i != len(tags)-1 {
			tagString += "|"
		}
	}

	tagString = "(" + tagString + ")+"

  fmt.Println(tagString)

	r, err := regexp.Compile(tagString)
	if err != nil {
		fmt.Println(err)
		return err
	}

	length := len(*notes)
	for i := 0; i < length; i++ {
		if !r.MatchString((*notes)[i]) {
			//this does not perserve order of notes, but is faster
			(*notes)[i] = (*notes)[len(*notes)-1]
			*notes = (*notes)[:len(*notes)-1]
			i--
			length--
		}
	}

	return nil
}

//FilterExactPhrase removes all notes where the given phrase does not occur.
//Any neccesary escapes are added to phrase to create regex.
//It returns an error if phrase cannot be converted to regex
func FilterExactPhrase(phrase string, notes *[]string) (err error) {

	//escpapes any special characters
	phrase = regexp.QuoteMeta(phrase)

	r, err := regexp.Compile(phrase)
	if err != nil {
		return err
	}

	length := len(*notes)
	for i := 0; i < length; i++ {
		if !r.MatchString((*notes)[i]) {
			//this does not perserve order of notes, but is faster
			(*notes)[i] = (*notes)[len(*notes)-1]
			*notes = (*notes)[:len(*notes)-1]
			i--
			length--
		}
	}

	return nil
}
