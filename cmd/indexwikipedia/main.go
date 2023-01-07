package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync/atomic"

	"github.com/dustin/go-wikiparse"
	"github.com/getumen/sakuin/analysis"
	"github.com/getumen/sakuin/analysis/charfilter"
	"github.com/getumen/sakuin/analysis/tokenizer"
	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/writer"
	"golang.org/x/sync/errgroup"
)

const tmpDir = "/tmp/sakuin/"

func downdload(url, fileName string) error {
	filePath := path.Join(tmpDir, fileName)
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("fail to request %s: %w", url, err)
		}
		defer resp.Body.Close()

		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("fail to create file %s: %w", filePath, err)
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return fmt.Errorf("fail to copy data %s: %w", fileName, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat error %s: %w", fileName, err)
	}
	return nil
}

func indexing(
	indexFileName, dataFileName string,
) error {
	ctx := context.Background()

	indexPath := path.Join(tmpDir, "index")

	err := os.MkdirAll(indexPath, 0700)
	if err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}
	storage, err := lsmstorage.NewStorage(indexPath)
	if err != nil {
		return fmt.Errorf("new storage error: %w", err)
	}
	defer storage.Close()

	docChan := make(chan []*document.Document)

	eg := new(errgroup.Group)

	eg.Go(func() error {
		defer close(docChan)

		parser, err := wikiparse.NewIndexedParser(
			path.Join(tmpDir, indexFileName),
			path.Join(tmpDir, dataFileName),
			runtime.NumCPU(),
		)
		if err != nil {
			return fmt.Errorf("new parser error: %w", err)
		}
		analyzer := analysis.NewAnalyzer(
			[]analysis.CharFilter{
				charfilter.NewUnicodeNFKCFilter(),
			},
			tokenizer.NewBigramTokenizer(),
			[]analysis.TokenFilter{},
		)
		documents := make([]*document.Document, 0)

		count := 0

		for count < 1000000 {
			page, err := parser.Next()
			if err != nil {
				break
			}

			titleTokens := analyzer.Analyze(page.Title)
			bodyTokens := analyzer.Analyze(page.Revisions[0].Text)

			documents = append(
				documents,
				document.NewDocument(
					page.ID,
					[]*document.Field{
						document.NewField("title", titleTokens.Terms()),
						document.NewField("body", bodyTokens.Terms()),
					},
				),
			)

			if len(documents) >= 1000 {
				count += len(documents)
				docChan <- documents
				documents = make([]*document.Document, 0)
			}
		}

		docChan <- documents

		return nil
	})

	var docNum uint64

	for i := 0; i < runtime.NumCPU(); i++ {
		eg.Go(func() error {
			for documents := range docChan {
				writer := writer.NewIndexWriter(storage)

				err = writer.CreateDocuments(ctx, documents)
				if err != nil {
					return fmt.Errorf("index docs error: %w", err)
				}
				atomic.AddUint64(&docNum, uint64(len(documents)))
				log.Printf("index %d docs", atomic.LoadUint64(&docNum))
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("worker error: %w", err)
	}
	return nil
}

func main() {

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	err := os.MkdirAll(tmpDir, 0700)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	eg := new(errgroup.Group)

	eg.Go(func() error {
		return downdload(
			"https://dumps.wikimedia.org/jawiki/latest/jawiki-latest-pages-articles-multistream-index.txt.bz2",
			"jawiki-latest-pages-articles-multistream-index.txt.bz2",
		)
	})

	eg.Go(func() error {
		return downdload(
			"https://dumps.wikimedia.org/jawiki/latest/jawiki-latest-pages-articles-multistream.xml.bz2",
			"jawiki-latest-pages-articles-multistream.xml.bz2",
		)
	})

	err = eg.Wait()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	err = indexing(
		"jawiki-latest-pages-articles-multistream-index.txt.bz2",
		"jawiki-latest-pages-articles-multistream.xml.bz2",
	)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
