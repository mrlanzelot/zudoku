package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Board type is an array of arrays with digits between 0-9, where 0 represents "no digit"
type Board [9][9]int

// ToString function for a variable with type Board
func (board Board) String() string {
	var boardString string
	for _, row := range board {
		boardString += fmt.Sprintf("%v\r\n", row)
	}
	return boardString
}

// CandidateBoard type is an array of arrays with a map that state the remaining possible digits
type CandidateBoard [9][9]map[int]bool

// createBoardFromFile takes a filename and returns an two dimensional array of digits
func createBoardFromFile(filename string) Board {
	var board Board

	for rowIndex, row := range readLines(filename) {
		textRow := strings.Split(strings.TrimSpace(row), " ")
		for colIndex, digitString := range textRow {
			digit, _ := strconv.Atoi(digitString)
			board[rowIndex][colIndex] = digit // board[y][x]
		}
	}

	return board
}

func createCandidateBoard(filename string) CandidateBoard {
	var candidate CandidateBoard
	for rowIndex, row := range readLines(filename) {
		textRow := strings.Split(strings.TrimSpace(row), " ")
		for colIndex, digitString := range textRow {
			digit, _ := strconv.Atoi(digitString)
			if digit == 0 {
				candidate[rowIndex][colIndex] = map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true}
			} else {
				candidate[rowIndex][colIndex] = map[int]bool{digit: true}
			}
		}
	}
	return candidate
}

// ToString function for a variable with type candidateBoard
func (candidateBoard CandidateBoard) String() string {
	var boardString string
	for _, row := range candidateBoard {
		for _, col := range row {
			if len(col) > 1 {
				boardString += fmt.Sprintf("[ ")
				for key := range col {
					boardString += fmt.Sprintf("%#v ", key)
				}
				boardString += fmt.Sprintf("] ")
			} else {
				for key := range col {
					boardString += fmt.Sprintf("%#v ", key)
				}
			}
		}
		boardString += fmt.Sprintf("\r\n")
	}
	return boardString
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}
	return strings.Split(string(content), "\r\n")
}

// returns 1 if the solution is complete
// returns 0 if the solution is correct, but not complete
// returns -1 if the solution is incorrect
func validate(solution Board) int {
	returnValue := validateRows(solution)

	if returnValue > -1 {
		returnValue = validateRows(transpose(solution))
	}

	if returnValue > -1 {
		returnValue = validateRows(transposeSquares(solution))
	}

	return returnValue
}

func validateRows(solution Board) int {
	returnValue := 0 // The solution is correct until proven incorrect
	correctRows := 0
	// check rows
	for rowIndex := 0; rowIndex < 9; rowIndex++ {
		sequence := solution[rowIndex]
		switch validateStatus := validateSequence(sequence); validateStatus {
		case -1:
			fmt.Println(rowIndex, sequence, " was incorrect")
			break
		case 0:
			fmt.Println(rowIndex, sequence, " is correct but not complete")
		case 1:
			fmt.Println(rowIndex, sequence, " is correct and complete")
			correctRows++
		}
	}
	fmt.Println(correctRows, " row(s) are complete")
	if correctRows == 9 {
		returnValue = 1
	}
	return returnValue
}

func validateSequence(sequence [9]int) int {
	returnValue := 0 // The sequence is correct until proven incorrect
	seq := map[int]bool{}
	for _, digit := range sequence {
		if seq[digit] == true {
			fmt.Println(digit, " appeard twice")
			returnValue = -1
			break
		}
		if digit > 0 {
			fmt.Println(digit, " was added to seq")
			seq[digit] = true
		}
	}
	if len(seq) == 9 {
		returnValue = 1
	}
	fmt.Println(seq)

	return returnValue
}

func transpose(solution Board) Board {
	var transposedBoard Board

	// populate new board where first row is the first column from the solution, ...
	for colIndex := 0; colIndex < 9; colIndex++ {
		for rowIndex := 0; rowIndex < 9; rowIndex++ {
			transposedBoard[colIndex][rowIndex] = solution[rowIndex][colIndex]
		}
	}

	return transposedBoard
}

func transposeSquares(solution Board) Board {
	var transposedBoard Board

	// populate new board where first row is the first square from the solution, ...
	for squareIndex := 0; squareIndex < 9; squareIndex++ {
		for seqIndex := 0; seqIndex < 9; seqIndex++ {
			transposedBoard[squareIndex][seqIndex] = solution[(seqIndex/3)+3*(squareIndex/3)][(seqIndex%3)+3*(squareIndex%3)]
		}
	}

	return transposedBoard
}

