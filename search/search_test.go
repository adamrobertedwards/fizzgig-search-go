package search

import (
	"testing"
)

type testCaseCommon struct {
	word           string
	comparisonWord string
}

type testCaseLevenschtein struct {
	testCaseCommon
	expected int
}

type testCaseSimilarity struct {
	testCaseCommon
	expected float64
}

type testCaseSearch struct {
	term      string
	answers   []string
	threshold float64
	expected  SearchResults
}

var testCasesLevenschtein = []testCaseLevenschtein{
	{testCaseCommon: testCaseCommon{word: "pineapple", comparisonWord: "apple"}, expected: 4},
	{testCaseCommon: testCaseCommon{word: "kitten", comparisonWord: "knitting"}, expected: 3},
	{testCaseCommon: testCaseCommon{word: "", comparisonWord: "knitting"}, expected: 8},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, test := range testCasesLevenschtein {
		distance := calculateLevenshteinDistance(test.word, test.comparisonWord)

		if distance != test.expected {
			t.Errorf(`calculateLevenshteinDistance: the levenschtein distance is incorrect. expected: %v, got: %v`, test.expected, distance)
		}
	}
}

func TestLevenshteinDistanceSlower(t *testing.T) {
	for _, test := range testCasesLevenschtein {
		distance := calculateLevenshteinDistanceSlower(test.word, test.comparisonWord)

		if distance != test.expected {
			t.Errorf(`calculateLevenshteinDistanceSlower: the levenschtein distance is incorrect. expected: %v, got: %v`, test.expected, distance)
		}
	}
}

func TestCalculateSimilarityScore(t *testing.T) {
	testCases := []testCaseSimilarity{
		{testCaseCommon: testCaseCommon{word: "apple", comparisonWord: "apple"}, expected: 1.},
		{testCaseCommon: testCaseCommon{word: "pineapple", comparisonWord: "apple"}, expected: .56},
		{testCaseCommon: testCaseCommon{word: "kitten", comparisonWord: "knitting"}, expected: .63},
		{testCaseCommon: testCaseCommon{word: "", comparisonWord: "knitting"}, expected: 0.},
	}

	for _, test := range testCases {
		score := calculateSimilarityScore(test.word, test.comparisonWord)

		if score != test.expected {
			t.Errorf(`calculateSimilarityScore: the similar score is incorrect. expected: %v, got: %v`, test.expected, score)
		}
	}
}

func TestSearch(t *testing.T) {
	answers := []string{"apple", "pineapple", "orange"}
	testCases := []testCaseSearch{
		{
			term:      "apple",
			answers:   answers,
			threshold: .5,
			expected: SearchResults{
				matches: []Match{
					{term: "apple", similarity: 1.},
					{term: "pineapple", similarity: .56},
				},
				total: 2,
			},
		},
		{
			term:      "apple",
			answers:   answers,
			threshold: 0.,
			expected: SearchResults{
				matches: []Match{
					{term: "apple", similarity: 1.},
					{term: "pineapple", similarity: .56},
					{term: "orange", similarity: .17},
				},
				total: 3,
			},
		},
		{
			term:      "kiwi",
			answers:   answers,
			threshold: .5,
			expected: SearchResults{
				matches: []Match{},
				total:   0,
			},
		},
	}

	for _, test := range testCases {
		results := Search(test.term, test.answers, test.threshold)

		if results.total != test.expected.total {
			t.Errorf(`Search: the search total is incorrect. expected: %v, got: %v`, test.expected.total, results.total)
		}

		for i, match := range results.matches {
			expectedMatch := test.expected.matches[i]
			if match.similarity != expectedMatch.similarity {
				t.Errorf(
					`Search: the similarity score for %v was incorrect. expected: %v, got: %v`,
					match.term, expectedMatch.similarity, match.similarity,
				)
			}
		}
	}
}
