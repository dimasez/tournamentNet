package tournamentNet

import (
	"math"
)

type net struct {
	teamsAmount int //Количество команд. Фактическое
	teamsMax    int //Мааксимальное количество команл в сетке
	matchAmount int //Количество матчей. Фактическое
	toursAmount int //Количество туров. Финал(1) - Полуфинал (2) - ....
}

func (n *net) CalcMatches() {
	n.matchAmount = n.teamsAmount - 1
}

func (n *net) CalcTours() int {
	n.CalcMatches()
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

func (n *net) GeneratePairs() [][]Match {

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
	firstNet.CalcMatches()
	firstNet.CalcTours()

	//fmt.Printf("При %v командах будет %v матчей в турнире\n", firstNet.teamsAmount, firstNet.matchAmount)
	//fmt.Println("Туров будет", firstNet.toursAmount)
	//fmt.Println("Максимальное количество команд", firstNet.teamsMax)
	//fmt.Println(firstNet.GeneratePairs())
	//fmt.Println("Будем передавать олимпийскую сетку")
	return firstNet.GeneratePairs()
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
	//fmt.Println("Создаем жеребьевку уровневой системы")
	//fmt.Printf("В турнире участвует %v команд\n", tournament.TeamsAmount)
	//fmt.Printf("В турнире %v уровней\n", tournament.LevelsAmount)
	//fmt.Printf("Каждый уровень содержит по %v команд\n", tournament.calcMaxTeamsOnLevel())
	//fmt.Println("Будем передавать уровневую сетку!!!")
	return tournament.fillLevels()
}
