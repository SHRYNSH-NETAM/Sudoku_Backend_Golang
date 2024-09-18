package controller

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type SudokuGrid struct {
	Grid         [][]int `json:"grid"`
	GridWithBlanks [][]int `json:"gridWithBlanks"`
	GridToBeFilled [][]int `json:"gridToBeFilled"`
	History [][]int `json:"history"`
	Mode string `json:"mode"`
}

func GetSudokuGrid(w http.ResponseWriter, r *http.Request) {
	size := 3
	gridSize := size * size
	sudokuGrid := make([][]int, gridSize)
	for i := range sudokuGrid {
		sudokuGrid[i] = make([]int, gridSize)
	}

	allNumbers := make([]int, gridSize)
	for i := range allNumbers {
		allNumbers[i] = i + 1
	}

	availableNumbers := make([][]int, gridSize*gridSize)
	for i := range availableNumbers {
		availableNumbers[i] = append([]int(nil), allNumbers...)
	}

	coordinatesOfPos := func(pos, gridSize int) (int, int) {
		return pos / gridSize, pos % gridSize
	}

	numberIsValid := func(pos, num int, grid [][]int, size int) bool {
		row, col := coordinatesOfPos(pos, gridSize)
		for i := 0; i < size*size; i++ {
			if (i != col && grid[row][i] == num) || (i != row && grid[i][col] == num) {
				return false
			}
		}

		squareRow, squareCol := (row/size)*size, (col/size)*size
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if tempRow, tempCol := squareRow+i, squareCol+j; tempRow != row || tempCol != col {
					if grid[tempRow][tempCol] == num {
						return false
					}
				}
			}
		}
		return true
	}

    var solveIsUnique func(grid [][]int, solutionsCount *int, size int)
	solveIsUnique = func(grid [][]int, solutionsCount *int, size int) {
		gridSize := size * size
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				if grid[i][j] == 0 {
					for num := 1; num <= 9; num++ {
						if numberIsValid(i*gridSize+j, num, grid, size) && *solutionsCount < 2 {
							grid[i][j] = num
							solveIsUnique(grid, solutionsCount, size)
							grid[i][j] = 0
						}
					}
					return
				}
			}
		}
		*solutionsCount++
	}

	createSudoku := func(grid [][]int) [][]int {
		pos := 0
		for pos < gridSize*gridSize {
			row, col := pos/gridSize, pos%gridSize

			if len(availableNumbers[pos]) == 0 {
				grid[row][col] = 0
				availableNumbers[pos] = append([]int(nil), allNumbers...)
				pos--
			} else {
				var newNumber int
				for len(availableNumbers[pos]) > 0 {
					rand.Seed(time.Now().UnixNano())
					newNumber = availableNumbers[pos][rand.Intn(len(availableNumbers[pos]))]
					if numberIsValid(pos, newNumber, grid, size) {
						break
					} else {
						availableNumbers[pos] = removeNumber(availableNumbers[pos], newNumber)
					}
				}

				if len(availableNumbers[pos]) == 0 {
					continue
				} else {
					grid[row][col] = newNumber
					availableNumbers[pos] = removeNumber(availableNumbers[pos], newNumber)
					pos++
				}
			}
		}
		return grid
	}

	sudokuGrid = createSudoku(sudokuGrid)

	calculateNumbersToDelete := func() int {
		mode := r.URL.Query().Get("mode")
		switch mode {
		case "easy":
			return rand.Intn(8) + 40
		case "medium":
			return rand.Intn(5) + 45
		case "hard":
			return rand.Intn(5) + 50
		case "extreme":
			return rand.Intn(5) + 55
		default:
			return rand.Intn(8) + 40
		}
	}

	solutionsCount := 0

	shuffleArray := func(array []int) {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(array), func(i, j int) { array[i], array[j] = array[j], array[i] })
	}

	createGridWithBlanks := func(completeSudoku [][]int) [][]int {
		numbersToDelete := calculateNumbersToDelete()

		gridWithBlanks := make([][]int, gridSize)
		for i := range gridWithBlanks {
			gridWithBlanks[i] = append([]int(nil), completeSudoku[i]...)
		}

		shuffledPositions := make([]int, gridSize*gridSize)
		for i := range shuffledPositions {
			shuffledPositions[i] = i
		}
		shuffleArray(shuffledPositions)

		k, deletedNumbers := 0, 0
		for k < gridSize*gridSize && deletedNumbers < numbersToDelete {
			row, col := coordinatesOfPos(shuffledPositions[k], gridSize)
			prevNum := gridWithBlanks[row][col]
			gridWithBlanks[row][col] = 0

			solveIsUnique(gridWithBlanks, &solutionsCount, size)
			if solutionsCount < 2 {
				deletedNumbers++
			} else {
				gridWithBlanks[row][col] = prevNum
			}
			solutionsCount = 0
			k++
		}
		return gridWithBlanks
	}

	gridWithBlanks := createGridWithBlanks(sudokuGrid)

	response := SudokuGrid{
		Grid:          sudokuGrid,
		GridWithBlanks: gridWithBlanks,
		GridToBeFilled: gridWithBlanks,
		History: [][]int{},
	}

	currentSudoku := models.Sudokugrid{
		SolvedGrid: sudokuGrid,
		UnSolvedGrid: gridWithBlanks,
		Time: time.Now(),
	}

	if !StoreSudoku(w, r, currentSudoku) {
		http.Error(w,"Error Uploading valid Sudoku in DB", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func removeNumber(slice []int, num int) []int {
	for i, v := range slice {
		if v == num {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func StoreSudoku(w http.ResponseWriter, r *http.Request, sudoku models.Sudokugrid) bool{
	authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return true
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if err != nil || !token.Valid {
			return false
		}

		email, ok := token.Claims.(jwt.MapClaims)["email"].(string)
		if !ok {
			return false
		}

		if Success := initializers.UpdateData(models.Fuser{Email: email},models.User{CurrentSudoku: sudoku}); !Success {
			return false
		}

		err = initializers.Add2Redis(email, sudoku);
		if err!=nil {
			log.Fatal(err)
		}

		return true
}

func ValidateSudoku(w http.ResponseWriter, r *http.Request){
	var recSudoku SudokuGrid

	if err := json.NewDecoder(r.Body).Decode(&recSudoku); err!=nil {
		http.Error(w,"Error while decoding received Sudoku", http.StatusInternalServerError)
		return
	}

	for  i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if recSudoku.GridToBeFilled[i][j] != recSudoku.Grid[i][j] {
				http.Error(w, "Incorrect Solution", http.StatusInternalServerError)
				return
			}
		}
	}

	Mistakes, Cheats := detectCheatnMistake(recSudoku.Grid, recSudoku.GridWithBlanks, recSudoku.History)
	response := []float64{float64(Mistakes), float64(Cheats)}

	w.Header().Set("Content-Type", "application/json")
	var jwtPayload models.Key = "jwtPayload"

	claims, ok := r.Context().Value(jwtPayload).(jwt.MapClaims)
	if !ok {
		json.NewEncoder(w).Encode(models.ResStruct{Result: response})
		http.Error(w, "Could not retrieve JWT Payload. Please Log in Again", http.StatusUnauthorized)
		return
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		json.NewEncoder(w).Encode(models.ResStruct{Result: response})
		http.Error(w, "Email not found in JWT payload. Please Log in Again", http.StatusUnauthorized)
		return
	}
	
	var CurrentSudoku models.Sudokugrid
	// start := time.Now().UnixNano() / int64(time.Millisecond)
	CurrentSudoku,err := initializers.Get2Redis(userEmail)
	if err!=nil {
		log.Println("CurrentSudoku did not found in Cache")
		result := initializers.FindData(models.Fuser{Email: userEmail})
		if result == nil {
			json.NewEncoder(w).Encode(models.ResStruct{Result: response})
			http.Error(w, "User data not found. Please Sign In", http.StatusNotFound)
			return
		}
		CurrentSudoku = result.CurrentSudoku
	}
	// end := time.Now().UnixNano() / int64(time.Millisecond)
	// diff := end - start
    // log.Printf("Duration(ms): %d", diff)

	for  i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if recSudoku.GridToBeFilled[i][j] != CurrentSudoku.SolvedGrid[i][j] {
				http.Error(w, "Incorrect Solution", http.StatusInternalServerError)
				return
			}
		}
	}

	Mistakes, Cheats = detectCheatnMistake(CurrentSudoku.SolvedGrid, CurrentSudoku.UnSolvedGrid, recSudoku.History)
	response = []float64{float64(Mistakes), float64(Cheats), time.Since(CurrentSudoku.Time).Hours()}

	valid := true
	if valid {
		if err := UpdateMyStatistics(userEmail,recSudoku.Mode); err != nil {
			http.Error(w, "Failed to Update Statistics", http.StatusInternalServerError)
			return
		}
	}
	
	if err := json.NewEncoder(w).Encode(models.ResStruct{Result: response}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func detectCheatnMistake(solsudoku [][]int, sudoku [][]int, history [][]int) (int,int){
	mistakes := 0
	cheats := 0
	possudoku,PossHashGrid := generatePossudoku(sudoku)

	for i := range history{
		value := history[i][0]
		row := history[i][1]
		col := history[i][2]

		if value==0 {
			sudoku[row][col] = value
			possudoku,PossHashGrid = generatePossudoku(sudoku)
		} else if (len(possudoku[row][col])==1 && possudoku[row][col][0]==value) || (PossHashGrid[row][value-1]==1 || PossHashGrid[9 + col][value-1]==1 || PossHashGrid[18 + ((row / 3) * 3 + (col / 3))][value-1]==1) {
			sudoku[row][col] = value
			possudoku,PossHashGrid = generatePossudoku(sudoku)
		} else if len(possudoku[row][col])==1 && possudoku[row][col][0]!=value {
			mistakes++
			// fmt.Printf("Mistake with %v at (%v,%v), Possible values:%v\n",value,row,col,possudoku[row][col])
			sudoku[row][col] = value
			possudoku,PossHashGrid = generatePossudoku(sudoku)
		} else {
			if solsudoku[row][col]==value {
				cheats++;
				// fmt.Printf("Cheat with %v at (%v,%v), Possible values:%v\n",value,row,col,possudoku[row][col])
			} else{
				flag:=0
				for _,val := range possudoku[row][col] {
					if val==value{
						flag=1
						break
					}
				}
				if flag==0 {
					mistakes++;
					// fmt.Printf("Mistake with %v at (%v,%v), Possible values:%v\n",value,row,col,possudoku[row][col])
				}
			}
			sudoku[row][col] = value
			possudoku,PossHashGrid = generatePossudoku(sudoku)
		}
	}

	return mistakes, cheats
}

func generatePossudoku(sudoku [][]int) ([][][]int,[27][9]int) {
	var PossHashGrid [27][9]int
	possudoku := make([][][]int, 9)
	for i := range sudoku {
		possudoku[i] = make([][]int, 9)
		for j := range sudoku[i] {
			possudoku[i][j] = make([]int, 0)
			if sudoku[i][j] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafe(sudoku, i, j, num) {
						PossHashGrid[i][num-1]++
						PossHashGrid[9 + j][num-1]++
						PossHashGrid[18 + ((i / 3) * 3 + (j / 3))][num-1]++
						possudoku[i][j] = append(possudoku[i][j], num)
					}
				}
			} else {
				possudoku[i][j] = append(possudoku[i][j], 0)
			}
		}
	}
	
	return possudoku,PossHashGrid
}

func isSafe(grid [][]int, row, col, num int) bool {
    for x := 0; x < 9; x++ {
        if grid[row][x] == num {
            return false
        }
    }

    for x := 0; x < 9; x++ {
        if grid[x][col] == num {
            return false
        }
    }

    startRow := row - row%3
    startCol := col - col%3
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if grid[i+startRow][j+startCol] == num {
                return false
            }
        }
    }
    return true
}