func solve(problem CandidateBoard) CandidateBoard {
	// Eliminate rows
	solution := eliminate(problem)

	// Eliminate cols
	solution = transposeCandidate(solution)
	solution = eliminate(solution)
	solution = transposeCandidate(solution)

	// Eliminate squars
	solution = transposeCandidateSquares(solution)
	solution = eliminate(solution)
	solution = transposeCandidateSquares(solution)

	return solution
}

func eliminate(problem CandidateBoard) CandidateBoard {
	var solution CandidateBoard
	for rowIndex := 0; rowIndex < 9; rowIndex++ {
		solution[rowIndex] = eliminateRow(problem[rowIndex])
	}
	return solution
}

// remove impossible candidates
func eliminateRow(problem [9]map[int]bool) [9]map[int]bool {
	var solution [9]map[int]bool
	// Get a list of candidates where the number of candidates in a cell is only one
	var correctDigits []int
	var indexToEliminate []int
	for index, candidates := range problem {
		if len(candidates) == 1 {
			for digit := range candidates {
				correctDigits = append(correctDigits, digit)
			}
		} else {
			indexToEliminate = append(indexToEliminate, index)
		}
		solution[index] = problem[index]
	}

	// Remove all single candidate digits from the cells with more than one candidate
	for _, index := range indexToEliminate {
		for _, digit := range correctDigits {
			delete(solution[index], digit)
		}
	}

	return solution
}

func transposeCandidate(candidate CandidateBoard) CandidateBoard {
	var transposedBoard CandidateBoard

	// populate new board where first row is the first column from the solution, ...
	for colIndex := 0; colIndex < 9; colIndex++ {
		for rowIndex := 0; rowIndex < 9; rowIndex++ {
			transposedBoard[colIndex][rowIndex] = candidate[rowIndex][colIndex]
		}
	}

	return transposedBoard
}

func transposeCandidateSquares(candidate CandidateBoard) CandidateBoard {
	var transposedBoard CandidateBoard

	// populate new board where first row is the first square from the solution, ...
	for squareIndex := 0; squareIndex < 9; squareIndex++ {
		for seqIndex := 0; seqIndex < 9; seqIndex++ {
			transposedBoard[squareIndex][seqIndex] = candidate[(seqIndex/3)+3*(squareIndex/3)][(seqIndex%3)+3*(squareIndex%3)]
		}
	}

	return transposedBoard
}

func validateCandidate(solution CandidateBoard) int {
	returnValue := validateCandidateRows(solution)

	/*	if returnValue > -1 {
			returnValue = validateRows(transpose(solution))
		}

		if returnValue > -1 {
			returnValue = validateRows(transposeSquares(solution))
		}
	*/
	return returnValue
}

func validateCandidateRows(solution CandidateBoard) int {
	returnValue := 0 // The solution is correct until proven incorrect
	correctRows := 0
	// check rows
	for rowIndex := 0; rowIndex < 9; rowIndex++ {
		sequence := solution[rowIndex]
		switch validateStatus := validateCandidateSequence(sequence); validateStatus {
		case -1:
			fmt.Println(rowIndex, sequence, " was incorrect")
			break
		case 0:
			fmt.Println(rowIndex, sequence, " is correct but not complete")
		case 1:
			fmt.Println(rowIndex, sequence, " is correct and complete")
			correctRows++
		}
	}
	fmt.Println(correctRows, " row(s) are complete")
	if correctRows == 9 {
		returnValue = 1
	}
	return returnValue
}

func validateCandidateSequence(sequence [9]map[int]bool) int {
	returnValue := 0 // The sequence is correct until proven incorrect
	/*	seq := map[int]bool{}
		for _, digit := range sequence {
			if seq[digit] == true {
				fmt.Println(digit, " appeard twice")
				returnValue = -1
				break
			}
			if digit > 0 {
				fmt.Println(digit, " was added to seq")
				seq[digit] = true
			}
		}
		if len(seq) == 9 {
			returnValue = 1
		}
		fmt.Println(seq)
	*/
	return returnValue
}

func main() {
	if len(os.Args) == 2 {
		problem := createCandidateBoard(os.Args[1])
		fmt.Println(problem)
		for i := 0; i < 20; i++ {
			problem = solve(problem)
			fmt.Println(i)
			fmt.Println(problem)
		}
		/*solutionStatus := validateCandidate(solution)

		if solutionStatus == -1 {
			fmt.Println("An incorrect solution is found:")
			fmt.Println(solution)
			return
		}
		if solutionStatus == 0 {
			fmt.Println("A correct solution is found, but it's not complete:")
			fmt.Println(solution)
			return
		}
		if solutionStatus == 1 {
			fmt.Println("A solution is found:")
			fmt.Println(solution)
		}*/
	}
}
