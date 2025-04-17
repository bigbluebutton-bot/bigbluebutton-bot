package bot

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	bbb "github.com/bigbluebutton-bot/bigbluebutton-bot/bbb"
	pad "github.com/bigbluebutton-bot/bigbluebutton-bot/pad"

	convert "github.com/benpate/convert"
)

func getCookieByName(cookies []*http.Cookie, name string) string {
	result := ""
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return result
}

type Language string

const (
	af     Language = "af"
	ar     Language = "ar"
	az     Language = "az"
	bg_BG  Language = "bg-BG"
	bn     Language = "bn"
	ca     Language = "ca"
	cs_CZ  Language = "cs-CZ"
	da     Language = "da"
	de     Language = "de"
	dv     Language = "dv"
	el_GR  Language = "el-GR"
	en     Language = "en"
	eo     Language = "eo"
	es     Language = "es"
	es_419 Language = "es-419"
	es_ES  Language = "es-ES"
	es_MX  Language = "es-MX"
	et     Language = "et"
	eu     Language = "eu"
	fa_IR  Language = "fa-IR"
	fi     Language = "fi"
	fr     Language = "fr"
	gl     Language = "gl"
	he     Language = "he"
	hi_IN  Language = "hi-IN"
	hr     Language = "hr"
	hu_HU  Language = "hu-HU"
	hy     Language = "hy"
	id     Language = "id"
	it_IT  Language = "it-IT"
	ja     Language = "ja"
	ka     Language = "ka"
	km     Language = "km"
	kn     Language = "kn"
	ko_KR  Language = "ko-KR"
	lo_LA  Language = "lo-LA"
	lt_LT  Language = "lt-LT"
	lv     Language = "lv"
	ml     Language = "ml"
	mn_MN  Language = "mn-MN"
	nb_NO  Language = "nb-NO"
	nl     Language = "nl"
	oc     Language = "oc"
	pl_PL  Language = "pl-PL"
	pt     Language = "pt"
	pt_BR  Language = "pt-BR"
	ro_RO  Language = "ro-RO"
	ru     Language = "ru"
	sk_SK  Language = "sk-SK"
	sl     Language = "sl"
	sr     Language = "sr"
	sv_SE  Language = "sv-SE"
	ta     Language = "ta"
	te     Language = "te"
	th     Language = "th"
	tr_TR  Language = "tr-TR"
	uk_UA  Language = "uk-UA"
	vi_VN  Language = "vi-VN"
	zh_CN  Language = "zh-CN"
	zh_TW  Language = "zh-TW"
)

func AllLanguages() []Language {
	return []Language{
		af, ar, az, bg_BG, bn, ca, cs_CZ, da, de, dv, el_GR, en, eo, es, es_419,
		es_ES, es_MX, et, eu, fa_IR, fi, fr, gl, he, hi_IN, hr, hu_HU, hy, id,
		it_IT, ja, ka, km, kn, ko_KR, lo_LA, lt_LT, lv, ml, mn_MN, nb_NO, nl,
		oc, pl_PL, pt, pt_BR, ro_RO, ru, sk_SK, sl, sr, sv_SE, ta, te, th, tr_TR,
		uk_UA, vi_VN, zh_CN, zh_TW,
	}
}

func LanguageShortToName(short Language) string {
	switch short {
	case af:
		return "Afrikaans"
	case ar:
		return "العربية"
	case az:
		return "Azərbaycan dili"
	case bg_BG:
		return "Български"
	case bn:
		return "বাংলা"
	case ca:
		return "Català"
	case cs_CZ:
		return "Čeština"
	case da:
		return "Dansk"
	case de:
		return "Deutsch"
	case dv:
		return "ދިވެހި"
	case el_GR:
		return "Ελληνικά"
	case en:
		return "English"
	case eo:
		return "Esperanto"
	case es:
		return "Español"
	case es_419:
		return "Español (Latinoamérica)"
	case es_ES:
		return "Español (España)"
	case es_MX:
		return "Español (México)"
	case et:
		return "eesti keel"
	case eu:
		return "Euskara"
	case fa_IR:
		return "فارسی"
	case fi:
		return "Suomi"
	case fr:
		return "Français"
	case gl:
		return "Galego"
	case he:
		return "עברית‏"
	case hi_IN:
		return "हिन्दी"
	case hr:
		return "Hrvatski"
	case hu_HU:
		return "Magyar"
	case hy:
		return "Հայերեն"
	case id:
		return "Bahasa Indonesia"
	case it_IT:
		return "Italiano"
	case ja:
		return "日本語"
	case ka:
		return "ქართული"
	case km:
		return "ភាសាខ្មែរ"
	case kn:
		return "ಕನ್ನಡ"
	case ko_KR:
		return "한국어 (韩国)"
	case lo_LA:
		return "ລາວ"
	case lt_LT:
		return "Lietuvių"
	case lv:
		return "Latviešu"
	case ml:
		return "മലയ" //
	case mn_MN:
		return "Монгол"
	case nb_NO:
		return "Norsk (bokmål)"
	case nl:
		return "Nederlands"
	case oc:
		return "Occitan"
	case pl_PL:
		return "Polski"
	case pt:
		return "Português"
	case pt_BR:
		return "Português (Brasil)"
	case ro_RO:
		return "Română"
	case ru:
		return "Русский"
	case sk_SK:
		return "Slovenčina (Slovakia)"
	case sl:
		return "Slovenščina"
	case sr:
		return "Српски"
	case sv_SE:
		return "Svenska"
	case ta:
		return "தமிழ்"
	case te:
		return "తెలుగు"
	case th:
		return "ภาษาไทย"
	case tr_TR:
		return "Türkçe"
	case uk_UA:
		return "Українська"
	case vi_VN:
		return "Tiếng Việt"
	case zh_CN:
		return "中文（中国）"
	case zh_TW:
		return "中文（台灣）"
	default:
		return ""
	}
}

