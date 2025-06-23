activate_env:
	. ./.venv/bin/activate

start_prod:
	make activate_env
	uvicorn src.main:app --port 8000 --log-level info

start_dev:
	make activate_env
	uvicorn src.main:app --reload --port 8000 --log-level debug

test:
	make activate_env
	python3 -m unittest tests.main