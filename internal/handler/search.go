package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"wordsearch/internal/response"
	"wordsearch/pkg/searcher"
)

func Search(w http.ResponseWriter, r *http.Request) {
	searcher := searcher.Searcher{
		FS: os.DirFS("../../examples/"),
	} 

	word := r.FormValue("word")
	fileNames, err := searcher.Search(word)
	if err != nil {
		writeError(w, err)
	}

	res := response.SearchResponse{Files: fileNames}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// Немного бойлерплейта из-за использования std net/http, но мне так захотелось
func writeError(w http.ResponseWriter, err error) {
	var res response.SearchResponse

	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		errs := uw.Unwrap()
		for _, v := range errs {
			errStruct := response.Error{
				Message: v.Error(),
			}
			res.Errors = append(res.Errors, errStruct)
		}
	} else {
		res.Errors = append(res.Errors, response.Error{Message: err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
