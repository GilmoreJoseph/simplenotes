package main

import (
    "log"
    "net/http"
    "net/url"
		"io/ioutil"
    "fmt"
    "encoding/json"
    "time"
)

//Message is just used to decode the json in put request
type Message struct {
    Note string
}

type Query struct {
    Start string
    End string
    Tags []string
    Phrase string
}

//handleQuery parses the query (a map) and returns notes.
func handleQuery(fileName string, q url.Values) (*[]string, error) {

    notes, err := SplitByDate(fileName)
    if(err != nil) {
        return nil, err
    }

    var query Query

    if start, ok := q["start"]; ok {
        query.Start = start[0]
    }
    if end, ok := q["end"]; ok {
        query.End = end[0]
    }
    if tags, ok := q["tag"]; ok {
        query.Tags = tags
    }
    if phrase, ok := q["phrase"]; ok {
        query.Phrase = phrase[0]
    }

    //filter by dates
    if(query.Start != "" && query.End != ""){
        s, err := time.Parse("Jan 2, 2006", query.Start)
        e, err := time.Parse("Jan 2, 2006", query.End)
        err = FilterDates(s, e, notes)
        if (err != nil){
            return nil, err
        }
    } else if (query.Start != "") {
        s, err := time.Parse("Jan 2, 2006", query.Start)
        err = FilterDates(s, time.Unix(1<<63-62135596801, 999999999), notes)
        if (err != nil){
            return nil, err
        }
    } else if (query.End != "") {
        e, err := time.Parse("Jan 2, 2006", query.End)
        err = FilterDates(time.Time{}, e, notes)
        if (err != nil){
            return nil, err
        }
    }

    //filter by tags
    if(len(query.Tags) != 0){
        err = FilterTags(query.Tags, notes)
        if(err != nil){
            return nil, err
        }
    }

    //filter by phrase
    if(query.Phrase != ""){
        err = FilterExactPhrase(query.Phrase, notes)
        if(err != nil){
            return nil, err
        }
    }

    return notes, nil

}

// notesHandler is the handler for al request to the /notes endpoint.
// Accpets only PUT or GET requests. PUT takes json with Note as the note body.
// for GET it calls handleQuery, and returns list of notes
func notesHandler(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
		noteFile := r.URL.Path[len("/notes/"):]

    switch r.Method {
        case "GET":
            notes, err := handleQuery(noteFile, r.URL.Query())
            if(err != nil){
				        w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(`{"Error": "Could not read. Could be your fault, could be my fault"}`))
            } else {
                json.NewEncoder(w).Encode(*notes)
                w.WriteHeader(http.StatusOK)
            }

        case "PUT":

            var mes Message
            err := json.NewDecoder(r.Body).Decode(&mes)
            if(err != nil || mes.Note == "") {
				        w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(`{"Error": "Could not decode json. expected Note"}`))
            } else {
                err := WriteNote(noteFile, mes.Note)
                if (err != nil){
				            w.WriteHeader(http.StatusInternalServerError)
                    w.Write([]byte(`{"Error": "Could not write to file"}`))
                }
                fmt.Println("wrote note: " + mes.Note)
                w.WriteHeader(http.StatusOK)
            }

        default:
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"Error": "method not found, use PUT or GET"}`))
    }

}

func main() {

    //home page is just simple html file (index.html)
    fs := http.FileServer(http.Dir("./static"))

    http.Handle("/", fs)
    http.HandleFunc("/notes/", notesHandler)

    log.Fatal(http.ListenAndServe(":80", nil))
}
