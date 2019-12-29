[![GoDoc](https://godoc.org/github.com/MagicalTux/idlock?status.svg)](https://godoc.org/github.com/MagicalTux/idlock)

# IntLock, etc

Very simple lock-by-id object for various purposes. Keep multiple locks under arbitrary ID values.

Example:

```Go
import (
	"github.com/MagicalTux/idlock"
	"fmt"
)

func test(lk *idlock.IntLock, c chan struct{}) {
	lk.Lock(2)
	defer lk.Unlock(2)

	fmt.Println("In test")

	close(c)
}

func main() {
	lk := idlock.NewInt()

	lk.Lock(1, 2)

	c := make(chan struct{})

	go test(lk, c)

	fmt.Println("Unlocking ...")
	lk.Unlock(2)

	<-c

	lk.Unlock(1)
}
```
