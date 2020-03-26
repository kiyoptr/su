package sharder

func removeIntSlice(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
