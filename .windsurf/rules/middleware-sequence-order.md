---
trigger: always_on
---

Expected sequence order of middlewares:

1. First the before middlewars are executed
globals middlewares before → domains middlewares before → groups middlewares before → routes middlewares before

2. Then Handler is executed
handler

3. Last the after middlewares are excuted
routes middlewares after → groups middlewares after → domains middlewares after → globals middlewares after


Note!!!
For routes at the same level as groups, they should only get global and domain middlewares