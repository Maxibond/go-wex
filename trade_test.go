package wex

import (
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const SLEEP = 5

var tapi = TradeAPI{}

func CheckAPI(t *testing.T) {
	time.Sleep(SLEEP * time.Second)
	if os.Getenv("API_KEY_TEST") == "" || os.Getenv("API_SECRET_TEST") == "" {
		t.SkipNow()
	}
}

func TestAccountInfo(t *testing.T) {

	CheckAPI(t)

	Convey("Account information data", t, func() {
		info, err := tapi.GetInfoAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"))
		Convey("No error should occur", func() {
			So(err, ShouldBeNil)
		})
		Convey("Account information fields should be returned", func() {
			So(info, ShouldHaveSameTypeAs, AccountInfo{})
			So(info.TransactionCount, ShouldEqual, 0)
			So(info.ServerTime, ShouldBeGreaterThan, 0)
			So(info.Funds["btc"], ShouldBeGreaterThanOrEqualTo, 0)
		})
	})

}

func TestActiveOrders(t *testing.T) {

	CheckAPI(t)

	Convey("Active orders data", t, func() {
		orders, err := tapi.ActiveOrdersAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), "btc_usd")

		if err != nil {
			Convey("If error is returned, it should be 'no orders'", func() {
				So(err, ShouldResemble, TradeError{msg: "no orders"})
			})
		} else {
			Convey("If no error is returned, 'order' should have length", func() {
				So(len(orders), ShouldBeGreaterThanOrEqualTo, 1)
			})
		}
	})
}

func TestOrderTrade(t *testing.T) {

	CheckAPI(t)

	Convey("Trade new order", t, func() {
		orderResponse, err := tapi.TradeAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), "btc_usd", "buy", 900, 1)

		if err != nil {
			Convey("If error is returned, it should be 'not enough USD'", func() {
				So(err, ShouldResemble, TradeError{msg: "It is not enough USD for purchase"})
			})
		} else {
			Convey("If no error is returned, 'btc_usd' amount should appear", func() {
				So(orderResponse.OrderID, ShouldBeGreaterThan, 0)
			})
		}
	})
}

func TestOrderInfo(t *testing.T) {

	CheckAPI(t)

	orderID := "1"
	Convey("Get order info", t, func() {
		orderResponse, err := tapi.OrderInfoAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), orderID)

		if err != nil {
			Convey("If error is returned, it should be 'invalid order'", func() {
				So(err, ShouldResemble, TradeError{msg: "invalid order"})
			})
		} else {
			Convey("If no error is returned, order information should be returned", func() {
				So(orderResponse[orderID], ShouldNotBeNil)
				So(orderResponse[orderID].Amount, ShouldBeGreaterThanOrEqualTo, 0)
			})
		}
	})
}

func TestCancelOrder(t *testing.T) {

	CheckAPI(t)

	orderID := "1"
	Convey("Cancel order", t, func() {
		orderResponse, err := tapi.CancelOrderAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), orderID)

		if err != nil {
			Convey("If error is returned, it should be 'bad status'", func() {
				So(err, ShouldResemble, TradeError{msg: "bad status"})
			})
		} else {
			Convey("If no error is returned, same order id should be returned", func() {
				So(strconv.Itoa(orderResponse.OrderID), ShouldEqual, orderID)
			})
		}
	})
}

func TestTradeHistory(t *testing.T) {

	CheckAPI(t)

	Convey("Trade history data", t, func() {

		filter := HistoryFilter{}
		tradeHistory, err := tapi.TradeHistoryAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), filter, "btc_usd")

		if err != nil {
			Convey("If error is returned, it should be 'no trades'", func() {
				So(err, ShouldResemble, TradeError{msg: "no trades"})
			})
		}

		Convey("Trade history should be retrieved", func() {
			So(tradeHistory, ShouldNotBeNil)
			So(len(tradeHistory), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestTransactionHistory(t *testing.T) {

	CheckAPI(t)

	Convey("Transaction history data", t, func() {

		filter := HistoryFilter{}
		transactionHistory, err := tapi.TransactionHistoryAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), filter)

		if err != nil {
			Convey("If error is returned, it should be 'no transactions'", func() {
				So(err, ShouldResemble, TradeError{msg: "no transactions"})
			})
		}

		Convey("Transaction history should be retrieved", func() {
			So(transactionHistory, ShouldNotBeNil)
			So(len(transactionHistory), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestWithdrawCoin(t *testing.T) {

	CheckAPI(t)

	Convey("Withdraw coin", t, func() {

		response, err := tapi.WithdrawCoinAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), "BTC", 0.001, "address")

		if err != nil {
			Convey("If error is returned, it should be 'api permission'", func() {
				So(err, ShouldResemble, TradeError{msg: "api key dont have withdraw permission"})
			})
		} else {
			Convey("If no error is returned, withdraw response should be returned", func() {
				So(response.TransactionID, ShouldBeGreaterThan, 0)
				So(response.AmountSent, ShouldBeGreaterThanOrEqualTo, 0)
			})
		}
	})
}

func TestCreateCoupon(t *testing.T) {

	CheckAPI(t)

	Convey("Create coupon", t, func() {

		response, err := tapi.CreateCouponAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), "BTC", 0.001)

		if err != nil {
			Convey("If error is returned, it should be 'api permission'", func() {
				So(err, ShouldResemble, TradeError{msg: "api key dont have coupon permission"})
			})
		} else {
			Convey("If no error is returned, withdraw response should be returned", func() {
				So(response.Coupon, ShouldNotBeBlank)
				So(response.TransactionID, ShouldBeGreaterThan, 0)
			})
		}
	})
}

func TestRedeemCoupon(t *testing.T) {

	CheckAPI(t)

	Convey("Redeem coupon", t, func() {

		response, err := tapi.RedeemCouponAuth(os.Getenv("API_KEY_TEST"), os.Getenv("API_SECRET_TEST"), "BTC-USD-XYZ")

		if err != nil {
			Convey("If error is returned, it should be 'api permission'", func() {
				So(err, ShouldResemble, TradeError{msg: "api key dont have coupon permission"})
			})
		} else {
			Convey("If no error is returned, withdraw response should be returned", func() {
				So(response.CouponCurrency, ShouldNotBeBlank)
				So(response.TransactionID, ShouldBeGreaterThan, 0)
			})
		}
	})
}
