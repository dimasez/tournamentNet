package tournamentNet

import (
	"fmt"
	"math"
	"sort"
)

type net struct {
	teamsAmount int //Количество команд. Фактическое
	teamsMax    int //Мааксимальное количество команл в сетке
	matchAmount int //Количество матчей. Фактическое
	toursAmount int //Количество туров. Финал(1) - Полуфинал (2) - ....
}

func (n *net) calcMatches() {
	n.matchAmount = n.teamsAmount - 1
}

func (n *net) calcTours() int {
	n.calcMatches()
	maxMatches := 1 //Количество матч для сетки с 1 туром
	if n.teamsAmount > 2 {
		for i := 2; i < 50; i++ {
			maxMatches = maxMatches*2 + 1
			if n.matchAmount <= maxMatches {
				n.toursAmount = i
				break
			}
		}
	} else {
		n.toursAmount = 1
	}
	n.teamsMax = int(math.Pow(float64(2), float64(n.toursAmount)))
	return n.toursAmount
}

func findEnemyIndex(teamIndex, tour int) int {
	enemyIndex := int(math.Pow(float64(2), float64(tour))) + 1 - teamIndex
	return enemyIndex
}

type Match struct {
	FirstTeamIndex  int
	SecondTeamIndex int
}

func (n *net) generatePairs() [][]Match {

	var netArray [][]Match
	var tour []Match

	tour = append(tour, Match{
		FirstTeamIndex:  1,
		SecondTeamIndex: 2,
	}) //Стартовая точка - финал из 2х команд - ЧАСТНЫЙ СЛУЧАЙ OnePair

	netArray = append(netArray, tour)

	for i := 1; i < n.toursAmount; i++ {
		prevTour := netArray[i-1]
		var tour []Match
		for _, m := range prevTour {
			tour = append(tour, Match{
				FirstTeamIndex:  m.FirstTeamIndex,
				SecondTeamIndex: findEnemyIndex(m.FirstTeamIndex, i+1),
			})
			tour = append(tour, Match{
				FirstTeamIndex:  m.SecondTeamIndex,
				SecondTeamIndex: findEnemyIndex(m.SecondTeamIndex, i+1),
			})
		}
		//fmt.Printf("Тур %v %v \n ----------\n", i+1, tour)
		netArray = append(netArray, tour)
	}

	return netArray
}

func GetOlympicNet(teamsAmount int) [][]Match {
	var firstNet = net{
		teamsAmount: teamsAmount,
	}
	firstNet.calcMatches()
	firstNet.calcTours()

	//fmt.Printf("При %v командах будет %v матчей в турнире\n", firstNet.teamsAmount, firstNet.matchAmount)
	//fmt.Println("Туров будет", firstNet.toursAmount)
	//fmt.Println("Максимальное количество команд", firstNet.teamsMax)
	//fmt.Println(firstNet.generatePairs())
	//fmt.Println("Будем передавать олимпийскую сетку")
	return firstNet.generatePairs()
}

type LevelsNet struct {
	TeamsAmount  int //Количество команд. Фактическое
	LevelsAmount int //Количество уровней
	//ToursAmount  int // Количество туров
}

func (l *LevelsNet) calcMaxTeamsOnLevel() int { //Вычисляем количество команд на один уровень
	maxTeamsForLevel := 2
	for i := 2; i < l.TeamsAmount; i++ {
		if maxTeamsForLevel*l.LevelsAmount < l.TeamsAmount {
			maxTeamsForLevel = maxTeamsForLevel + 2
		}
	}

	return maxTeamsForLevel
}

func (l *LevelsNet) fillLevels() [][]Match {

	var netArray [][]Match
	//var LevelTeams []Match
	a := l.calcMaxTeamsOnLevel()

	for i := 0; i < l.LevelsAmount; i++ {
		var tour []Match
		firstTeamIndex := a*i + 1        //Определение индекса первой команды на уровне
		controlSum := firstTeamIndex + a //Определение контрольной суммы на уровне для генерации пар
		for j := 1; j <= a/2; j++ {
			tour = append(tour, Match{
				FirstTeamIndex:  controlSum - a + j - 1,
				SecondTeamIndex: controlSum - j,
			})

		}
		//fmt.Println(tour)
		//fmt.Println("Контрольная сумма", controlSum)
		netArray = append(netArray, tour)
	}
	return netArray
}

func GetLevelsNet(teamsAmount, levelsAmount int) [][]Match {
	var tournament = LevelsNet{
		TeamsAmount:  teamsAmount,
		LevelsAmount: levelsAmount,
	}
	return tournament.fillLevels()
}

//Создание новых пар команд/матчей в 1+ туре на одном уровне уровневой системы
func NewLevelPairs(teamsArray []int) []Match {
	sort.Ints(teamsArray)
	var level []Match
	for i := 1; i <= len(teamsArray)/2; i++ {
		level = append(level, Match{
			FirstTeamIndex:  teamsArray[i-1],
			SecondTeamIndex: teamsArray[len(teamsArray)-i],
		})
	}
	return level
}

//Нахождение степени 2 для турниных пар - в зависимости от количества команд
func findDegreeIndex(teamsAmount int) int {
	degreeIndex := 1
	for i := 1; i < 40; i++ {
		if teamsAmount <= int(math.Pow(float64(2), float64(i))) && teamsAmount > int(math.Pow(float64(2), float64(i-1))) {
			degreeIndex = i
		}
	}
	return degreeIndex
}

