package main

import (
	"bufio"
	// "compress/gzip"
	gzip "github.com/klauspost/pgzip"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LOG_LINE = 10000000
)

func check_result(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "list of fastq.gz files is required\n")
        	flag.PrintDefaults()
	        os.Exit(1)
	}

	for _, filename := range os.Args[1:] {
		fmt.Fprintf(os.Stderr, "%s: processing: %s\n", time.Now().String(), filename)

		var targets map[string] *gzip.Writer = make(map[string] *gzip.Writer)

		file, err := os.Open(filename)
		check_result(err)
		defer file.Close()

		gr, err := gzip.NewReader(file)
		check_result(err)
		defer gr.Close()

		scanner := bufio.NewScanner(gr)
		var record [4]string
                var line int = 0
		var written int = 0

		for scanner.Scan() {
			record[line % 4] = scanner.Text()
			if line % 4 == 3 {
				fields := strings.Split(record[0], ":")
				lane := fields[3]
				_, ok := targets[lane]
				if !ok {
					filename_components := strings.Split(filename, "_")
					target_filename := fmt.Sprintf("%s_%s_%s", strings.Join(filename_components[0:len(filename_components)-1], "_"), lane, filename_components[len(filename_components)-1])
					target_file, err := os.Create(target_filename)
					check_result(err)
					defer target_file.Close()

					targets[lane] = gzip.NewWriter(target_file)
					fmt.Fprintf(os.Stderr, "%s: opened %s for writing\n", time.Now().String(), target_filename)
				}
				for i := 0; i < 4; i++ {
					_, err := targets[lane].Write([]byte(fmt.Sprintf("%s\n", record[i])))
					check_result(err)
				}
				written += 4
			}
			line++
			if line % LOG_LINE == 0 {
				fmt.Fprintf(os.Stderr, "%s: processed %d lines\n", time.Now().String(), line)
			}
		}

		// finished with this file
		for lane := range targets {
			targets[lane].Flush()
			targets[lane].Close()
		}


		fmt.Fprintf(os.Stderr, "%s: %s: finished processing %d lines. wrote %d\n", time.Now().String(), filename, line, written)


		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} // for filename

} // main
