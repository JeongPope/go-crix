package exchange

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jeongpope/go-crix/logger"
	"github.com/jeongpope/go-crix/model"
	"github.com/jeongpope/go-crix/utils"
)

type Upbit struct {
	exchange
}

// UpbitTickerEvent define websocket ticker statistics event
type UpbitTickerEvent struct {
	Type               string  `json:"type"`                  // 타입(ticker : 현재가)
	Code               string  `json:"code"`                  // 마켓 코드 (ex. KRW-BTC)
	OpeningPrice       float64 `json:"opening_price"`         // 시가
	HighPrice          float64 `json:"high_price"`            // 고가
	LowPrice           float64 `json:"low_price"`             // 저가
	TradePrice         float64 `json:"trade_price"`           // 현재가
	PrevClosingPrice   float64 `json:"prev_closing_price"`    // 전일 종가
	Change             string  `json:"change"`                // 전일 대비
	ChangePrice        float64 `json:"change_price"`          // 부호 없는 전일 대비 값
	SignedChangePrice  float64 `json:"signed_change_price"`   // 전일 대비 값
	ChangeRate         float64 `json:"change_rate"`           // 부호 없는 전일 대비 등락율
	SignedChangeRate   float64 `json:"signed_change_rate"`    // 전일 대비 등락율
	TradeVolume        float64 `json:"trade_volume"`          // 가장 최근 거래량
	AccTradeVolume     float64 `json:"acc_trade_volume"`      // 누적 거래량(UTC 0시 기준)
	AccTradeVolume24h  float64 `json:"acc_trade_volume_24h"`  // 24시간 누적 거래량
	AccTradePrice      float64 `json:"acc_trade_price"`       // 누적 거래대금(UTC 0시 기준)
	AccTradePrice24h   float64 `json:"acc_trade_price_24h"`   // 24시간 누적 거래대금
	TradeDate          string  `json:"trade_date"`            // 최근 거래 일자(UTC)
	TradeTime          string  `json:"trade_time"`            // 최근 거래 시각(UTC)
	TradeTimestamp     float64 `json:"trade_timestamp"`       // 체결 타임스탬프 (milliseconds)
	AskBid             string  `json:"ask_bid"`               // 매수/매도 구분
	AccAskVolume       float64 `json:"acc_ask_volume"`        // 누적 매도량
	AccBidVolume       float64 `json:"acc_bid_volume"`        // 누적 매수량
	Highest52WeekPrice float64 `json:"highest_52_week_price"` // 52주 최고가
	Highest52WeekDate  string  `json:"highest_52_week_date"`  // 52주 최고가 달성일
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`  // 52주 최저가
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`   // 52주 최저가 달성일
	TradeStatus        string  `json:"trade_status"`          // 거래상태
	MarketState        string  `json:"market_status"`         // 거래상태
	MarketStateForIOS  string  `json:"market_state_for_ios"`  // 거래 상태
	IsTradingSuspended bool    `json:"is_trading_suspended"`  // 거래 정지 여부
	DelistingDate      string  `json:"delisting_date"`        // 상장폐지일
	MarketWarning      string  `json:"market_warning"`        // 유의 종목 여부
	Timestamp          uint    `json:"timestamp"`             // 타임스탬프 (milliseconds)
	StreamType         string  `json:"stream_type"`           // 스트림 타입
}

// UpbitTicketField define upbit websocket ticket field statistics event
type UpbitTicketField struct {
	Ticket string `json:"ticket,omitempty"`
}

// UpbitTypeField define upbit websocket type field statistics event
type UpbitTypeField struct {
	Type  string   `json:"type,omitempty"`
	Codes []string `json:"codes,omitempty"`
}

// UpbitMarket define upbit market information
type UpbitMarket struct {
	Market        string `json:"market"`         // 업비트 제공 시장 정보
	Korean        string `json:"korean_name"`    // 한글명
	English       string `json:"english_name"`   // 영문명
	MarketWarning string `json:"market_warning"` // 유의 종목 여부
}

func (ex *Upbit) Initialize(currencies *[]string) []string {
	logger.Log.Info("[upbit.go] Start Initialize()")

	codes, coins := ex.getMarkets()
	msg := ex.makeSubsMessage(codes)

	ex.tickerEndpoint = tickerURL
	ex.c = nil
	ex.reconnectLock = &sync.Mutex{}
	ex.subsMessage = append(ex.subsMessage, &msg)
	ex.chanSendMessage = nil

	ex.supportAssets = append(ex.supportAssets, coins...)
	ex.tickers = make(map[string]model.Ticker)
	ex.updateLock = nil

	ex.initSnapshot(codes)

	logger.Log.Info("[upbit.go] End Initialize()")

	return ex.supportAssets
}

