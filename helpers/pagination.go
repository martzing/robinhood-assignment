package helpers

func Paginate[T interface{}](items *[]T, size int64) (int64, bool) {
	lenItems := int64(len(*items))
	if lenItems < size+1 {
		return lenItems, false
	}
	*items = (*items)[:lenItems-1]
	return lenItems - 1, true
}
