# Simple test data

Given original data

![original](simple.png)

```console
$ < simple.json | jq -c '.[]' | sigo -q x,y | jq -s > simple_sigo.json
```

![masked](simple-sigo.png)
