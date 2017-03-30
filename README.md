# fastqsplit

Splits fastq files by lane.

Takes gzipped fastq files of the form "prefix_suffix" and writes gzipped files of the form "prefix_lane_suffix".

## Usage

fastqsplit fastq_file(s)

## Notes

* fastq files must be gzipped
* Does approximately 10M lines every 3 minutes (depending on CPU speed)

## Building
go build fastqsplit.go

