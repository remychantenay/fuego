package document

// WriteOption allows to define the expected behaviour when writing to Firetstore
type WriteOption int

// Override - the existing data (if any) will be overridden
var Override WriteOption = 1

// Merge - the new data will be MERGED to the existing (if any)
//
// Note: if used while updating a field that contains a Map or an Array
// we will remove potential duplicates
var Merge WriteOption = 2

// Append - the new data will be ADDED to the existing (if any)
//
// Note: if used while updating a field that contains a Map or an Array
// we will append the data
var Append WriteOption = 3
