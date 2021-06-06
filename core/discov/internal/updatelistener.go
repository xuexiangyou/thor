package internal

type (
	KV struct {
		Key string
		Val string
	}

	UpdateListener interface {
		OnAdd(kv KV)
		OnDelete(kv KV)
	}
)
