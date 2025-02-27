package search

import (
	"math"
)

type Match struct {
	term       string
	similarity float64
}

type SearchResults struct {
	term    string
	matches []Match
	total   int
}

// calculateLevenshteinDistance is a function that returns the Levenschtein Distance between two strings in an optimised way, by tracking and processing the previous and current row of a cost matrix.
func calculateLevenshteinDistance(str string, comparisonStr string) int {
	strLen, comparisonStrLen := len(str), len(comparisonStr)
	prevRow := make([]int, comparisonStrLen+1)
	currentRow := make([]int, comparisonStrLen+1)

	// Initialise the previous row
	for col := 0; col <= comparisonStrLen; col++ {
		prevRow[col] = col
	}

	// Populate the matrix to calculate the distances
	for row := 1; row <= strLen; row++ {
		// Initialise the first column of the current row, for the "deletion" case
		currentRow[0] = row

		for col := 1; col <= comparisonStrLen; col++ {
			substitutionCost := 1

			// If the characters are equal, the cost should be 0
			if str[row-1] == comparisonStr[col-1] {
				substitutionCost = 0
			}

			// Calculate the minimum of either insertion, deletion or substitution
			currentRow[col] = min(
				currentRow[col-1]+1,
				prevRow[col]+1,
				prevRow[col-1]+substitutionCost,
			)
		}

		// Swap the previous and current rows for the next iteration
		// It's important to swap both, because without doing so, prevRow and currentRow will point
		// to the same slice in memory which would be overwritten in the next iteration.
		// Re-assigning is also cheaper than creating a brand new slice on each iteration
		prevRow, currentRow = currentRow, prevRow
	}

	// The distance is stored in the last column of the prevRow, after being swapped
	return prevRow[comparisonStrLen]
}

// calculateLevenshteinDistanceSlower is a function that returns the Levenschtein Distance between two strings in an unoptimised way, by generating a full cost matrix of changes.
func calculateLevenshteinDistanceSlower(str string, comparisonStr string) int {
	strLen, comparisonStrLen := len(str), len(comparisonStr)
	matrix := make([][]int, strLen+1)

	for i := range matrix {
		matrix[i] = make([]int, comparisonStrLen+1)
	}

	// Initialise the matrix values
	for row := 0; row <= strLen; row++ {
		matrix[row][0] = row
	}

	for col := 0; col <= comparisonStrLen; col++ {
		matrix[0][col] = col
	}

	// Populate the matrix to calculate the distances
	for row := 1; row <= strLen; row++ {
		for col := 1; col <= comparisonStrLen; col++ {
			substitutionCost := 1

			// If the characters are equal, the cost should be 0
			if str[row-1] == comparisonStr[col-1] {
				substitutionCost = 0
			}

			// Calculate the minimum of either insertion, deletion or substitution
			matrix[row][col] = min(
				matrix[row][col-1]+1,
				matrix[row-1][col]+1,
				matrix[row-1][col-1]+substitutionCost,
			)
		}
	}

	// Return the bottom corner of the matrix for the distance
	return matrix[strLen][comparisonStrLen]
}

// calculateSimilarityScore is a function that calculates and returns a similarity score between 0 - 1 for two given strings.
func calculateSimilarityScore(str string, comparisonStr string) float64 {
	distance := calculateLevenshteinDistance(str, comparisonStr)
	maxLength := max(
		len(str),
		len(comparisonStr),
	)

	// The similarity score is based on the number of letters that didn't need to change
	// to make the strings equal, compared to the total string length
	score := (float64(maxLength) - float64(distance)) / float64(maxLength)

	// Round the score, to return the value to 2 decimal places
	return math.Round(score*100) / 100
}

// Search is a public function that takes a search term, a slice of possible answers and a threshold, and returns a sub set of answers which are matched to the search term with a similarity score, above the provided threshold.
func Search(term string, answers []string, threshold float64) SearchResults {
	matches := []Match{}

	for _, answer := range answers {
		similarity := calculateSimilarityScore(term, answer)

		// Exclude any terms that don't meet the minimum score threshold
		if similarity <= threshold {
			continue
		}

		matches = append(matches, Match{
			term:       answer,
			similarity: similarity,
		})
	}

	total := len(matches)

	return SearchResults{
		term,
		matches,
		total,
	}
}
