# tools-go

Some useful tools that I have used for my own project, install with:

```bash
go get github.com/chikaku/tools-go
```

### Overview

- tools/network:
  - tcp tunnel
  - get public mapped address by STUN
- tools/rt(**Linux Only**):
  - get runtime information like Goid/GoStack/GoStackSize/NumTimer
  - get `Offset` from unexposed structure (can be used for runtime structure)
