set dotenv-load

alias c := cargo-check
alias r := cargo-run
alias w := cargo-watch
alias wc := watch-css

default:
    @just --list

cargo-check:
    cargo c

cargo-run:
    cargo r

cargo-watch:
    cargo watch -x run

cargo-watch-html:
    cargo watch --workdir=templates -x run

watch-css:
    cd tailwind && npm run watch-css
