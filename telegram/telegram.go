package telegram

import (
	"bytes"
	"fmt"
	gh "github.com/google/go-github/github"
	"github.com/proshik/githubstatbot/github"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"sort"
	"sync"
)

const (
	starCountTextFormat = "В репозиториях пользователя *%s* нашлось звезд в количество *%d* шт."
	forkCountTextFormat = "Репозитории пользователя *%s* форкнули целых *%d* раз"
)

var (
	//Chains
	//-commands
	startC    = make(chan tgbotapi.Update)
	authC     = make(chan tgbotapi.Update)
	languageC = make(chan tgbotapi.Update)
	starC     = make(chan tgbotapi.Update)
	forkC     = make(chan tgbotapi.Update)
	cancelC   = make(chan tgbotapi.Update)
	//-send message
	messages = make(chan tgbotapi.Chattable)
	//Randomize
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type language struct {
	Title      string
	Percentage float32
}

type userRepos struct {
	username string
	repos    []*github.Repo
}

type messageError struct {
	message string
}

func (e *messageError) Error() string { return e.message }

func (b *Bot) ReadUpdates() {
	//create timeout value
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	//read updates from telegram server
	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	//handle commands from channels
	go func() {
		for {
			select {
			case u := <-startC:
				messages <- startCommand(&u)
			case u := <-authC:
				messages <- authCommand(&u, b)
			case u := <-languageC:
				messages <- languageCommand(&u, b)
			case u := <-starC:
				messages <- calcCountCommand(&u, b, stargazersCountValue, starCountTextFormat)
			case u := <-forkC:
				messages <- calcCountCommand(&u, b, forkCountValue, forkCountTextFormat)
			case u := <-cancelC:
				messages <- cancelCommand(&u, b)
			}
		}
	}()
	//Отправка сообщений пользователям.
	//Отдельно от предыдущего блока т.к. нельзя в select нельзя обрабатывать каналы команд из которох читается(*С)
	//и куда записыватеся(messages)
	go func() {
		for res := range messages {
			b.bot.Send(res)
		}
	}()
	//read updates and send to channels
	for update := range updates {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				startC <- update
			case "auth":
				authC <- update
			case "language":
				languageC <- update
			case "star":
				starC <- update
			case "fork":
				forkC <- update
			case "cancel":
				cancelC <- update
			default:
				//show access commands
				startC <- update
			}
		} else {
			startC <- update
		}
	}
}

func (b *Bot) InformAuth(chatId int64, result bool) {
	if result {
		messages <- tgbotapi.NewMessage(chatId, "GitHub аккаунт был успешно подключен!")
	} else {
		messages <- tgbotapi.NewMessage(chatId, "Произошла ошибка при подключении GitHub аккаунта!")
	}
}

func startCommand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	buf := bytes.NewBufferString("Телеграм бот для отображения статистики GitHub аккаунта\n")
	//descriptions of commands
	buf.WriteString("\n")
	buf.WriteString("Вы можете управлять мной, отправляя следующие команды:\n\n")
	buf.WriteString("[/auth]() - авторизация через OAuth\n")
	buf.WriteString("[/language]() - статистика используемых языков в репозиториях\n")
	buf.WriteString("[/language]() <username> - статистика используемых языков в репозиториях указанного пользователя\n")
	buf.WriteString("[/star]() - статистика пожалованных звездочек в репозиториях\n")
	buf.WriteString("[/star]() <username> - статистика пожалованных звездочек в репозиториях указанного пользователя\n")
	buf.WriteString("[/fork]() - статистика форков пользовательских репозиториев\n")
	buf.WriteString("[/fork]() <username> - статистика форков репозиториев указанного пользователя\n")
	buf.WriteString("[/cancel]() - отмена авторизации\n")
	//create message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, buf.String())
	msg.ParseMode = "markdown"
	return msg
}