func (ex *Upbit) Execute() (err error) {
	logger.Log.Info("[upbit.go] Start Execute()")

	handler := func(event *UpbitTickerEvent) {
		name := event.Code[4:]

		if ex.tickers[name].Price != event.TradePrice {
			tempTicker := model.Ticker{
				Exchange:       "UPBIT",
				Currency:       event.Code[4:],
				Price:          event.TradePrice,
				YesterdayPrice: event.PrevClosingPrice,
				Change:         event.SignedChangePrice,
				ChangeRate:     event.SignedChangeRate,
				Volume:         uint(event.AccTradePrice24h),
			}

			ex.tickers[name] = tempTicker
			logger.Log.Info("[TICKER] ", tempTicker)

			ex.chanSendMessage <- tempTicker
		}
	}

	errHandler := func(err error) {
		logger.Log.Error("Func subscribeTicker() return error : ", err)
	}

	// Serve
	wsHandler := func(message []byte) {
		event := new(UpbitTickerEvent)
		err := json.Unmarshal(message, event)

		if err != nil {
			errHandler(err)
			return
		}

		handler(event)
	}

	return websocketServe(ex.c, ex.reconnectLock,
		ex.tickerEndpoint, ex.subsMessage, wsHandler, errHandler)
}

func (ex *Upbit) Release() {
	logger.Log.Info("[upbit.go] Start Release()")

	if ex.c != nil {
		ex.c.Close()
	}

	logger.Log.Info("[upbit.go] End Release())")
}

func (ex *Upbit) initSnapshot(codes []string) {
	logger.Log.Info("[upbit.go] Start initSnapshot()")

	var markets string
	for _, v := range codes {
		markets += v + ","
	}

	markets = markets[:len(markets)-1]

	resp, err := http.Get(snapshotURL + markets)
	if err != nil {
		logger.Log.Error("Upbit updateSnapshot() failed")
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Failed read response body
		logger.Log.Error("Upbit updateSnapshot() read response body return err, retry after 10 second")
		return
	}

	// Parse JSON
	var f []interface{}
	json.Unmarshal(data, &f)
	if err != nil {
		logger.Log.Error("Upbit updateSnapshot() error parsing JSON: ", err)
		return
	}

	// Loop through the items
	for _, value := range f {
		dataMap := value.(map[string]interface{})

		currency := dataMap["market"].(string)[len("KRW-"):]

		ex.tickers[currency] = model.Ticker{
			Exchange:       "UPBIT",
			Currency:       currency,
			Price:          utils.ToFloat64(dataMap["trade_price"]),
			YesterdayPrice: utils.ToFloat64(dataMap["prev_closing_price"]),
			Change:         utils.ToFloat64(dataMap["signed_change_price"]),
			ChangeRate:     utils.ToFloat64(dataMap["signed_change_rate"]),
			Volume:         uint(utils.ToFloat64(dataMap["acc_trade_price_24h"])),
		}
	}

	logger.Log.Info("[upbit.go] initSnapshot close")
}

// -----
func (ex *Upbit) getMarkets() ([]string, []string) {
	var resp *http.Response
	var err error
	var data []byte
	var codes []string
	var currencies []string

	var markets []UpbitMarket

	for {
		resp, err = http.Get(marketURL)
		if err != nil {
			logger.Log.Error("Upbit getMarkets() restServe return err, retry after 10 second")

			time.Sleep(time.Second * 10)
			continue
		} else {
			data, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				logger.Log.Error("Upbit getMarkets() read response body return err, retry after 10 second")

				resp.Body.Close()
				time.Sleep(time.Second * 10)
			} else {
				err = json.Unmarshal(data, &markets)
				if err != nil {
					logger.Log.Error("Upbit getMarkets() unmarshal response body return err, retry after 10 second")

					resp.Body.Close()
					time.Sleep(time.Second * 10)
				}

				for _, market := range markets {
					if strings.HasPrefix(market.Market, "KRW-") {
						codes = append(codes, market.Market)
						currencies = append(currencies, market.Market[4:])
					}
				}
			}
		}

		break
	}
	defer resp.Body.Close()

	return codes, currencies
}

func makeSignature() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	key, _ := base64.StdEncoding.DecodeString("t" + timestamp)
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(signature_key))

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	logger.Log.Info("Upbit signature : ", signature)

	return signature
}

func (ex *Upbit) makeSubsMessage(codes []string) []byte {
	// Subscribe message format
	// [{ticket}, {type}, {format}]
	msg := []interface{}{}
	i := UpbitTicketField{makeSignature()}
	j := UpbitTypeField{
		Type:  "ticker",
		Codes: codes,
	}
	msg = append(msg, i, j)
	tickerMsg, _ := json.Marshal(msg)

	return tickerMsg
}