func (c *Client) CreateCapture(short Language, external bool, host string, port int) (*pad.Pad, error) {
	lang := LanguageShortToName(short)

	//Subscribe to captions, pads and pads-sessions
	//Subscribe to captions
	if err := c.ddpSubscribe(bbb.CaptionsSub, nil); err != nil {
		return nil, err
	}
	captionsCollection := c.ddpClient.CollectionByName("captions")
	captionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads
	if err := c.ddpSubscribe(bbb.PadsSub, nil); err != nil {
		return nil, err
	}
	padsCollection := c.ddpClient.CollectionByName("pads")
	padsCollection.AddUpdateListener(c.ddpEventHandler)

	//Subscribe to pads-sessions
	if err := c.ddpSubscribe(bbb.PadsSessionsSub, nil); err != nil {
		return nil, err
	}
	padsSessionsCollection := c.ddpClient.CollectionByName("pads-sessions")
	padsSessionsCollection.AddUpdateListener(c.ddpEventHandler)

	//Create caption and add this bot as owner to it
	_, err := c.ddpCall(bbb.CreateGroupCall, string(short), "captions", lang)
	if err != nil {
		return nil, err
	}

	_, err = c.ddpCall(bbb.UpdateCaptionsOwnerCall, string(short), lang)
	if err != nil {
		return nil, err
	}

	//Get padID
	var padId string
	getPadIDtry := 0
	for {
		getPadIDtry++
		result, err := c.ddpCall(bbb.GetPadIdCall, string(short))
		if err != nil {
			return nil, err
		}

		if getPadIDtry > 10 {
			return nil, errors.New("timeout to call getPadId: " + err.Error())
		}

		if result == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		padId = result.(string)
		break
	}
	fmt.Println("padID: " + padId)

	_, err = c.ddpCall(bbb.CreateSessionCall, string(short))
	if err != nil {
		return nil, err
	}

	//Get sessionID
	var sessionID string
	getsessionIDtry := 0
	loop := true
	for loop {
		getsessionIDtry++
		result := padsSessionsCollection.FindAll()

		for _, element0 := range result {
			if element1, found := element0["sessions"]; found {
				if reflect.TypeOf(element1).Kind() == reflect.Slice {
					s := reflect.ValueOf(element1)
					for i := 0; i < s.Len(); i++ {
						element2 := s.Index(i)
						if element2.Kind() == reflect.Interface { //is Interface
							element3 := reflect.ValueOf(element2.Interface())
							if element3.Kind() == reflect.Map {
								for _, e := range element3.MapKeys() {
									element4 := element3.MapIndex(e)
									if convert.String(e) == string(short) {
										sessionID = convert.String(element4)
										loop = false
									}
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(100 * time.Millisecond)

		if (getsessionIDtry % 10) == 9 {
			fmt.Println("Retry to create and subscribe to pads-sessions")
			_, err = c.ddpCall(bbb.CreateSessionCall, string(short))
			if err != nil {
				fmt.Println("Failed to create pad session")
			}

			if err := c.ddpSubscribe(bbb.PadsSessionsSub, nil); err != nil {
				fmt.Println("Failed to subscribe to pads-sessions")
			}
		}

		if getsessionIDtry > 100 {
			return nil, errors.New("timeout to get sessionID")
		}
	}
	fmt.Println("sessionID: " + sessionID)

	capturePad := pad.NewPad(string(short), lang, c.PadURL, c.PadWSURL, c.SessionToken, padId, sessionID, c.SessionCookie, external, host, port)
	if err := capturePad.Connect(); err != nil {
		return nil, err
	}

	// Add capturePad to the list of pads
	c.padMutex.Lock()
	c.captures = append(c.captures, capturePad)
	c.padMutex.Unlock()

	capturePad.OnDisconnect(func() {
		// Remove capturePad from the list of pads
		c.padMutex.Lock()
		for i, p := range c.captures {
			if p == capturePad {
				c.captures = append(c.captures[:i], c.captures[i+1:]...)
				break
			}
		}
		c.padMutex.Unlock()
	})

	return capturePad, nil
}

func (c *Client) GetCaptures() []*pad.Pad {
	c.padMutex.Lock()
	defer c.padMutex.Unlock()

	return c.captures
}
