activate_env:
	. ./.venv/bin/activate

test:
	make activate_env
	python3 -m unittest tests.main