//Нахождение максимального ранга команды в зависимости от степени 2 - например 2-в-4 = 16
func CalcMaxTeamRange(degreeIndex int) int {
	maxRangeTeam := int(math.Pow(float64(2), float64(degreeIndex)))
	return maxRangeTeam
}

//Генерация нижней сетки турнира Double-Elimination - ОДНОГО ТУРА - в зависимости от того - он НЕЧЕТНЫЙ или  чОтный
func GenerateBottomDETour(isOdd bool, teamsAmount int) []Match {
	var tour []Match
	a := findDegreeIndex(teamsAmount)
	//Два алгоритма для нечетного и четного тура
	if isOdd == true {
		controlSum := 3*int(math.Pow(float64(2), float64(a-1))) + 1
		maxRangeTeam := int(math.Pow(float64(2), float64(a))) //Считаем максимальный ранг команды
		//Количество матчей - maxRang команды / 4
		for i := 0; i < maxRangeTeam/4; i++ {
			FindFirstTeamIndex := controlSum - maxRangeTeam + i
			tour = append(tour, Match{
				FirstTeamIndex:  FindFirstTeamIndex,
				SecondTeamIndex: controlSum - FindFirstTeamIndex,
			})
		}
		fmt.Println("Контрольная сумма. Тур НЕЧЕТНЫЙ", controlSum)
		//Для четного тура
	} else {
		controlSum := int((math.Pow(float64(2), float64(a))) + 1)
		maxRangeTeam := int(math.Pow(float64(2), float64(a))) - int(math.Pow(float64(2), float64(a)))/4 //Считаем максимальный ранг команды
		//Количество матчей - maxRang команды / 3
		for i := 0; i < maxRangeTeam/3; i++ {
			FindFirstTeamIndex := controlSum - maxRangeTeam + i
			tour = append(tour, Match{
				FirstTeamIndex:  FindFirstTeamIndex,
				SecondTeamIndex: controlSum - FindFirstTeamIndex,
			})
		}

	}
	return tour
}

//Double Elemination - cоздаем нижнюю сетку
func GenerateBottomDENet(teamsAmount int) [][]Match {
	//Определяем степень ДВОЙКИ в зависимости от кол-ва команд
	degreeIndex := findDegreeIndex(teamsAmount)
	//Определяем максимальный ранг команды
	maxRangeTeam := CalcMaxTeamRange(degreeIndex)
	var bottomDENet [][]Match

	//Итерируем по степени два по не достинет степени 2

	//fmt.Println("------- новая нижняя сетка -------")
	for i := degreeIndex; i >= 2; i = i - 1 {

		//Генерим тур нижней сетки
		bottomDENet = append(bottomDENet, GenerateBottomDETour(true, maxRangeTeam))
		bottomDENet = append(bottomDENet, GenerateBottomDETour(false, maxRangeTeam))
		maxRangeTeam = maxRangeTeam / 2
	}
	//fmt.Println(bottomDENet)
	//fmt.Println("------- конец сетки -------")

	return bottomDENet
}

//Структура для уровневой системы - задается количество команд + количество матчей на уровне
type LevelsNetVersion2 struct {
	TeamsAmount        int //Количество команд. Фактическое
	GamesOnLevelAmount int //Количество Игр на уровне
}

//Определяем количество уровней для уровневой системы v2
func (l *LevelsNetVersion2) calcLevels() int {
	levelsAmount := 1
	for i := 1; i <= l.TeamsAmount/2; i++ {
		if levelsAmount*l.GamesOnLevelAmount*2 < l.TeamsAmount {
			levelsAmount++
		} else {
			break
		}
	}
	return levelsAmount
}

//Создание уровневой системы - формирование пар команд - заполнение уровней
func (l *LevelsNetVersion2) fillLevels() [][]Match {
	var netArray [][]Match
	a := l.calcLevels()

	for i := 0; i < a; i++ {
		var tour []Match
		firstTeamIndex := 1 + i*l.GamesOnLevelAmount*2        //Определение индекса первой команды на уровне
		controlSum := firstTeamIndex + l.GamesOnLevelAmount*2 //Определение контрольной суммы на уровне для генерации пар
		for j := 1; j <= l.GamesOnLevelAmount; j++ {
			tour = append(tour, Match{
				FirstTeamIndex:  controlSum - l.GamesOnLevelAmount*2 + j - 1,
				SecondTeamIndex: controlSum - j,
			})
		}
		netArray = append(netArray, tour)
	}
	return netArray
}

func FillLevelsNetType(teamsAmount, gamesOnLevel int) [][]Match {
	var netArray [][]Match
	var l = LevelsNetVersion2{
		TeamsAmount:        teamsAmount,
		GamesOnLevelAmount: gamesOnLevel,
	}
	a := l.calcLevels()

	for i := 0; i < a; i++ {
		var tour []Match
		firstTeamIndex := 1 + i*l.GamesOnLevelAmount*2        //Определение индекса первой команды на уровне
		controlSum := firstTeamIndex + l.GamesOnLevelAmount*2 //Определение контрольной суммы на уровне для генерации пар
		for j := 1; j <= l.GamesOnLevelAmount; j++ {
			tour = append(tour, Match{
				FirstTeamIndex:  controlSum - l.GamesOnLevelAmount*2 + j - 1,
				SecondTeamIndex: controlSum - j,
			})
		}
		netArray = append(netArray, tour)
	}
	return netArray
}
