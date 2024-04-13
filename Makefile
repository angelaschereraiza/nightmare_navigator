.ONESHELL:
.PHONY: deploy
deploy:
	go build -o nightmare_navigator main.go 
	sync
	strip nightmare_navigator
	ssh aiza.ch mkdir -p nightmare_navigator
	ssh aiza.ch	pkill -f nightmare_navigator || true
	scp nightmare_navigator aiza.ch:~/nightmare_navigator
	ssh -T aiza.ch << 'EOF'
	cd nightmare_navigator
	./nightmare_navigator > nightmare_navigator.log 2>&1 &
	EOF

.PHONY: stop
stop:
	ssh aiza.ch	pkill -f nightmare_navigator || true