package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/ahmdrz/goinsta.v2"
)

//Command BOT

type CMD struct {
	Name  string `json:"name"`
	Jawab string `json:"jawab"`
}

func Command() []CMD {
	f_cmd, err := os.Open("cmd.json")
	if err != nil {
		log.Println(err.Error())
	}
	b_cmd, err := ioutil.ReadAll(f_cmd)
	if err != nil {
		log.Println(err.Error())
	}
	cmds := []CMD{}

	if err := json.Unmarshal(b_cmd, &cmds); err != nil {
		log.Println(err.Error())
	}

	return cmds

}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Baca Config Auth Instagram

func Config() *Auth {
	cfg, err := os.Open("config.json")
	if err != nil {
		return nil
	}
	read, err := ioutil.ReadAll(cfg)
	if err != nil {
		return nil
	}
	auth := &Auth{}

	if err := json.Unmarshal(read, &auth); err != nil {
		return nil
	}
	return auth
}

//Login Ke Instagram
func Login(user, pass string) (*goinsta.Instagram, bool) {
	ig := goinsta.New(user, pass)
	if err := ig.Login(); err != nil {
		return nil, false
	}
	return ig, true
}

func init() {
	log.SetFlags(log.Ltime | log.Ldate)
	log.SetPrefix("Instabot : ")
}

func main() {
	cfg := Config()
	ig, ok := Login(cfg.Username, cfg.Password)
	switch ok {
	case true:
		defer ig.Logout()
		fmt.Println("Hello", ig.Account.Username)
		fmt.Println("\n================[Olshop Asistan]====================\n")
		fmt.Println("")
		pr := func(t int) {
			if err := ig.Inbox.Sync(); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if ig.Inbox.Next() {
				for i := 0; i < ig.Inbox.UnseenCount; i++ {
					ibx := ig.Inbox.Conversations[i]

					for _, usr := range ibx.Users {
						fmt.Println("From :", usr.Username)
					}
					for _, itm := range ibx.Items {
						fmt.Println("Type :", itm.Type)
						fmt.Println("Text :", itm.Text)
						fmt.Println("")

						switch itm.Type {
						case "text":
							switch FindCommand(strings.ToLower(itm.Text)) {
							case true:
								pesan := ResponeMessage(strings.ToLower(itm.Text))
								if err := ibx.Send(pesan); err != nil {
									log.Println(err.Error())
								}
							default:
								var PesanSalah = `🙏🙏🙏 ` + "\n" + `Mohon Maaf saat ini bot tidak mengerti dengan pesan:` + "\n\n" + `"` + itm.Text + `"` + "\n\n" + `Balas Chat ini dengan format /help untuk mendapatkan bantuan` + "\n\n" + ig.Account.Username + "\n" + "Terima Kasih."
								if err := ibx.Send(PesanSalah); err != nil {
									log.Println(err.Error())
								}
							}

						case "media":
							if err := ibx.Send("Tolong kirim gambar/video di wa saja 🙏!!!\n\nBalas Chat ini dengan format /wa untuk mendapatkan nomer Whatsapp kami"); err != nil {
								log.Println(err.Error())
							}
						case "action_log":
							if err := ibx.Send("😄"); err != nil {
								log.Println(err.Error())
							}

						case "link":
							if err := ibx.Send("Terimakasih atas link yang anda bagikan 🙏 , kami akan segera mengunjungi link tersebut😄!!!"); err != nil {
								log.Println(err.Error())
							}
						case "like":
							if err := ibx.Like(); err != nil {
								log.Println(err.Error())
							}
						default:
							if err := ibx.Like(); err != nil {
								log.Println(err.Error())
							}
						}
					}
				}
			}
		}

		cout := 1
		pr(cout)
		cout++
		Tiker := time.NewTicker(time.Second).C
		for {
			select {
			case <-Tiker:
				pr(cout)
				cout++
			default:
				if !ok {
					if err := ig.Login(); err != nil {
						log.Println(err.Error())
					}
				}
			}

		}
	default:
		log.Println("Cannot Login")
		os.Exit(0)

	}

}

func FindCommand(m string) bool {
	cmd := Command()
	for _, c := range cmd {
		if c.Name == m {
			return true
		}
	}
	return false
}

func ResponeMessage(m string) string {
	cmd := Command()
	cek := FindCommand(m)
	if cek {
		for _, c := range cmd {
			if c.Name == m {
				return c.Jawab
			}
		}
	}
	return ""
}
