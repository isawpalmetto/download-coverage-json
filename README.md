# download-coverage-json
Downloads the provider and formulary json files from the coverage index json file

Healthcare issuers provide index files with links to the provider and forumlary coverage for their plans.  This tool downloads each of those files given the index file.

## Usage
The program takes the url for the index file and the destination to download the files.  It will download the files to the destination entered in a folder marked 'providers' or 'drugs'

`download-coverage-json [index-url] [file-destination]`