func authCommand(update *tgbotapi.Update, bot *Bot) tgbotapi.Chattable {
	//check, maybe user already authorize
	token, err := bot.tokenStore.Get(update.Message.Chat.ID)
	if err != nil {
		return errorMessage(update)
	}
	if token != "" {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже авторизованы!")
	}
	//generate state for url string for auth in github
	state := randStringRunes(20)
	//save to store [state]chatId
	bot.stateStore.Add(state, update.Message.Chat.ID)
	//build url
	authUrl := bot.oAuth.BuildAuthUrl(state)
	//build text for message
	buf := bytes.NewBufferString("Для авторизации перейдите по следующей ссылке:\n")
	buf.WriteString("\n")
	buf.WriteString(authUrl + "\n")
	//build message with url for user
	text := buf.String()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	return msg
}

func languageCommand(update *tgbotapi.Update, bot *Bot) tgbotapi.Chattable {
	//found token by chatId in store
	token, err := bot.tokenStore.Get(update.Message.Chat.ID)
	if err != nil || token == "" {
		log.Printf("Token=%s, Error: %v\n", token, err)
		return errorMessage(update)
	}
	//client to github
	client := github.NewClient(token)
	//get userRepos
	userRepos, err := repos(client, update.Message.CommandArguments())
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	}
	//concurrent receipt language info in repositories of an user
	wg := sync.WaitGroup{}
	languageChan := make(chan map[string]int)
	for _, repo := range userRepos.repos {
		wg.Add(1)
		go func(wg *sync.WaitGroup, r *github.Repo) {
			defer wg.Done()
			//receipt language info
			lang, err := client.Language(userRepos.username, *r.Name)
			if err != nil {
				log.Printf("Error on request language for user=%s, repo=%s", userRepos.username, *r.Name)
			}
			languageChan <- lang
		}(&wg, repo)

	}
	//wait before not will be receipt language info
	go func() {
		wg.Wait()
		close(languageChan)
	}()
	//calculate sum of a bytes in each repository by language name
	statistics := make(map[string]int)
	for stat := range languageChan {
		for k, v := range stat {
			statistics[k] = statistics[k] + v
		}
	}
	//create text messages for user
	var msg tgbotapi.MessageConfig
	if len(statistics) != 0 {
		percentages := calcLanguagePercentages(statistics)
		text := createLangStatText(userRepos.username, percentages)
		//create messages
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("У пользователя: *%s* нет репозиториев\n", userRepos.username))
	}
	msg.ParseMode = "markdown"
	return msg
}

func calcCountCommand(u *tgbotapi.Update, b *Bot, count func(r *github.Repo) *int, textFormat string) tgbotapi.Chattable {
	//found token by chatId in store
	token, err := b.tokenStore.Get(u.Message.Chat.ID)
	if err != nil || token == "" {
		return errorMessage(u)
	}
	//client to github
	client := github.NewClient(token)
	//request username and his repos
	userRepos, err := repos(client, u.Message.CommandArguments())
	if err != nil {
		return tgbotapi.NewMessage(u.Message.Chat.ID, err.Error())
	}
	//struct for calc count totalCount
	type countLock struct {
		sync.Mutex
		count int
	}
	//variable for calc count values
	var totalCount countLock
	//concurrent receipt calcCount values in repositories of an user
	wg := sync.WaitGroup{}
	for _, repo := range userRepos.repos {
		wg.Add(1)
		go func(wg *sync.WaitGroup, r *github.Repo) {
			defer wg.Done()
			//receipt language info
			r, err := client.Repo(userRepos.username, *r.Name)
			if err != nil {
				log.Printf("Error on request count for user=%s, repo=%s", userRepos.username, *r.Name)
			}
			totalCount.Lock()
			totalCount.count += *count(r)
			totalCount.Unlock()
		}(&wg, repo)

	}
	//wait before not calculate all count values by user repos
	wg.Wait()
	//create text messages for user
	text := fmt.Sprintf(textFormat, userRepos.username, totalCount.count)
	message := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	message.ParseMode = "markdown"
	return message
}

