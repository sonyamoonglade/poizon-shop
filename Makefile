.SILENT:
.PHONY:

run-clothes-bot:
	cd apps/clothes_bot && make

run-api:
	cd apps/api && make

ci:
	cd apps/api && make ci && cd ../..
	cd apps/clothes_bo && make ci && cd ../..
	cd apps/household_bot && make ci

test:
	cd apps/api && make unit-test
	cd apps/clothes_bot && make unit-test
	cd apps/household_bot && make unit-test
