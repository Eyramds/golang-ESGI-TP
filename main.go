package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type player struct {
	nom, prenom, pseudo string
	score               int
	hp                  int
}

type game struct {
	number int
	winner int
	rounds int
	turns  int
}

func roll() (int, int, int) {
	return rand.Intn(6), rand.Intn(20), rand.Intn(100)
}

func printRoll() {

	de1, de2, de3 := roll()

	fmt.Println("|", de1, "|", "|", de2, "|", "|", de3, "|")

}

func printTime() {

	fmt.Println("Nous somme le", time.Now().Day(), time.Now().Month(), ", il est", time.Now().Format("15h04"))

}

func printImput() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("|")

	for {
		saisie, _ := reader.ReadString('\n')

		if strings.ToLower(saisie) != "no" {
			fmt.Println("Dans ma valise, il y a…", saisie)
		} else {
			break
		}
	}

}

func (g *game) newNumber() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	g.number = r.Intn(100) + 1
	fmt.Println("number", g.number)
}

func promptPlayer(reader *bufio.Reader, playerNb int) player {
	fmt.Println("Veuillez entrer votre nom.")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSuffix(nom, "\n")
	fmt.Println("Veuillez entrer votre prénom.")
	prenom, _ := reader.ReadString('\n')
	prenom = strings.TrimSuffix(prenom, "\n")
	fmt.Println("Veuillez entrer votre pseudo.")
	pseudo, _ := reader.ReadString('\n')
	pseudo = strings.TrimSuffix(pseudo, "\n")
	fmt.Printf("Vous êtes %s %s, mais ici on vous appelle %s, le joueur %d.\n",
		prenom,
		nom,
		pseudo,
		playerNb,
	)
	return player{
		nom:    nom,
		prenom: prenom,
		pseudo: pseudo,
		score:  0,
		hp:     100,
	}
}

func describePlayer(p player, playerNb int) {
	fmt.Printf(
		"Le joueur %d s’appelle %s %s, il porte le pseudo %s et possède un score de %d.\n",
		playerNb, p.prenom, p.nom, p.pseudo, p.score,
	)
}

func promptAddMorePlayer(reader *bufio.Reader) bool {
	fmt.Println("Voulez-vous ajouter un nouveau joueur ? (oui/non)")
	yesOrNo, _ := reader.ReadString('\n')
	return yesOrNo == "oui\n"
}

func gameStart(reader *bufio.Reader, players *[]player) {
	var g game

	g.newNumber()
	i := g.winner
	for {
		for i < len(*players) {
			if i == g.winner {
				g.turns++
			}
			fmt.Println("Au tour de", (*players)[i].pseudo)
			if guessNumber(reader, &g, (*players)[i]) {
				(*players)[i].score += 10 - g.turns
				g.winner = i
				g.rounds++
				printScores(players)
				if g.rounds >= 5 {
					return
				}
				g.turns = -1
				i--
				g.newNumber()
			}
			i++
		}
		i = 0
	}
}

func guessNumber(reader *bufio.Reader, g *game, p player) bool {
	fmt.Println("Proposez un nombre entre 1 et 100!")
	playerGuessInput, _ := reader.ReadString('\n')
	playerGuessInput = strings.TrimSpace(strings.TrimSuffix(playerGuessInput, "\n"))

	numberGuess, _ := strconv.Atoi(playerGuessInput)
	switch {
	case numberGuess == g.number:
		fmt.Println("Bien deviné ! Le nombre était bien", numberGuess)
		return true
	case numberGuess < g.number:
		fmt.Println("Trop petit ! Le nombre n'est pas", numberGuess)
	case numberGuess > g.number:
		fmt.Println("Trop grand ! Le nombre n'est pas", numberGuess)
	}
	return false
}

func printScores(players *[]player) {
	fmt.Println("Voici les scores:")
	for i := 0; i < len(*players); i++ {
		fmt.Printf("%s: %d\n", (*players)[i].pseudo, (*players)[i].score)
	}
}

