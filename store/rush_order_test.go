package store

import (
	"blackledger/core"
	"testing"
)

func TestRushOrderReservationAndReceiptSurviveReload(t *testing.T) {
	s := testStore(t)
	if err := s.Change(func(w *core.World) error {
		w.Player.Location = "laundry"
		w.Event = nil
		w.Plots = nil
		w.Tasks = nil
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	w, _ := s.Read()
	cash := w.Player.Cash
	command := core.Command{Kind: "rushorder", Target: "laundry", RequestID: "rush-order-deduplicated", Revision: w.Revision}
	first, err := s.Command(command)
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.Command(command)
	if err != nil || string(first) != string(second) {
		t.Fatal("duplicate receipt changed", err)
	}
	loaded, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Player.Cash != cash+core.RushOrderPay || loaded.Properties["laundry"].RushOrderDay == 0 {
		t.Fatal("payment/reservation lost")
	}
	if _, err = s.Command(core.Command{Kind: "rushorder", Target: "laundry", RequestID: "rush-order-second-attempt", Revision: loaded.Revision}); err == nil {
		t.Fatal("reload permits second same-day order")
	}
}

func TestVersion18UpgradeKeepsExistingLaundryResources(t *testing.T) {
	s := testStore(t)
	if err := s.Change(func(w *core.World) error {
		w.Version = 18
		w.Player.Location = "laundry"
		w.Properties["laundry"].Supply = 11
		w.Properties["laundry"].Condition = 67
		w.Properties["laundry"].RushOrderDay = 0
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	w, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if w.Version != core.SaveVersion || w.Properties["laundry"].Supply != 11 || w.Properties["laundry"].Condition != 67 || w.Properties["laundry"].RushOrderDay != 0 {
		t.Fatal("upgrade reset resources or invented a reservation")
	}
	if why := w.RushOrderReadiness("laundry"); why != "" {
		t.Fatal("old save cannot take its first order:", why)
	}
}
