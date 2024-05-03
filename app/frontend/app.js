import van from "vanjs-core"
const {p, div} = van.tags

const Hello = () => {
    const dom = div();

    return div(dom, p("hello world"));
}

van.add(document.body, Hello())