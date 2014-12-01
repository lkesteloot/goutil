// Copyright 2014 Lawrence Kesteloot

package sortutil

import (
	"os"
	"sort"
)

// An os.FileInfo slice that implements sort.Interface, comparing numbers in
// filenames properly.
type NumericalFileInfoSlice []os.FileInfo

// Return the length of the slice.
func (fi NumericalFileInfoSlice) Len() int {
	return len(fi)
}

// Return whether filename at i is less than the one at j. Embedded numbers are
// handled properly, meaning that "B2" is less than "B10".
func (fi NumericalFileInfoSlice) Less(i, j int) bool {
	return CompareStringsNumerically(fi[i].Name(), fi[j].Name()) < 0
}

// Swap objects are positions i and j.
func (fi NumericalFileInfoSlice) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

// Sort FileInfo objects in place, putting numbers in their proper order.
func SortFileInfoNumerically(fi []os.FileInfo) {
	// Use the methods on NumericalStringSlice to compare.
	sort.Sort(NumericalFileInfoSlice(fi))
}
