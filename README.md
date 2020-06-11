# testadsprocessor
The project is used to process the test ads before uploading to S3.
It mainly does following:
1. Unwrap the adm from ad response.
2. The processed files will be put under a new folder at the same place of original folder but with postfix _processed.
3. All the characters excepts 0-9, a-z, A-Z, _ are replaced by _ in the file name.

How to run:
1. Get the code.
2. Update config.ini based on your test ads path.
3. Run cmd: go run cmd/testadsprocessor/testadsprocessor.go.
4. Check the output.

