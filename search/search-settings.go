package search

// ExcerptPadding How many words to put to the left and right of the search term when generating an excerpt
const ExcerptPadding = 20

// ProximitySearchRadius When doing loose term search, how far to look for the next term
const ProximitySearchRadius = 5

// ExcerptOverlapThreshold When merging document search result, disregard excerpts in the requested merge that are
// closer than this threshold
const ExcerptOverlapThreshold = 20
