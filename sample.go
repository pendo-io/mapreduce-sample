package mapreduce

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"fmt"
	"github.com/pendo-io/mapreduce"
	"net/http"
	"strconv"
	"strings"
)

type sampleUniqueWordCount struct {
	mapreduce.FileLineInputReader
	mapreduce.BlobstoreWriter
	mapreduce.StringKeyHandler
	mapreduce.IntValueHandler
	mapreduce.BlobIntermediateStorage
	mapreduce.AppengineTaskQueue
}

func (uwc sampleUniqueWordCount) Map(item interface{}) ([]mapreduce.MappedData, error) {
	line := item.(string)
	words := strings.Split(line, " ")
	result := make([]mapreduce.MappedData, 0, len(words))
	for _, word := range words {
		if len(word) > 0 {
			result = append(result, mapreduce.MappedData{word, 1})
		}
	}

	return result, nil
}

func (uwc sampleUniqueWordCount) Reduce(key interface{}, values []interface{}) (result interface{}, err error) {
	return fmt.Sprintf("%s: %d", key, len(values)), nil
}

func run(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	u := sampleUniqueWordCount{}

	job := mapreduce.MapReduceJob{
		MapReducePipeline: u,
		Inputs:            mapreduce.FileLineInputReader{[]string{"testdata/pandp-1", "testdata/pandp-2", "testdata/pandp-3", "testdata/pandp-4", "testdata/pandp-5"}},
		Outputs:           mapreduce.BlobstoreWriter{2},
		UrlPrefix:         "/mr/test",
		OnCompleteUrl:     "/done",
	}

	if jobId, err := mapreduce.Run(context, job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `running job <a href="status?id=%d">%d</a>`, jobId, jobId)
	}
}

func done(w http.ResponseWriter, r *http.Request) {
}

func blob(w http.ResponseWriter, r *http.Request) {
	elements := strings.Split(r.URL.Path, "/")
	blobKey := appengine.BlobKey(elements[len(elements)-1])
	blobstore.Send(w, blobKey)
}

func status(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	if idStr := r.FormValue("id"); idStr == "" {
		fmt.Fprintf(w, "no id given\n")
		return
	} else if val, err := strconv.ParseInt(idStr, 10, 64); err != nil {
		fmt.Fprintf(w, "bad id\n")
		return
	} else {
		key := datastore.NewKey(context, mapreduce.JobEntity, "", val, nil)
		var job mapreduce.JobInfo
		if err := datastore.Get(context, key, &job); err != nil {
			fmt.Fprintf(w, "failed to load job: %s\n", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<p>Job Stage: %s\n", job.Stage)
		if job.Stage == mapreduce.StageDone {
			fmt.Fprintf(w, "\n")
			result, err := mapreduce.GetJobResults(context, key)
			if err != nil {
				fmt.Fprintf(w, "<p>Failed to load task status: %s\n", err)
			} else {
				fmt.Fprintf(w, "<ul>\n")
				for _, result := range result {
					fmt.Fprintf(w, `<li>result: <a href="blob/%s">%s</a></li>`+"\n",
						result, result)
				}
				fmt.Fprintf(w, "</ul>\n")
			}
		}
	}
}

func init() {
	pipeline := sampleUniqueWordCount{}

	http.Handle("/mr/test/", mapreduce.MapReduceHandler("/mr/test", &pipeline, appengine.NewContext))
	http.HandleFunc("/run", run)
	http.HandleFunc("/done", done)
	http.HandleFunc("/status", status)
	http.HandleFunc("/blob/", blob)
}
