package price

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/onionj/pricebot/utils"
)

type Detail struct {
	Price    string `json:"p"`
	DateTime string `json:"ts"`
	Time     string `json:"t"`
}

type CurrentData struct {
	Dollar Detail `json:"price_dollar_rl"`
	Eur    Detail `json:"price_eur"`
	GBP    Detail `json:"price_gbp"`
	CAD    Detail `json:"price_cad"`
	AUD    Detail `json:"price_aud"`
	AED    Detail `json:"price_aed"`
	TRY    Detail `json:"price_try"`

	Tether   Detail `json:"crypto-tether-irr"`
	BitCoin  Detail `json:"crypto-bitcoin"`
	Ethereum Detail `json:"crypto-ethereum"`

	SekeB   Detail `json:"sekeb"`
	SekeE   Detail `json:"sekee"`
	Nim     Detail `json:"nim"`
	Rob     Detail `json:"rob"`
	RobDown Detail `json:"rob_down"`

	Geram18 Detail `json:"geram18"`
	Mesghal Detail `json:"mesghal"`
	Ons     Detail `json:"ons"`
}

type Price struct {
	Current      CurrentData `json:"current"`
	LastRefresh  time.Time
	JLastRefresh utils.JDate
}

func NewPrice() *Price {
	return &Price{}
}

func (p *Price) Refresh() error {
	loc, _ := time.LoadLocation("Asia/Tehran")
	ltime := time.Now().In(loc)

	// ‍‍`what` just for deactivate cache!
	url := fmt.Sprintf("https://call3.tgju.org/ajax.json?what=%d", ltime.Unix())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept-Language", "fa-IR")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(body, p); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	p.LastRefresh = ltime
	p.JLastRefresh = utils.GregorianToJalali(p.LastRefresh.Year(), int(p.LastRefresh.Month()), p.LastRefresh.Day())
	return nil
}

func (p Price) prettyNumber(i int) string {
	s := strconv.Itoa(i)
	r1 := ""
	idx := 0

	// Reverse and interleave the separator.
	for i = len(s) - 1; i >= 0; i-- {
		idx++
		if idx == 4 {
			idx = 1
			r1 = r1 + ","
		}
		r1 = r1 + string(s[i])
	}

	// Reverse back and return.
	r2 := ""
	for i = len(r1) - 1; i >= 0; i-- {
		r2 = r2 + string(r1[i])
	}
	return r2
}

func (p Price) toToman(rilaString string) string {
	rilaInt, err := strconv.Atoi(strings.Replace(rilaString, ",", "", 10))

	if err != nil {
		return "0"
	}
	return p.prettyNumber(rilaInt / 10)
}

func (p Price) String() string {

	return fmt.Sprintf(`ا📆 اخرین بروزرسانی: %02d:%02d:%02d %s

ا🇺🇸 دلار امریکا (%s) ⬅️ %s تومان
ا🇪🇺 یورو اروپا (%s) ⬅️ %s تومان
ا🇬🇧 پوند انگلیس (%s) ⬅️ %s تومان
ا🇨🇦 دلار کانادا (%s) ⬅️ %s تومان
ا🇦🇺 دلار استرالیا (%s) ⬅️ %s تومان
ا🇦🇪 درهم امارات (%s) ⬅️ %s تومان
ا🇹🇷 لیر ترکیه (%s) ⬅️ %s تومان

ا👑 بیتکوین (%s) ⬅️ %s دلار
ا🇺🇸 تتر (%s) ⬅️ %s تومان
ا💠 اتریوم (%s) ⬅️ %s دلار

ا🪙 سکه بهار آزادی (%s) ⬅️ %s تومان
ا🪙 سکه امامی (%s) ⬅️ %s تومان
ا🪙 نیم سکه (%s) ⬅️ %s تومان
ا🪙 رب سکه (%s) ⬅️ %s تومان
ا🪙 رب سکه قبل ۸۶ (%s) ⬅️ %s تومان

ا💰 طلا گرمی (%s) ⬅️ %s تومان
ا💰 مثقال طلا (%s) ⬅️ %s تومان
ا💰 انس طلا (%s) ⬅️ %s دلار`,
		p.LastRefresh.Hour(), p.LastRefresh.Minute(), p.LastRefresh.Second(), p.JLastRefresh.String(),
		p.Current.Dollar.Time, p.toToman(p.Current.Dollar.Price),
		p.Current.Eur.Time, p.toToman(p.Current.Eur.Price),
		p.Current.GBP.Time, p.toToman(p.Current.GBP.Price),
		p.Current.CAD.Time, p.toToman(p.Current.CAD.Price),
		p.Current.AUD.Time, p.toToman(p.Current.AUD.Price),
		p.Current.AED.Time, p.toToman(p.Current.AED.Price),
		p.Current.TRY.Time, p.toToman(p.Current.TRY.Price),

		p.Current.BitCoin.Time, p.Current.BitCoin.Price,
		p.Current.Tether.Time, p.toToman(p.Current.Tether.Price),
		p.Current.Ethereum.Time, p.Current.Ethereum.Price,

		p.Current.SekeB.Time, p.toToman(p.Current.SekeB.Price),
		p.Current.SekeE.Time, p.toToman(p.Current.SekeE.Price),
		p.Current.Nim.Time, p.toToman(p.Current.Nim.Price),
		p.Current.Rob.Time, p.toToman(p.Current.Rob.Price),
		p.Current.RobDown.Time, p.toToman(p.Current.RobDown.Price),

		p.Current.Geram18.Time, p.toToman(p.Current.Geram18.Price),
		p.Current.Mesghal.Time, p.toToman(p.Current.Mesghal.Price),
		p.Current.Ons.Time, p.Current.Ons.Price,
	)
}
