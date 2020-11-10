// Extend the github.com/magefile/mage/sh package with helper functions.
//
// Similar to sh.RunV that modifies how stdout and stderr are handled,
// this package provides sister methods RunS / OutputS (silent) and
// RunE / OutputE (print stderr only).
package shx
