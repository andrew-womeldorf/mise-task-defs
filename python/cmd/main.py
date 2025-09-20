import typer

from foo import hello, goodbye

app = typer.Typer()


@app.command()
def hi(name: str):
    hello(name)


@app.command()
def bye(name: str, formal: bool = False):
    goodbye(name, format)


if __name__ == "__main__":
    app()
