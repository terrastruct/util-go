package main

import (
	"os"
	"time"

	"oss.terrastruct.com/util-go/cmdlog"
	"oss.terrastruct.com/util-go/xos"
)

func main() {
	l := cmdlog.New(xos.NewEnv(os.Environ()), os.Stderr)
	l = l.WithCCPrefix("lochness")
	l = l.WithCCPrefix("imgbundler")
	l = l.WithCCPrefix("cache")

	l.NoLevel.Println("Somehow, the world always affects you more than you affect it.")

	l.SetDebug(true)
	l.Debug.Println("Man is a rational animal who always loses his temper when he is called upon.")

	l.SetDebug(false)
	l.Debug.Println("You can never trust a woman; she may be true to you.")

	l.SetTS(true)
	l.Success.Println("An alcoholic is someone you don't like who drinks as much as you do.")
	l.Info.Println("There once was this swami who lived above a delicatessan.")

	l.SetTSFormat(time.UnixDate)
	l.Warn.Println("Telephone books are like dictionaries -- if you know the answer before.")

	l.SetTS(false)
	l.Error.Println("Nothing can be done in one trip.")
}