func printWinner(players []player) {
	var winner player
	for _, p := range players {
		if p.score > winner.score {
			winner = p
		}
	}
	fmt.Printf("Le gagnant est %s avec %d points!\n", winner.pseudo, winner.score)
}

func readHighScores() string {
	saveData, err := os.ReadFile("./save.data")
	if err == nil {
		fmt.Println("Voici les meilleurs scores:")
		fmt.Println(string(saveData))
		return string(saveData)
	}
	return ""
}

func saveScore(previousHighScore string, winner player) {
	type scoreEntry struct {
		pseudo string
		score  int
	}

	var keptScores []scoreEntry
	for _, entry := range strings.Split(previousHighScore, "\n") {
		entryParts := strings.Split(entry, ":")
		if len(entryParts) < 2 {
			break
		}
		s, _ := strconv.Atoi(entryParts[1])
		keptScores = append(keptScores, scoreEntry{pseudo: entryParts[0], score: s})
	}
	if len(keptScores) >= 3 {
		toBeReplacedIndex := 0
		for i := range keptScores {
			if keptScores[toBeReplacedIndex].score > keptScores[i].score {
				toBeReplacedIndex = i
			}
		}
		keptScores[toBeReplacedIndex].pseudo = winner.pseudo
		keptScores[toBeReplacedIndex].score = winner.score
	} else {
		keptScores = append(keptScores, scoreEntry{winner.pseudo, winner.score})
	}

	saveFile, err := os.OpenFile("./save.data", os.O_RDWR|os.O_CREATE, 0755)
	defer saveFile.Close()

	w := bufio.NewWriter(saveFile)
	if err == nil {
		for i := range keptScores {
			fmt.Fprintf(w, "%s:%d\n", keptScores[i].pseudo, keptScores[i].score)
		}
	}
	w.Flush()
}

func getWinner(players []player) player {
	winnerIndex := 0
	for i := range players {
		if players[i].score > players[winnerIndex].score {
			winnerIndex = i
		}
	}
	return players[winnerIndex]
}

func main() {
	// reader := bufio.NewReader(os.Stdin)

	// previousHighScore := readHighScores()
	// var players []player
	// for i := 1; ; i++ {
	// 	players = append(players, promptPlayer(reader, 1))
	// 	fmt.Printf("Il y a %d joueurs.\n", i)
	// 	for pi, obj := range players {
	// 		describePlayer(obj, pi+1)
	// 	}
	// 	if !promptAddMorePlayer(reader) {
	// 		break
	// 	}
	// }
	// gameStart(reader, &players)
	// saveScore(previousHighScore, getWinner(players))

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/entries", EntryHandler)
	http.HandleFunc("/add", AddHandler)
	http.ListenAndServe(":9000", nil)

}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Hello world")
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
	}
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "%s", time.Now().Format("15h04"))
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
	}
}

type Author struct {
	author string
	entry  string
}

func EntryHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		file, err := os.Open("entries.txt")
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)

		fmt.Fprintf(w, "%v", buf.String())

	case http.MethodPost:

	}
}

func AddHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:

		http.ServeFile(w, req, "add.html")

	case http.MethodPost:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}

		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}

		if a := req.FormValue("author"); a != "" && req.FormValue("entry") != "" {

			authorEntry := Author{
				author: req.FormValue("author"),
				entry:  req.FormValue("entry"),
			}

			file, err := os.OpenFile("entries.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintln(w, "Something went bad")
				log.Fatal(err)
			}

			if authorEntry.author != os.DevNull && authorEntry.entry != os.DevNull {

			}

			if _, err := file.Write([]byte(authorEntry.author + ":" + authorEntry.entry + "\n")); err != nil {
				log.Fatal(err)
			}

			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}

		http.ServeFile(w, req, "add.html")

	}
}
