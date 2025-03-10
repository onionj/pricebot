package price

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	url := fmt.Sprintf("https://call5.tgju.org/ajax.json?what=%d", ltime.Unix())

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

func (p Price) String() string {

	return fmt.Sprintf(`
ا🇺🇸 دلار امریکا (%s): %s ریال
ا🇪🇺 یورو اروپا (%s): %s ریال
ا🇬🇧 پوند انگلیس (%s): %s ریال
ا🇨🇦 دلار کانادا (%s): %s ریال
ا🇦🇺 دلار استرالیا (%s): %s ریال
ا🇦🇪 درهم امارات (%s): %s ریال

ا👑 بیتکوین (%s): %s دلار
ا🇺🇸 تتر (%s): %s ریال
ا💠 اتریوم (%s): %s دلار

ا🪙 سکه بهار آزادی (%s): %s ریال
ا🪙 سکه امامی (%s): %s ریال
ا🪙 نیم سکه (%s): %s ریال
ا🪙 رب سکه (%s): %s ریال
ا🪙 رب سکه قبل ۸۶ (%s): %s ریال

ا💰 طلا گرمی (%s): %s ریال
ا💰 مثقال طلا (%s): %s ریال
ا💰 انس طلا (%s): %s دلار

ا📆 اخرین بروزرسانی: %02d:%02d:%02d %s`,
		p.Current.Dollar.Time, p.Current.Dollar.Price,
		p.Current.Eur.Time, p.Current.Eur.Price,
		p.Current.GBP.Time, p.Current.GBP.Price,
		p.Current.CAD.Time, p.Current.CAD.Price,
		p.Current.AUD.Time, p.Current.AUD.Price,
		p.Current.AED.Time, p.Current.AED.Price,

		p.Current.BitCoin.Time, p.Current.BitCoin.Price,
		p.Current.Tether.Time, p.Current.Tether.Price,
		p.Current.Ethereum.Time, p.Current.Ethereum.Price,

		p.Current.SekeB.Time, p.Current.SekeB.Price,
		p.Current.SekeE.Time, p.Current.SekeE.Price,
		p.Current.Nim.Time, p.Current.Nim.Price,
		p.Current.Rob.Time, p.Current.Rob.Price,
		p.Current.RobDown.Time, p.Current.RobDown.Price,

		p.Current.Geram18.Time, p.Current.Geram18.Price,
		p.Current.Mesghal.Time, p.Current.Mesghal.Price,
		p.Current.Ons.Time, p.Current.Ons.Price,
		p.LastRefresh.Hour(), p.LastRefresh.Minute(), p.LastRefresh.Second(), p.JLastRefresh.String(),
	)
}
