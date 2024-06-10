package exception

func PanicIfNeeded(err any) {
	if err != nil {
		panic(err)
	}
}
