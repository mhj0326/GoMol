package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

// ReadBLOSUM62 reads the BLOSUM62 matrix from a CSV file
func ReadBLOSUM62() error {
	file, err := os.Open("blosum62.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	matrix, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for i, row := range matrix {
		if i == 0 {
			continue // Skip the header row
		}

		for j, cell := range row {
			if j == 0 || i == j {
				continue // Skip the header column and diagonal
			}

			score, err := strconv.Atoi(cell)
			if err != nil {
				return err
			}

			BLOSUM62[AminoPair{First: rune(matrix[0][j][0]), Second: rune(matrix[i][0][0])}] = score
			BLOSUM62[AminoPair{First: rune(matrix[i][0][0]), Second: rune(matrix[0][j][0])}] = score
		}
	}
	return nil
}

// score returns the BLOSUM62 score for a pair of amino acids
func score(a, b rune) int {
	return BLOSUM62[AminoPair{a, b}]
}

// max returns the maximum value from a slice of integers
func max(values ...int) (maxVal int, maxIndex int) {
	maxVal = values[0]
	maxIndex = 0
	for i, v := range values {
		if v > maxVal {
			maxVal = v
			maxIndex = i
		}
	}
	return maxVal, maxIndex
}

// needlemanWunsch performs the Needleman-Wunsch algorithm for sequence alignment
// needlemanWunsch performs the Needleman-Wunsch algorithm for sequence alignment
func NeedlemanWunsch(seq1, seq2 string) (string, string, string, float64) {
	gapPenalty := -10 // Define gap penalty
	m, n := len(seq1), len(seq2)
	dp := make([][]int, m+1) // Initialize the scoring matrix
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	// Initialize first row and column of the scoring matrix
	for i := 0; i <= m; i++ {
		dp[i][0] = i * gapPenalty
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j * gapPenalty
	}
	// Fill the scoring matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			match := dp[i-1][j-1] + score(rune(seq1[i-1]), rune(seq2[j-1]))
			delete := dp[i-1][j] + gapPenalty
			insert := dp[i][j-1] + gapPenalty
			dp[i][j], _ = max(match, delete, insert)
		}
	}
	//to find the best alignment and calculate alignment score
	align1, align2, matchLine := "", "", ""
	matchingCount := 0   // Count of matching residues
	alignmentLength := 0 // Total length of the alignment
	i, j := m, n
	for i > 0 && j > 0 {
		scoreCurrent := dp[i][j]
		scoreDiagonal := dp[i-1][j-1]
		scoreLeft := dp[i-1][j]
		if scoreCurrent == scoreDiagonal+score(rune(seq1[i-1]), rune(seq2[j-1])) {
			if seq1[i-1] == seq2[j-1] {
				matchingCount++
				matchLine = "|" + matchLine // Mark as match
			} else {
				matchLine = " " + matchLine // Mark as mismatch
			}
			align1 = string(seq1[i-1]) + align1
			align2 = string(seq2[j-1]) + align2
			i--
			j--
		} else if scoreCurrent == scoreLeft+gapPenalty {
			matchLine = " " + matchLine // Mark gap in seq2
			align1 = string(seq1[i-1]) + align1
			align2 = "-" + align2
			i--
		} else {
			matchLine = " " + matchLine // Mark gap in seq1
			align1 = "-" + align1
			align2 = string(seq2[j-1]) + align2
			j--
		}
		alignmentLength++
	}
	// Handle any remaining characters in seq1 or seq2
	for i > 0 {
		align1 = string(seq1[i-1]) + align1
		align2 = "-" + align2
		matchLine = " " + matchLine // Mark gap in seq2
		i--
		alignmentLength++
	}
	for j > 0 {
		align1 = "-" + align1
		align2 = string(seq2[j-1]) + align2
		matchLine = " " + matchLine // Mark gap in seq1
		j--
		alignmentLength++
	}
	// Calculate the percentage similarity
	percentSimilarity := 0.0
	if alignmentLength > 0 {
		percentSimilarity = float64(matchingCount) / float64(alignmentLength) * 100
	}
	return align1, align2, matchLine, percentSimilarity
}
