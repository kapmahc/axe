dist=dist

build:
	cargo build --release
	mkdir -p $(dist)
	-cp -r target/release/axe package.json package-lock.json db locales templates themes README.md LICENSE $(dist)/
	strip -s $(dist)/axe
	cd $(dist) && tar cfJ ../$(dist).tar.xz *


clean:
	cargo clean
	-rm -r $(dist) $(dist).tar.xz


init:
	rustup update
	cargo update
	npm install
