# This is a monorepo for 3 different apps for business. 
Shopping in [Poison.com](https://poison.com) with Telegram

Master branch of repo is considered production ready and currently deployed version
## Contains:

- API server to administrate system, see apps/api
- [Household bot](https://t.me/xKK_ru_techbot). Sells technical furniture (*computers, phones, laptops etc...*),
  see apps/household_bot
- [Clothing bot](https:/t.me/xKK_ru_bot). Sells clothes (*sneakers, bags, jackets, etc...*), see apps/clothing_bot

Consider household_bot to be more modern implementation. It's more clean, more naive and better techniques were applied, comparing it to clothing_bot. 

Application uses **MongoDB** as it's primary database.\
For http request/response handlers it's **Fiber** framework.\
To cache hot data **Redis** is used, as well as message broker.\
Apps communicate by **Redis Pub/Sub**\
Moreover, we store static files and thumbnails in **Yandex Cloud S3**,
uploads are done via self-written *Yandex S3 Client*\
As for deploying our app we use Yandex Compute Cloud VMs\
We use React and ElectronJS for client's to admin panel. This code is not public.

Majority of synchronous domain code is written in semi-functional style. It's clean, descriptive and expressive.\
Domain is well architected. Everything goes around it. Codebase is highly testable, every dependency can be mocked out.
