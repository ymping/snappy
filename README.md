# snappy client

A compression/decompression client for Snappy

## Usage
> command arg `-o, --output : Output directory` default to current working directory

Compress files
```shell
$ snappy -c -o /path/to/output/dir file1.txt file2.txt file3.txt
```

Decompress files
```shell
$ snappy -x -o /path/to/output/dir file1.snappy file2.snappy file3.snappy
```
