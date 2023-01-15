package maleo

import (
	"reflect"
	"testing"
	"time"
)

func TestMessageParameters(t *testing.T) {
	mal, _ := NewTestingMaleo()
	paramsGen := func() *MessageParameters {
		return &MessageParameters{
			ForceSend:  false,
			Messengers: Messengers{},
			Benched:    Messengers{newMockMessengerWithName(0, "zz")},
			Cooldown:   0,
			Maleo:      mal,
		}
	}
	trueFunc := func(Messenger) bool {
		return true
	}
	params := paramsGen()
	defer func() {
		if t.Failed() {
			t.Logf("params: %+v", params)
		}
	}()
	opts := Option.Message().
		ForceSend(true).
		Filter(trueFunc, trueFunc).
		IncludeBenched().
		IncludeBenchedFilter(trueFunc).
		IncludeBenchedName("zz").
		IncludeBenchedPrefix("whoo").
		IncludeBenchedSuffix("yes").
		IncludeBenchedContains("zzz").
		Messengers(newMockMessengerWithName(0, "mock"), newMockMessengerWithName(0, "mock2")).
		Cooldown(time.Second).
		FilterName().
		FilterName("mock")

	opts.Apply(params)

	if !params.ForceSend {
		t.Error("ForceSend should be true")
	}
	if len(params.Messengers) != 1 {
		t.Error("Messengers should have 1 element")
	}
	if params.Messengers[0].Name() != "mock" {
		t.Error("Messenger should be mock")
	}
	if params.Cooldown != time.Second {
		t.Error("Cooldown should be 1 second")
	}
	params2 := params.clone()
	if !reflect.DeepEqual(params, params2) {
		t.Error("clone should be equal")
	}

	opts = opts.Cooldown(time.Second * -1)
	opts.Apply(params)
	if params.Cooldown != 15*time.Minute {
		t.Error("Cooldown should be 15 minutes")
	}

	opts.Include(newMockMessengerWithName(0, "mock3")).Apply(params)
	if len(params.Messengers) != 2 {
		t.Error("Messengers should have 2 elements")
	}

	opts.ExcludeName("mock3").Apply(params)
	if len(params.Messengers) != 1 {
		t.Error("Messengers should have 1 elements")
	}
	opts.Include(newMockMessengerWithName(0, "mock5")).ExcludePrefix("mock").Apply(params)
	if len(params.Messengers) != 0 {
		t.Error("Messengers should have 0 elements")
	}
	opts.Messengers(newMockMessengerWithName(0, "mock5")).ExcludeSuffix("5").Apply(params)
	if len(params.Messengers) != 0 {
		t.Error("Messengers should have 0 elements")
	}

	var ctx MessageContext = errorMessageContext{
		Error: nil,
		param: params,
	}
	if ctx.Err() != nil {
		t.Error("Err should be nil")
	}
	if ctx.ForceSend() != true {
		t.Error("ForceSend should be true")
	}
	if ctx.Cooldown() != time.Minute*15 {
		t.Error("Cooldown should be 15 minutes")
	}
	if !reflect.DeepEqual(ctx.Maleo(), mal) {
		t.Error("Maleo should be equal")
	}
	ctx = entryMessageContext{nil, params}
	if ctx.Err() != nil {
		t.Error("Err should be nil")
	}
	if ctx.ForceSend() != true {
		t.Error("ForceSend should be true")
	}
	if ctx.Cooldown() != time.Minute*15 {
		t.Error("Cooldown should be 15 minutes")
	}
	if !reflect.DeepEqual(ctx.Maleo(), mal) {
		t.Error("Maleo should be equal")
	}
}
