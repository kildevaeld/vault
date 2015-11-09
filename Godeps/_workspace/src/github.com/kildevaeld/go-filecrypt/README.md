# go-filecrypt
Streaming file encryption based on nacl/secretbox

# Usage:
**Encrypt file**
```go


key := filecrypt.Key([]byte("some super secret key"))

file, _ := os.Open("main.go")
defer file.Close()

encrypted, _ := os.Create("main.vau")
defer file.Close()
filecrypt.Encrypt(encrypted, file, &key)

encrypted.Seek(0, 0)

buf := bytes.NewBuffer(nil)
filecrypt.Decrypt(buf, encrypted, &key)

fmt.Printf("%v", string(buf.Bytes()))
```
** Encrypt message **
```go
key := filecrypt.Key([]byte("some super secret key"))

emsg, ee := filecrypt.EncryptMessage([]byte("a message", &key, nil)

msg, ed := filecrypt.DecryptMessage(emsg, &key)

fmt.Printf(string(msg))

```