func stargazersCountValue(r *github.Repo) *int {
	return r.StargazersCount
}

func forkCountValue(r *github.Repo) *int {
	return r.ForksCount
}

func cancelCommand(update *tgbotapi.Update, bot *Bot) tgbotapi.Chattable {
	//check on exists token in store
	token, err := bot.tokenStore.Get(update.Message.Chat.ID)
	if err != nil {
		return errorMessage(update)
	}
	if token == "" {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Вы не авторизованы!")
	}
	//delete token by chatId. Exactly remove user from store
	bot.tokenStore.Delete(update.Message.Chat.ID)

	mess := tgbotapi.NewMessage(update.Message.Chat.ID, "GitHub аккаунт отключен!")

	log.Printf("Was cancel authentication user with id=%d", update.Message.Chat.ID)

	return mess
}

func repos(client *github.Client, username string) (*userRepos, error) {
	var userRepos *userRepos
	var err error
	if username == "" {
		userRepos, err = reposAuthUser(client)
	} else {
		userRepos, err = reposSpecificUser(username, client)
	}
	if err != nil {
		return nil, err
	}
	return userRepos, nil
}

func reposAuthUser(client *github.Client) (*userRepos, error) {
	//receipt info for user
	username, err := client.User()
	if err != nil {
		return nil, &messageError{"Ошибка получения данных. Выполните повторную авторизацию"}
	}
	//receipt user repositories
	repos, err := client.Repos(username)
	if err != nil {
		return nil, &messageError{"Не найдены репозитории пользователя"}
	}

	return &userRepos{username, repos}, nil
}

func reposSpecificUser(username string, client *github.Client) (*userRepos, error) {
	//receipt info for user
	username, err := client.SpecificUser(username)
	if err != nil {
		if ghErr, ok := err.(*gh.ErrorResponse); ok {
			if ghErr.Response.StatusCode == 404 {
				return nil, &messageError{"Пользовать с указанным ником не найден"}
			}
		} else {
			return nil, &messageError{"Ошибка получения данных. Выполните повторную авторизацию"}
		}
	}
	//receipt user repositories
	repos, err := client.Repos(username)
	if err != nil {
		return nil, &messageError{"Не найдены репозитории пользователя"}
	}

	return &userRepos{username, repos}, nil
}

func calcLanguagePercentages(languages map[string]int) []*language {
	result := make([]*language, 0)
	//calculate total sum byte by all languages
	var totalSum float32
	for _, v := range languages {
		totalSum += float32(v)
	}

	var totalByteOtherLanguages int
	for key, value := range languages {
		repoPercent := float32(value) * (float32(100) / totalSum)
		roundRepoPercent := round(repoPercent, 0.1)
		if roundRepoPercent >= 0.1 {
			result = append(result, &language{key, roundRepoPercent})
		} else {
			totalByteOtherLanguages += value
		}
	}
	//sort found languages by percentage
	sort.Slice(result, func(i, j int) bool { return result[i].Percentage > result[j].Percentage })
	//calculate percentage for language with less then 0.1% from total count
	if totalByteOtherLanguages != 0 {
		percent := round(float32(totalByteOtherLanguages)*(float32(100)/totalSum), 0.1)
		if percent != 0.0 {
			result = append(result, &language{"Other languages", percent})
		}
	}
	return result
}

func errorMessage(update *tgbotapi.Update) tgbotapi.Chattable {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо выполнить авторизацию. Команда /auth")
}

func createLangStatText(username string, statistics []*language) string {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("Username: *%s*\n", username))
	buf.WriteString("\n")
	for _, l := range statistics {
		buf.WriteString(fmt.Sprintf("*%s* %.1f%%\n", l.Title, l.Percentage))
	}
	return buf.String()
}

func round(x, unit float32) float32 {
	if x > 0 {
		return float32(int32(x/unit+0.5)) * unit
	}
	return float32(int32(x/unit-0.5)) * unit
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
