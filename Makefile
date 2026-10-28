.PHONY: check acceptance audit
check:
	sh scripts/check.sh
acceptance:
	python3 scripts/acceptance.py
audit:
	sh scripts/audit.sh
