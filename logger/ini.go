package logger

func init() {
	if err := InitZapCore(nil); err != nil {
		panic(err)
	}
}
