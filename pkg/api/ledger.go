package api

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// InsertMarker as used in original java code
const InsertMarker = "<!-- INSERT NEW TRANSACTION HERE:"

// LedgerPoster needs instance to sync with
type LedgerPoster struct {
	waitGroup sync.WaitGroup
	path      string
}

// NewLedgerPoster is a LedgerPoster constructor
func NewLedgerPoster(path string) *LedgerPoster {
	return &LedgerPoster{
		sync.WaitGroup{},
		path,
	}
}

// NOTE: https://stackoverflow.com/questions/34395060/how-to-implement-a-timeout-when-using-sync-waitgroup-wait?lq=1

// LedgerPost inserts into file
func (l *LedgerPoster) LedgerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newLines, err := parseNewLines(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.waitGroup.Wait()
	l.waitGroup.Add(1)

	writeToFile(l, "", newLines)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(make([]byte, 0))
	if err != nil {
		log.Error("failed to write", err)
	}
}

func parseNewLines(r *http.Request) ([]string, error) {

	var lines []string

	tid := r.Form.Get("t_id")
	date := r.Form.Get("date")
	cls := r.Form.Get("class")
	author := r.Form.Get("author")
	descr := r.Form.Get("description")

	lines = append(lines,
		fmt.Sprintf("\t<transaction id=\"%s\" date=\"%s\" class=\"%s\" author=\"%s\" description=\"%s\">\n",
			tid, date, cls, author, descr))

	balance := 0.0

	for i := 1; i < 10; i++ {
		acc := r.Form.Get("a" + strconv.Itoa(i))
		if acc == "" {
			continue
		}
		t := r.Form.Get("r" + strconv.Itoa(i))
		d := r.Form.Get("d" + strconv.Itoa(i))
		amount, err := strconv.ParseFloat(d, 64)
		if err != nil {
			return nil, err
		}

		balance = balance + amount

		lines = append(lines,
			fmt.Sprintf("\t\t<set account=\"%s\" amount=\"%s\" date=\"%s\">%s</set>\n",
				acc, d, date, t))
	}

	lines = append(lines, "\t</transaction>")

	if math.Floor(balance*100) != 0 {
		return lines, errors.New("transaction is not balanced")
	}

	return lines, nil
}

func writeToFile(l *LedgerPoster, f string, newlines []string) error {

	lines, err := readAndInsert(f, newlines)
	if err == nil {
		err = writeLines(f, lines)
	}

	l.waitGroup.Done()

	return err

}

func readAndInsert(f string, newlines []string) ([]string, error) {

	var lines []string

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, InsertMarker) {
			for _, l := range newlines {
				lines = append(lines, l)
			}
		}

		lines = append(lines, line)
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

// writes lines into named file
func writeLines(f string, rows []string) error {

	file, err := os.OpenFile(f, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, row := range rows {
		_, err = w.WriteString(row)
		if err != nil {
			return err
		}
		err = w.WriteByte('\n')
		if err != nil {
			return err
		}
	}

	return w.Flush()

}
