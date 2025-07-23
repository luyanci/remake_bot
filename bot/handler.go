package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	bot    *tele.Bot
	logger *zap.Logger

	mutex sync.Mutex

	remake *Remake
}

func NewHandler(bot *tele.Bot, logger *zap.Logger, remake *Remake) *Handler {
	return &Handler{
		bot:    bot,
		logger: logger,
		remake: remake,
	}
}

func (h *Handler) RegisterAll() {
	h.bot.Handle("/remake", h.CommandRemake)
	h.bot.Handle("/remake_data", h.CommandRemakeData)
	h.bot.Handle("/eat", h.CommandEat)
	h.bot.Handle("/jeff", h.CommandJeff)
}

func (h *Handler) getRandomCountry() Country {
	// 生成随机数
	randomNum := rand.Int63n(h.remake.TotalPopulation)

	// 根据随机数获取对应的国家
	index := 0
	for i, country := range h.remake.CountryList {
		if randomNum < country.Population {
			index = i
			break
		}
		randomNum -= country.Population
	}

	return h.remake.CountryList[index]
}

func (h *Handler) CommandRemake(c tele.Context) error {
	msg := c.Message()

	remakeData := []string{"男孩子", "女孩子", "MtF", "FtM", "MtC", "萝莉", "正太", "武装直升机", "沃尔玛购物袋", "星巴克", "无性别", "扶她", "死胎", "xyn", "Furry", "变态", "鲨鲨", "鸽子", "狗狗", "海鸥", "猫猫", "鼠鼠", "猪猪", "薯条", "GG Bond", "老色批", "柚子厨","小杂鱼~","小八嘎"}
	remakeLocate := []string{"首都", "省会", "直辖市", "市区", "县城", "自治区", "农村", "大学", "沙漠"}

	remakeResult := rand.Intn(len(remakeData))
	remakeResult_Locate := rand.Intn(len(remakeLocate))
	randomCountry := h.getRandomCountry()

	func() {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		_, hasKey := h.remake.RemakeCount[c.Sender().ID]
		if !hasKey {
			h.remake.RemakeCount[c.Sender().ID] = new(RemakeData)
		}
		oldGender := h.remake.RemakeCount[c.Sender().ID].count
		h.remake.RemakeCount[c.Sender().ID] = &RemakeData{
			country: randomCountry.CountryName,
			locate:  remakeLocate[remakeResult_Locate],
			gender:  remakeData[remakeResult],
			count:   oldGender + 1,
		}
	}()

	text := fmt.Sprintf("重生成功！您出生在 %s 的 %s ，是 %s 喵。", randomCountry.CountryName, remakeLocate[remakeResult_Locate], remakeData[remakeResult])

	_, err := c.Bot().Reply(msg, text)
	if err != nil {
		return err
	}

	if c.Chat().Type == tele.ChatPrivate {
		return nil
	}

	time.AfterFunc(5*time.Second, func() {
		if err != nil {
			return
		}
	})
	return nil
}

func (h *Handler) CommandRemakeData(c tele.Context) error {

	msg := c.Message()

	var text string
	userData, hasKey := h.remake.RemakeCount[c.Sender().ID]
	if hasKey {
		text = fmt.Sprintf("您现在是 %s %s 的 %s ，共 remake 了 %d 次", userData.country, userData.locate, userData.gender, userData.count)
	} else {
		text = "您还没有 remake 过呢，快 /remake 吧"
	}

	_, err := c.Bot().Reply(msg, text)
	if err != nil {
		return err
	}

	if c.Chat().Type == tele.ChatPrivate {
		return nil
	}

	time.AfterFunc(10*time.Second, func() {
		// err = c.Bot().Delete(reply)
		// err = c.Bot().Delete(msg)
		if err != nil {
			return
		}
	})
	return nil
}

func (h *Handler) CommandEat(c tele.Context) error {
	//if !(c.Chat().Type == tele.ChatPrivate || c.Chat().ID == -1001965344356) {
	//	fmt.Println(c.Chat().ID)
	//	return nil
	//}

	var userList []string
	chatID := c.Chat().ID
	chat, err := h.bot.ChatByID(chatID)
	if err != nil {
		return err
	}
	if chat.Type == tele.ChatPrivate {
		// 如果是私聊，只添加发送者自己
		userList = append(userList, c.Sender().FirstName)
	} else {
		if c.Message().ReplyTo != nil && c.Message().ReplyTo.Text != "" {
			// 如果是回复消息且不是topic消息(回复文本不为空)，添加回复的用户
			userList = append(userList, c.Message().ReplyTo.Sender.FirstName)
		} else {
			// 吃管理吧
			members, err := h.bot.AdminsOf(chat)
			if err != nil {
				return err
			}
			for _, member := range members {
				var name string
				if member.User.ID == c.Sender().ID || member.User.ID == c.Bot().Me.ID {
					continue
				}
				if member.User.FirstName != "" {
					name = member.User.FirstName
				} else if member.User.Username != "" {
					name = member.User.Username
				} else {
					continue
				}
				userList = append(userList, name)
			}
			if len(userList) == 0 {
				// 如果没有管理员，添加发送者自己
				userList = append(userList, c.Sender().FirstName)
			}
		}
	}

	method := []string{"炒", "蒸", "煮", "红烧", "爆炒", "烤", "炸", "煎", "炖", "焖", "炖", "卤"}

	loc := time.FixedZone("Asia/Shanghai", 8*60*60)
	now := time.Now().In(loc)
	// 获取时间段
	hour := now.Hour()
	var hourText string
	switch {
	case hour >= 6 && hour < 10:
		hourText = "早上"
	case hour >= 10 && hour < 14:
		hourText = "中午"
	case hour >= 14 && hour < 17:
		hourText = "下午"
	case hour >= 18 && hour <= 21:
		hourText = "晚上"
	default:
		hourText = "宵夜"
	}

	var name string
	if strings.Contains(c.Sender().FirstName, " | ") {
		name = strings.Split(c.Sender().FirstName, " | ")[0]
	} else {
		name = c.Sender().FirstName
	}

	result := fmt.Sprintf("%s 今天%s吃 %s %s", name, hourText, method[rand.Intn(len(method))], userList[rand.Intn(len(userList))])
	return c.Reply(result)
}

func (h *Handler) CommandJeff(c tele.Context) error {
	text, err := getRandomLine()
	username := ""
	if err != nil {
		return nil
	}
	if c.Message().ReplyTo != nil {
		if c.Message().ReplyTo.Text != "" {
			// topic message cannot get reply sender correctly
			username = c.Message().ReplyTo.Sender.FirstName
		} else {
			username = c.Sender().FirstName
		}
	} else {
		username = c.Sender().FirstName
	}

	result := replacePlaceholder(text, username)
	return c.Reply(result)
}
