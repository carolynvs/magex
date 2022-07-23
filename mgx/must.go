package mgx

import "github.com/magefile/mage/mg"

// Must stops execution by throwing a panic when an error occurs.
//
// This may be used to keep your magefile brief, and mimic set -euo in an
// equivalent bash script. This pattern works well in magefile targets only, not
// helper functions. For helper functions, return an error so that you can write
// tests and allow the calling function to handle the error.
func Must(err error) {
	if err != nil {
		panic(mg.Fatal(1, err))
	}
}
