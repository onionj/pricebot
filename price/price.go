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

// Make these package variables so they can be modified in tests
var (
	httpClient = &http.Client{}
	baseURL    = "https://call3.tgju.org/ajax.json"
)

type Detail struct {
	Price            string  `json:"p"`
	DateTime         string  `json:"ts"`
	Time             string  `json:"t"`
	ChangePercentage float64 `json:"dp"`
	ChangeDirection  string  `json:"dt"` // low, high
}

func (d Detail) FormatChange() string {

	percentage := fmt.Sprintf("%.2f%%", d.ChangePercentage)
	FormatPercentage := "(%s%s%s)"

	if d.ChangePercentage == 0 {
		return "⬅️"
	}

	switch d.ChangeDirection {
	case "high":
		return fmt.Sprintf(FormatPercentage, percentage, "+", "🟢")
	case "low":
		return fmt.Sprintf(FormatPercentage, percentage, "-", "🔴")
	default:
		return "⬅️"
	}
}

type CurrentData struct {
	Dollar Detail `json:"price_dollar_rl"`
	Eur    Detail `json:"price_eur"`
	GBP    Detail `json:"price_gbp"`
	CAD    Detail `json:"price_cad"`
	AUD    Detail `json:"price_aud"`
	AED    Detail `json:"price_aed"`
	TRY    Detail `json:"price_try"`
	SEK    Detail `json:"price_sek"`
	CNY    Detail `json:"price_cny"`
	SAR    Detail `json:"price_sar"`
	IQD    Detail `json:"price_iqd"`

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
	url := fmt.Sprintf("%s?what=%d", baseURL, ltime.Unix())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept-Language", "fa-IR")

	resp, err := httpClient.Do(req)
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

ا🇺🇸 دلار امریکا (%s) %s *%s* تومان
ا🇪🇺 یورو اروپا (%s) %s *%s* تومان
ا🇬🇧 پوند انگلیس (%s) %s *%s* تومان
ا🇨🇦 دلار کانادا (%s) %s *%s* تومان
ا🇦🇺 دلار استرالیا (%s) %s *%s* تومان
ا🇦🇪 درهم امارات (%s) %s *%s* تومان
ا🇹🇷 لیر ترکیه (%s) %s *%s* تومان
ا🇸🇪 کرون سوئد (%s) %s *%s* تومان
ا🇨🇳 یوان چین (%s) %s *%s* تومان
ا🇸🇦 ریال عربستان (%s) %s *%s* تومان
ا🇮🇶 دینار عراق (%s) %s *%s* ریال

ا👑 بیتکوین (%s) %s *%s* دلار
ا🇺🇸 تتر (%s) %s *%s* تومان
ا💠 اتریوم (%s) %s *%s* دلار

ا🪙 سکه بهار آزادی (%s) %s *%s* تومان
ا🪙 سکه امامی (%s) %s *%s* تومان
ا🪙 نیم سکه (%s) %s *%s* تومان
ا🪙 ربع سکه (%s) %s *%s* تومان
ا🪙 ربع سکه قبل ۸۶ (%s) %s *%s* تومان

ا💰 طلا گرمی (%s) %s *%s* تومان
ا💰 مثقال طلا (%s) %s *%s* تومان
ا💰 انس طلا (%s) %s *%s* دلار`,
		p.LastRefresh.Hour(), p.LastRefresh.Minute(), p.LastRefresh.Second(), p.JLastRefresh.String(),
		p.Current.Dollar.Time, p.Current.Dollar.FormatChange(), p.toToman(p.Current.Dollar.Price),
		p.Current.Eur.Time, p.Current.Eur.FormatChange(), p.toToman(p.Current.Eur.Price),
		p.Current.GBP.Time, p.Current.GBP.FormatChange(), p.toToman(p.Current.GBP.Price),
		p.Current.CAD.Time, p.Current.CAD.FormatChange(), p.toToman(p.Current.CAD.Price),
		p.Current.AUD.Time, p.Current.AUD.FormatChange(), p.toToman(p.Current.AUD.Price),
		p.Current.AED.Time, p.Current.AED.FormatChange(), p.toToman(p.Current.AED.Price),
		p.Current.TRY.Time, p.Current.TRY.FormatChange(), p.toToman(p.Current.TRY.Price),
		p.Current.SEK.Time, p.Current.SEK.FormatChange(), p.toToman(p.Current.SEK.Price),
		p.Current.CNY.Time, p.Current.CNY.FormatChange(), p.toToman(p.Current.CNY.Price),
		p.Current.SAR.Time, p.Current.SAR.FormatChange(), p.toToman(p.Current.SAR.Price),
		p.Current.IQD.Time, p.Current.IQD.FormatChange(), p.Current.IQD.Price,

		p.Current.BitCoin.Time, p.Current.BitCoin.FormatChange(), p.Current.BitCoin.Price,
		p.Current.Tether.Time, p.Current.Tether.FormatChange(), p.toToman(p.Current.Tether.Price),
		p.Current.Ethereum.Time, p.Current.Ethereum.FormatChange(), p.Current.Ethereum.Price,

		p.Current.SekeB.Time, p.Current.SekeB.FormatChange(), p.toToman(p.Current.SekeB.Price),
		p.Current.SekeE.Time, p.Current.SekeE.FormatChange(), p.toToman(p.Current.SekeE.Price),
		p.Current.Nim.Time, p.Current.Nim.FormatChange(), p.toToman(p.Current.Nim.Price),
		p.Current.Rob.Time, p.Current.Rob.FormatChange(), p.toToman(p.Current.Rob.Price),
		p.Current.RobDown.Time, p.Current.RobDown.FormatChange(), p.toToman(p.Current.RobDown.Price),

		p.Current.Geram18.Time, p.Current.Geram18.FormatChange(), p.toToman(p.Current.Geram18.Price),
		p.Current.Mesghal.Time, p.Current.Mesghal.FormatChange(), p.toToman(p.Current.Mesghal.Price),
		p.Current.Ons.Time, p.Current.Ons.FormatChange(), p.Current.Ons.Price,
	)
}
